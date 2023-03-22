# wfile
A simple file watcher for Go.

`wfile` compares MD5 checksums of a given directory for changes, and emits a change event when a change occurs. 
## Listening for changes
```go
package main

import "github.com/aboxofsox/wfile"

func main(){
    wfile.Listen("root")
}
```

## Subscribing to events
Events are triggered when a change is detected. Currently, there are only three different events; `CHANGE`, `NOCHANGE`, `ERROR`. Errors live in their own channel.

To "subscribe" to these events:
```go
package main

import "github.com/aboxofsox/wfile"

func main() {
    w := &wfile.Watcher{
        interval: time.Millisecond * 500,
        files: make(chan wfile.File),
        events: make(chan wfile.Event),
        errors :make(chan error),
        ffs: wfile.NewFS("root"),
    }

    go w.Watch()
    go w.Subscribe()

    for {
        go func(){
            select {
                case event := <-w.events:
                    // handle event
                case err := <-w.errors:
                    // handle error
            }
        }()
    }


}
```

