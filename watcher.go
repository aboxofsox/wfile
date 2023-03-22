package wfile

import (
	"fmt"
	"time"
)

// Watcher holds all relevant data for the watching mechanism.
type Watcher struct {
	interval time.Duration
	files    chan File
	events   chan Event
	errors   chan error
	ffs      *FS
}

// Watch() will "watch" a given directory for changes.
func (w *Watcher) Watch() {
	ticker := time.NewTicker(w.interval)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				w.Walk()
			}
		}
	}()

	time.Sleep(1600 * time.Millisecond)
	ticker.Stop()
	done <- true
}

// Walk() reads files in a given directory, does the checksum comparisons,
// and sends events to the Event channel of Watcher.
func (w *Watcher) Walk() {
	w.ffs.Update()

	for _, f := range w.ffs.Iter() {
		file := f.(*File)
		sum, err := Checksum(file.path)
		if err != nil {
			fmt.Println("file checksum error:", err)
			w.files <- *file
			w.events <- Event{name: "error", path: file.path, code: ERROR, error: err}
			w.errors <- err
		}
		if file.last != sum {
			file.last = sum
			w.events <- Event{name: "change", path: file.path, code: CHANGE, error: nil}
		}
	}
}

// Subscribe() will listen for events emitted from the Watcher.
func (w *Watcher) Subscribe() {
	for event := range w.events {
		switch event.code {
		case CHANGE:
			fmt.Println("change detected in:", event.path)
		case NOCHANGE:
			fmt.Println("no change")
		case ERROR:
			fmt.Println("an error occurred:", event.path, event.error)
		}
	}
}
