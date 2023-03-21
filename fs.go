package wfile

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type EventCode int

const (
	sysInterval  = time.Millisecond * 500
	syncInterval = time.Millisecond * 1000
	maxFileSize  = 40 << 10
)

type DirHandler func(path string, d fs.DirEntry) error

// FS stores the information we need for the cache
type FS struct {
	root     string
	interval time.Duration
	files    *sync.Map
}

// File is just a file path to watch.
// also tracks the last checksum
type File struct {
	path string
	last [md5.Size]byte
}

// NewFS acts as a cache.
// keeping track of the files to watch.
func NewFS(root string) *FS {
	return &FS{
		root:  root,
		files: &sync.Map{},
	}
}

// Add adds an item to the cache
func (ffs *FS) Add(key, value any) {
	if _, exists := ffs.files.Load(key); exists {
		return
	}
	ffs.files.Store(key, value)
}

// Delete removes an item from the cache
func (ffs *FS) Delete(key any) {
	if _, exists := ffs.files.Load(key); !exists {
		return
	}
	ffs.files.Delete(key)
}

// Iter creates a map[any]any representation of the cache
func (ffs *FS) Iter() map[any]any {
	mp := make(map[any]any)

	ffs.files.Range(func(key, value any) bool {
		mp[key] = value
		return true
	})
	return mp
}

// Exists checks if a key exists in the cache
func (ffs *FS) Exists(key any) bool {
	_, exists := ffs.files.Load(key)
	return exists
}

// Size returns the item count of the cache
func (ffs *FS) Size() int {
	return len(ffs.Iter())
}

// WalkDir walks through the root directory, recursively
func WalkDir(path string, handler DirHandler) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("walkdir read path error: %v", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if err := WalkDir(entryPath, handler); err != nil {
				return fmt.Errorf("walkdir recursion error: %v", err)
			}
		} else {
			if err := handler(entryPath, entry); err != nil {
				return fmt.Errorf("walkdir handler error: %v", err)
			}
		}
	}
	return nil
}

// Update runs and collects any new or removed files in the root directory
func (ffs *FS) Update() {
	err := WalkDir(ffs.root, func(path string, d fs.DirEntry) error {
		sum, _ := Checksum(path)
		if !d.IsDir() {
			ffs.Add(path, &File{
				path: path,
				last: sum,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println("fs update error:", err)
		return
	}
}
