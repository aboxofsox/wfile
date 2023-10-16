package wfile

import (
	"context"
	"sync"
	"time"
)

type Listener struct {
	Cancel  context.CancelFunc
	ctx     context.Context
	wg      *sync.WaitGroup
	ticker  *time.Ticker
	watcher *Watcher
	handler HandlerFunc
}

type HandlerFunc func(e Event)

// NewListener creates a new Listener
func NewListener(ctx context.Context, root string, handler HandlerFunc) *Listener {
	_ctx, cancel := context.WithCancel(ctx)
	return &Listener{
		Cancel: cancel,
		ctx:    _ctx,
		wg:     new(sync.WaitGroup),
		ticker: time.NewTicker(500 * time.Millisecond), // 🧯💨🔥
		watcher: &Watcher{
			events:  make(chan Event),
			errors:  make(chan error),
			monitor: newMonitor(root),
		},
		handler: handler,
	}
}

// Watch is a method for the Listener struct. It starts a loop that continuously checks for signals from
// the context and ticker attached to the listener.
// If the context is done, it returns, effectively stopping the loop.
// If the ticker sends a signal, it increments a WaitGroup counter and starts two goroutines.
// The first goroutine checks if there isn't an error in the context and if so, reduces the
// WaitGroup counter and invokes the watch function of the watcher attached to the listener.
// The second routine does a similar context error check and, if no error is found, it invokes
// the subscribe function of the watcher by passing the context and an event handler.
// This way, the Watch function provides periodic checks and calls for watcher's functions,
// while providing a mechanism for cleanly stopping the function via context.
func (l *Listener) Watch() {
	defer l.ticker.Stop()

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-l.ticker.C:
			l.wg.Add(1)

			// there can only be one call to l.wg.Done()
			go func() {
				go func() {
					if l.ctx.Err() == nil {
						defer l.wg.Done()
						l.watcher.watch(l.ctx)
					}
				}()
				if l.ctx.Err() == nil {
					l.watcher.subscribe(l.ctx, l.handler)
				}
			}()
		}
	}
}
