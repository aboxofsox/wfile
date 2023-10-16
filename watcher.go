package wfile

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Watcher is a struct that represents a file system Watcher. It contains channels for sending events and errors,
// as well as a monitor that tracks the state of the file system.
type Watcher struct {
	events  chan Event // Channel for sending events
	errors  chan error // Channel for sending errors
	monitor *monitor   // monitor that tracks the state of the file system
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

// watch periodically monitors the file system for changes until a signal is received
// on the done channel. It creates a time.Ticker with a 500 millisecond duration and
// enters a loop that repeatedly selects between receiving from the done channel and
// the Ticker's channel. If a signal is received on the done channel, the Ticker is stopped
// and the function returns. Otherwise, it calls the walk method to scan the file system
// for changes. After each cycle of the loop, the function sleeps for 500 milliseconds to prevent
// excessive CPU usage.
func (w *Watcher) watch(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			w.monitor.purge()
			return
		case <-ticker.C:
			w.walk()
		}
		time.Sleep(time.Millisecond * 500)
	}
}

// walk scans the file system for changes since the last time it was called and
// sends events for any modified files. It first updates the file system state by
// refreshing the monitor object associated with the Watcher. It then loops through
// all the files in the monitor's toMap, which is a map of file paths to
// file objects. For each file, it calculates its checksum using the checksum
// function, which returns an Error if the file can't be read. If the file's checksum
// is different from its last known checksum, it updates the file object and sends
// a CHANGE event to the Watcher's events channel. If there is an Error calculating
// the checksum, it sends an ERROR event to the Watcher's errors channel.
func (w *Watcher) walk() {
	w.monitor.refresh()
	w.processFiles()
}

// refershMonitor calls a refresh method on the associated monitor of the Watcher.
func (w *Watcher) refershMonitor() {
	w.monitor.refresh()
}

// fileChanged checks if the content of a file has changed.
// It does so by comparing the current checksum of the file
// to its previous state. Any error occurred during this process will be pushed to the error channel.
func (w *Watcher) fileChanged(f *file) (result bool) {
	result = false
	sum, err := checksum(f.path)
	if err != nil {
		fmt.Println("file checksum Error:", err)
		w.errors <- err
	}
	if f.last != sum {
		result = true
	}
	return
}

// exists checks if the given file exists.
// If not, the method removes the file path from the monitor and reports an error through the error channel.
func (w *Watcher) exists(f *file) {
	if _, err := os.Stat(f.path); os.IsNotExist(err) {
		w.monitor.delete(f.path)
		w.errors <- err
	}
}

// update modifies the last known checksum of the file and emits a change event.
// Any errors occurred during the checksum calculation process are ignored and will not affect the operation.
func (w *Watcher) update(f *file) {
	sum, _ := checksum(f.path)
	f.last = sum
	w.events <- Event{Name: "change", Path: f.path, Code: CHANGE, Error: nil}
}

// process performs an existence check and
// updates the file if a change in its content has been detected.
func (w *Watcher) process(f *file) {
	w.exists(f)
	if w.fileChanged(f) {
		w.update(f)
	}
}

// processFiles goes through every currently tracked file
// and performs an existence and update check for each one.
func (w *Watcher) processFiles() {
	for _, f := range w.monitor.toMap() {
		w.process(f.(*file))
	}
}

// subscribe creates a subscription for event listeners using a handler function.
// The handler will receive events from the Watcher until context cancellation occurs.
func (w *Watcher) subscribe(ctx context.Context, handler func(e Event)) {
	for event := range w.events {
		select {
		case <-ctx.Done():
			return
		default:
			handler(event)
		}
	}
}
