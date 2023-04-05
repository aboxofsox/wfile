package wfile

import (
	"fmt"
	"sync"
	"time"
)

// Listen will start the file listening process.
// the default polling interval is 500ms.
func Listen(m *Monitor) {
	watcher := &Watcher{
		events:  make(chan Event),
		monitor: m,
	}

	done := make(chan bool)
	wg := new(sync.WaitGroup)

	for {
		wg.Add(1)
		go watcher.Watch(done)
		go func() {
			//defer wg.Done()
			for event := range watcher.events {
				//wg.Add(1)
				switch event.code {
				case CHANGE:
					fmt.Println("change detected")
					break
				case NOCHANGE:
					break
				case ERROR:
					fmt.Println(event.error)
					break
				}
			}
		}()
		time.Sleep(time.Millisecond * 1600)
		wg.Wait()
	}
}
