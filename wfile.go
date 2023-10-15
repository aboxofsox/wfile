package wfile

import (
	"context"
	"sync"
	"time"
)

type Listener struct {
	ctx     context.Context
	root    string
	handler func(e Event)
	wg      *sync.WaitGroup
	watcher *Watcher
	ticker  *time.Ticker
}

func Newlistener(ctx context.Context, root string, handler func(e Event)) *Listener {
	return &Listener{
		ctx:     ctx,
		root:    root,
		handler: handler,
		wg:      new(sync.WaitGroup),
		watcher: &Watcher{
			events:  make(chan Event),
			monitor: newMonitor(root),
		},
		ticker: time.NewTicker(500 * time.Millisecond),
	}
}

func (l *Listener) Watch() {
	defer l.ticker.Stop()

	for {
		go l.handleEvent()
		l.wg.Wait()
	}
}

func (l *Listener) handleEvent() {
	l.wg.Add(1)
	defer l.wg.Done()

	select {
	case <-l.ctx.Done():
		return
	case <-l.ticker.C:
		l.watchAndSubscribe()
	}
}

func (l *Listener) watchAndSubscribe() {
	defer l.wg.Done()
	l.watcher.watch(l.ctx)
	l.watcher.subscribe(l.ctx, l.handler)
}

// Listen starts monitoring the directory at the specified root Path for changes at the specified interval.
// When a change is detected, the handler function is called with the details of the event.
// Listening is terminated when ctx.Done() is triggered.
//func Listen(ctx context.Context, root string, handler func(e Event)) {
//	wg := new(sync.WaitGroup)
//
//	watcher := &Watcher{
//		events:  make(chan Event),
//		monitor: newMonitor(root),
//	}
//
//	ticker := time.NewTicker(500 * time.Millisecond)
//	defer ticker.Stop()
//
//	for {
//		wg.Add(1)
//		select {
//		case <-ctx.Done():
//			return
//		case <-ticker.C:
//			go func() {
//				defer wg.Done()
//				watcher.watch(ctx)
//			}()
//			go func() {
//				defer wg.Done()
//				watcher.subscribe(ctx, handler)
//			}()
//		}
//		wg.Wait()
//	}
//}
