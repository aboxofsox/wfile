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

type monitor struct {
	root     string
	interval time.Duration
	files    *sync.Map
}

type file struct {
	path string
	last [md5.Size]byte
}

// newMonitor acts as a cache.
// keeping track of the files to watch.
func newMonitor(root string) *monitor {
	return &monitor{
		root:  root,
		files: &sync.Map{},
	}
}

// add adds an item to the cache
func (m *monitor) add(key, value any) {
	if _, exists := m.files.Load(key); exists {
		fmt.Println(key, "exists")
		return
	}
	m.files.Store(key, value)
}

// delete removes an item from the cache
func (m *monitor) delete(key any) {
	if _, exists := m.files.Load(key); !exists {
		return
	}
	m.files.Delete(key)
}

// toMap creates a map[any]any representation of the cache
func (m *monitor) toMap() map[any]any {
	mp := make(map[any]any)

	m.files.Range(func(key, value any) bool {
		mp[key] = value
		return true
	})
	return mp
}

// exists checks if a key exists in the cache
func (m *monitor) exists(key any) bool {
	_, exists := m.files.Load(key)
	return exists
}

// size returns the item count of the cache
func (m *monitor) size() int {
	return len(m.toMap())
}

// walkDir walks through the root directory, recursively
func walkDir(path string, handler func(path string, d fs.DirEntry) error) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("walkdir read Path Error: %v", err)
	}

	for _, entry := range entries {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			if err := walkDir(entryPath, handler); err != nil {
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

// refresh runs and checks for any new or removed files in the root directory
func (m *monitor) refresh() {
	err := walkDir(m.root, func(path string, d fs.DirEntry) error {
		sum, _ := checksum(path)

		if !d.IsDir() && !m.exists(path) {
			m.add(path, &file{
				path: path,
				last: sum,
			})
		}
		return nil
	})
	if err != nil {
		fmt.Println("monitor update Error:", err)
		return
	}
}

func (m *monitor) purge() {
	for k := range m.toMap() {
		m.delete(k)
	}
}
