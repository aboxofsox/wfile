package wfile

import (
	"fmt"
	"time"
)

// Watcher holds all relevant data for the watching mechanism.
type Watcher struct {
	events  chan Event
	errors  chan error
	monitor *Monitor
}

type EventCode int

const (
	CHANGE EventCode = iota
	NOCHANGE
	ERROR
)

// Event represents the event data when changes occur.
type Event struct {
	name  string
	code  EventCode
	path  string
	error error
}

// Watch will "watch" a given directory for changes.
func (w *Watcher) Watch(done chan bool) {
	ticker := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			w.Walk()
		}
		time.Sleep(time.Millisecond * 1600)
	}

}

// Walk reads files in a given directory, does the checksum comparisons,
// and sends events to the Event channel of Watcher.
func (w *Watcher) Walk() {
	w.monitor.Refresh()

	for _, f := range w.monitor.ExportFileMap() {
		file := f.(*File)
		sum, err := Checksum(file.path)
		if err != nil {
			fmt.Println("file checksum error:", err)
			w.errors <- err
		}
		if file.last != sum {
			file.last = sum
			w.events <- Event{name: "change", path: file.path, code: CHANGE, error: nil}
		}
	}
}

// Subscribe will listen for events emitted from the Watcher.
func (w *Watcher) Subscribe() {
	for event := range w.events {
		switch event.code {
		case CHANGE:
			fmt.Println("change detected in:", event.path)
			break
		case NOCHANGE:
			fmt.Println("no change")
			break
		case ERROR:
			fmt.Println("an error occurred:", event.path, event.error)
			break
		}
	}
}
