package wfile

import (
	"sync"
	"time"
)

// Listen will start the file listening process.
// the default polling interval is 1600ms.
func Listen(root string, interval time.Duration, handler EventHandler) {
	m := NewMonitor(root)

	watcher := &Watcher{
		Events:  make(chan Event),
		Monitor: m,
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
