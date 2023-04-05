# wfile
A simple file watcher for Go.

`wfile` compares MD5 checksums of a given directory for changes, and emits a change event when a change occurs.

## Subscribing to events
Events are triggered when a change is detected. Currently, there are only three different events; `CHANGE`, `NOCHANGE`, `ERROR`. Errors live in their own channel.

To "subscribe" to these events:
```go
package main

import "github.com/aboxofsox/wfile"

func main() {
    watcher := &Watcher{
		interval: time.Millisecond * 500,
		events:   make(chan Event),
		monitor:  m,
	}

	done := make(chan bool)
	wg := new(sync.WaitGroup)

	for {
		wg.Add(1)
		go watcher.Watch(done, wg

        // listen for any change events
		go func() {
			defer wg.Done()
			for event := range watcher.events {
				wg.Add(1)
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
```


