package wfile

import (
	"context"
	"fmt"
	"os"
	"time"
)

// watcher is a struct that represents a file system watcher. It contains channels for sending events and errors,
// as well as a monitor that tracks the state of the file system.
type watcher struct {
	events  chan Event // Channel for sending events
	errors  chan error // Channel for sending errors
	monitor *monitor   // monitor that tracks the state of the file system
}

type EventCode int

// These constants define the possible event codes that can be sent by a watcher.
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
// for changes. After each cycle of the loop, the function sleeps for 1.6 seconds to prevent
// excessive CPU usage.
func (w *watcher) watch(ctx context.Context) {
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
		time.Sleep(time.Millisecond * 1600)
	}
}

// walk scans the file system for changes since the last time it was called and
// sends events for any modified files. It first updates the file system state by
// refreshing the monitor object associated with the watcher. It then loops through
// all the files in the monitor's toMap, which is a map of file paths to
// file objects. For each file, it calculates its checksum using the checksum
// function, which returns an Error if the file can't be read. If the file's checksum
// is different from its last known checksum, it updates the file object and sends
// a CHANGE event to the watcher's events channel. If there is an Error calculating
// the checksum, it sends an ERROR event to the watcher's errors channel.
func (w *watcher) walk() {
	w.monitor.refresh()

	for _, f := range w.monitor.toMap() {
		file := f.(*file)

		if _, err := os.Stat(file.path); os.IsNotExist(err) {
			w.monitor.delete(file.path)
			break
		}

		sum, err := checksum(file.path)
		if err != nil {
			fmt.Println("file checksum Error:", err)
			w.errors <- err
		}
		if file.last != sum {
			file.last = sum
			w.events <- Event{Name: "change", Path: file.path, Code: CHANGE, Error: nil}
		}
	}
}

// subscribe will listen for events emitted from the watcher.
func (w *watcher) subscribe(ctx context.Context, handler func(e Event)) {
	for event := range w.events {
		select {
		case <-ctx.Done():
			return
		default:
			handler(event)
		}
	}
}
