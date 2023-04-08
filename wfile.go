package wfile

import (
	"sync"
	"time"
)

// Listen starts monitoring the directory at the specified root Path for changes at the specified interval.
// When a change is detected, the handler function is called with the details of the event.
func Listen(root string, interval time.Duration, handler EventHandler) {
	watcher := &Watcher{
		Events:  make(chan Event),
		Monitor: NewMonitor(root),
	}

	done := make(chan bool)
	wg := new(sync.WaitGroup)

	for {
		wg.Add(1)
		go watcher.Watch(done)
		go watcher.Subscribe(handler)
		time.Sleep(interval)
		wg.Wait()
	}
}
