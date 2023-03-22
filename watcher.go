package wfile

import (
	"fmt"
	"time"
)

type Watcher struct {
	interval time.Duration
	files    chan File
	events   chan Event
	errors   chan error
	ffs      *FS
}

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

func (w *Watcher) Close() {
}
