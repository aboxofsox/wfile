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
        Events: make(chan wfile.Event),
        Monitor: wfile.NewMonitor("."),
    }

    done := make(chan bool)
    wg := new(sync.WaitGroup)

    for {
        wg.Add(1)
        go watcher.Watch(done) // <- watch for changes
        go watcher.Subscribe(handler) // <- subscribe to events
        time.Sleep(time.Millisecond * 1600)
        wg.Wait()
    }
}

func handler(event wfile.Event) {
    switch event.Code {
        case wfile.CHANGE:
            // do something
        case wfile.NOCHANGE:
            // do something else
        case wfile.ERROR:
            // handle error
    }
}
```
Alternatively:
```go
package main

import "github.com/aboxofsox/wfile"

func main() {
    wfile.Listen(".", time.Millisecond * 1600, func(e wfile.Event){
        switch e.code {
            case wfile.CHANGE:
                // do soemthing
            case wfile.NOCHANGE:
                // do something
            case wfile.ERROR:
                // handle error
        }
    })
}
```