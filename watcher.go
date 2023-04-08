package wfile

import (
	"fmt"
	"time"
)

// Watcher is a struct that represents a file system watcher. It contains channels for sending events and errors,
// as well as a Monitor that tracks the state of the file system.
type Watcher struct {
	Events  chan Event // Channel for sending events
	Errors  chan error // Channel for sending errors
	Monitor *Monitor   // Monitor that tracks the state of the file system
}

type EventCode int

// These constants define the possible event codes that can be sent by a Watcher.
//
// CHANGE indicates that the file system has detected a change in a file that is being monitored.
// NOCHANGE indicates that the file system has not detected any changes since the last check.
// ERROR indicates that an Error has occurred while monitoring the file system.
const (
	CHANGE EventCode = iota
	NOCHANGE
	ERROR
)

// Event represents a file system event, including the Name of the event, the type of event
// (as an EventCode), the Path of the affected file, and an optional Error value.
type Event struct {
	Name  string
	Code  EventCode
	Path  string
	Error error
}

// EventHandler type defines the EventHandler function
// making it easier to pass around as function parameters
type EventHandler func(Event)

// Watch periodically monitors the file system for changes until a signal is received
// on the done channel. It creates a time.Ticker with a 500 millisecond duration and
// enters a loop that repeatedly selects between receiving from the done channel and
// the Ticker's channel. If a signal is received on the done channel, the Ticker is stopped
// and the function returns. Otherwise, it calls the Walk method to scan the file system
// for changes. After each cycle of the loop, the function sleeps for 1.6 seconds to prevent
// excessive CPU usage.
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

// Walk scans the file system for changes since the last time it was called and
// sends events for any modified files. It first updates the file system state by
// refreshing the Monitor object associated with the Watcher. It then loops through
// all the files in the Monitor's ExportFileMap, which is a map of file paths to
// file objects. For each file, it calculates its checksum using the Checksum
// function, which returns an Error if the file can't be read. If the file's checksum
// is different from its last known checksum, it updates the file object and sends
// a CHANGE event to the Watcher's Events channel. If there is an Error calculating
// the checksum, it sends an ERROR event to the Watcher's Errors channel.
func (w *Watcher) Walk() {
	w.Monitor.Refresh()

	for _, f := range w.Monitor.ExportFileMap() {
		file := f.(*File)
		sum, err := Checksum(file.path)
		if err != nil {
			fmt.Println("file checksum Error:", err)
			w.Errors <- err
		}
		if file.last != sum {
			file.last = sum
			w.Events <- Event{Name: "change", Path: file.path, Code: CHANGE, Error: nil}
		}
	}
}

// Subscribe will listen for Events emitted from the Watcher.
func (w *Watcher) Subscribe(handler EventHandler) {
	for event := range w.Events {
		handler(event)
	}
}
