# wfile
A simple Go package for file watching.

## Usage
`wfile` listens to a given directory for any changes to it's content by generating MD5 hashes and comparing them.

```go
package main

import "github.com/aboxofsox/wfile"

func main(){
	wfile.Listen()
}
```

`wfile.Listen()` will spawn goroutines to do a few things:
- walk through the given directory and update a "cache" of files with new or removed files.
- generate an MD5 hash, store it, walk through the directory again, hash again, and compare the two hashes.
- if a difference is detected, a `CHANGE` event is triggered.
### Define a Watcher
You can also define your own `Watcher` with your own parameters:
```go
package main

import "github.com/aboxofsox/wfile"

func main() {
	watcher := &wfile.Watcher{
		ffs: wfile.NewFS("root_dir"),
		interval: time.Millisecond * 100,
		files: make(chan File),
		events: make(chan Event),
    }
	
	for {
		go watcher.Watch()
		go watcher.Subscribe()
    }
}
```