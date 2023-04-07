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
```
or

```go
import "github.com/aboxofsox/wfile"

func main() {
    watcher := &Watcher{
		Events:  make(chan Event),
		Monitor: m,
	}

	done := make(chan bool)
	wg := new(sync.WaitGroup)

	for {
		wg.Add(1)
		go watcher.Watch(done)
		go func(){
            for event := range watcher.Events {
                switch event.code{
                    case CHANGE:
                        // do something
                        break
                    case NOCHANGE:
                        // do something
                        break
                    case ERROR:
                        // handle error
                        break
                }
            }
        }()
		time.Sleep(time.Millisecond * 1600)
		wg.Wait()
	}
}
```
