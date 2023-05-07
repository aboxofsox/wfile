package wfile

import (
	"context"
	"sync"
	"time"
)

// Listen starts monitoring the directory at the specified root Path for changes at the specified interval.
// When a change is detected, the handler function is called with the details of the event.
// Listening is terminated when ctx.Done() is triggered.
func Listen(root string, ctx context.Context, handler func(e Event)) {
	wg := new(sync.WaitGroup)

	watcher := &watcher{
		events:  make(chan Event),
		monitor: newMonitor(root),
	}

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		wg.Add(1)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			go func() {
				defer wg.Done()
				watcher.watch(ctx)
			}()
			go func() {
				defer wg.Done()
				watcher.subscribe(ctx, handler)
			}()
		}
		wg.Wait()
	}
}
