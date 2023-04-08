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

type DirHandler func(path string, d fs.DirEntry) error

// Monitor stores the information we need for the cache
type Monitor struct {
	root     string
	interval time.Duration
	files    *sync.Map
}

// File is just a file Path to watch.
// also tracks the last checksum
type File struct {
	path string
	last [md5.Size]byte
}

// NewMonitor acts as a cache.
// keeping track of the files to watch.
func NewMonitor(root string) *Monitor {
	return &Monitor{
		root:  root,
		files: &sync.Map{},
	}
}

// Add adds an item to the cache
func (m *Monitor) Add(key, value any) {
	if _, exists := m.files.Load(key); exists {
		fmt.Println(key, "exists")
		return
	}
	m.files.Store(key, value)
}

// Delete removes an item from the cache
func (m *Monitor) Delete(key any) {
	if _, exists := m.files.Load(key); !exists {
		return
	}
	m.files.Delete(key)
}

// ExportFileMap creates a map[any]any representation of the cache
func (m *Monitor) ExportFileMap() map[any]any {
	mp := make(map[any]any)

	m.files.Range(func(key, value any) bool {
		mp[key] = value
		return true
	})
	return mp
}

// Exists checks if a key exists in the cache
func (m *Monitor) Exists(key any) bool {
	_, exists := m.files.Load(key)
	return exists
}

// Size returns the item count of the cache
func (m *Monitor) Size() int {
	return len(m.ExportFileMap())
}

// WalkDir walks through the root directory, recursively
func WalkDir(path string, handler DirHandler) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("walkdir read Path Error: %v", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if err := WalkDir(entryPath, handler); err != nil {
				return fmt.Errorf("walkdir recursion Error: %v", err)
			}
		} else {
			if err := handler(entryPath, entry); err != nil {
				return fmt.Errorf("walkdir handler Error: %v", err)
			}
		}
	}
	return nil
}

// Refresh runs and checks for any new or removed files in the root directory
func (m *Monitor) Refresh() {
	err := WalkDir(m.root, func(path string, d fs.DirEntry) error {
		sum, _ := Checksum(path)
		if !d.IsDir() && !m.Exists(path) {
			m.Add(path, &File{
				path: path,
				last: sum,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println("fs update Error:", err)
		return
	}
}
