# wfile
A simple file watcher for Go.

`wfile` uses MD5 hashes to detect changes to a file's contents. While `wfile` can be useful for smaller projects, it might not work out so well if you have a lot of files you to need to listen to.

## Usage
```go
package main

import (
    "github.com/aboxofsox/wfile"
    "context"
 )

func main() {
    ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
    defer cancel()

    wfile.Listen( ctx, "some-dir",func(e wfile.Event) {
        if e.Code == wfile.CHANGE {
            fmt.Println("change detected")
        }
        if e.Code == wfile.ERROR {
            fmt.Println(e.Error.Error())
        }
    })
}
```
If you want to listen for changes indefinitely:
```go
package main

import (
    "github.com/aboxofsox/wfile"
    "context"
 )

func main() {
    wfile.Listen(context.TODO(), "some-dir", func(e wfile.Event) {
        if e.Code == wfile.CHANGE {
            fmt.Println("change detected")
        }
        if e.Code == wfile.ERROR {
            fmt.Println(e.Error.Error())
        }
    })
}
```

