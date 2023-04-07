package wfile

import (
	"sync"
	"time"
)

// Listen will start the file listening process.
// the default polling interval is 500ms.
func Listen(m *Monitor) {
	watcher := &Watcher{
		Events:  make(chan Event),
		Monitor: m,
	}

	done := make(chan bool)
	wg := new(sync.WaitGroup)

	for {
		wg.Add(1)
		go watcher.Watch(done)
		go watcher.Subscribe()
		time.Sleep(time.Millisecond * 1600)
		wg.Wait()
	}
}
