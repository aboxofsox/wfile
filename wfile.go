package wfile

import (
	"fmt"
	"time"
)

func Listen(root string) {
	w := &Watcher{
		interval: time.Millisecond * 500,
		ffs:      NewFS(root),
		files:    make(chan File),
		events:   make(chan Event),
	}

	for {
		go func() {
			for {
				w.Watch()
				select {
				case event := <-w.events:
					fmt.Println(event.name)
				case err := <-w.errors:
					fmt.Println("event error:", err.Error())
				}
			}
		}()
		time.Sleep(time.Millisecond * 1000)
	}

}
