# wfile
A simple file watcher for Go.

`wfile` uses MD5 hashes to detect changes to a file's contents. While `wfile` can be useful for smaller projects, it might not work out so well if you have a lot of files you to need to listen to.

## Usage
```go
package main

import (
    "github.com/aboxofsox/wfile"
    "context"
	"fmt"
	"time"
 )

func main() {
	listener := wfile.NewListener(context.Background(), "root", func(e wfile.Event) {
		switch e.Code {
		case wfile.CHANGE:
			fmt.Println("change detected")
		case wfile.ERROR:
			fmt.Println(e.Error.Error())
        }
    })
	
	go func (){
		time.AfterFunc(time.Second * 30, func() {
			listener.Cancel()
        })
    }()
	
	listener.Watch()
	
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
	listener := wfile.NewListener(context.Background(), "root", func(e wfile.Event) {
		switch e.Code {
		case wfile.CHANGE:
			fmt.Println("change detected")
		case wfile.ERROR:
			fmt.Println(e.Error.Error())
		}
	})
	
	listener.Watch()
}
```

