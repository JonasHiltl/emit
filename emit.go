package emit

import (
	"errors"
	"sync"
	"time"
)

type key[T any] struct{}

type Emitter struct {
	subs map[any][]chan any
	mu   sync.RWMutex
}

func (e *Emitter) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, chans := range e.subs {
		for _, ch := range chans {
			close(ch)
		}
	}
}

func (e *Emitter) init() {
	if e.subs == nil {
		e.subs = make(map[any][]chan any)
	}
}

func Emit[T any](e *Emitter, v T) error {
	if e == nil {
		return errors.New("emitter is nil")
	}
	e.mu.RLock()
	defer e.mu.RUnlock()

	subs, ok := e.subs[key[T]{}]
	if ok {
		for _, cn := range subs {
			select {
			case cn <- v:
				// message sent
			default:
				// message dropped
			}
		}
	}

	return nil
}

type OnOptions struct {
	timeout   time.Duration
	batchSize int
}

func WithTimeout(maxTimeout time.Duration) func(*OnOptions) {
	return func(s *OnOptions) {
		s.timeout = maxTimeout
	}
}

func WithBatchSize(maxBatchSize int) func(*OnOptions) {
	return func(s *OnOptions) {
		s.batchSize = maxBatchSize
	}
}

func On[T any](e *Emitter, fn func([]T), options ...func(*OnOptions)) {
	opts := &OnOptions{}
	for _, o := range options {
		o(opts)
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	e.init()

	if _, ok := e.subs[key[T]{}]; !ok {
		e.subs[key[T]{}] = make([]chan any, 0, 3) // default capacity of 3
	}

	ch := make(chan any)
	e.subs[key[T]{}] = append(e.subs[key[T]{}], ch)

	go collect(ch, opts.batchSize, opts.timeout, fn)
}

// collect reads from a channel and calls fn when batchSize is met or timeout is triggered
func collect[T any](ch <-chan any, batchSize int, timeout time.Duration, fn func([]T)) {
	var ticker *time.Ticker
	if timeout > 0 {
		ticker = time.NewTicker(timeout)
	} else {
		ticker = new(time.Ticker)
	}
	defer ticker.Stop()

	for batchSize <= 1 { // sanity check,
		for v := range ch {
			if v, ok := v.(T); ok {
				fn([]T{v})
			}
		}

		return
	}

	// batchSize > 1
	batch := make([]T, 0, batchSize)
	for {
		select {
		case <-ticker.C:
			if len(batch) > 0 {
				fn(batch)
				// reset
				batch = make([]T, 0, batchSize)
				if timeout > 0 {
					ticker.Reset(timeout)
				}
			}
		case v, ok := <-ch:
			if !ok { // closed
				return
			}
			if v, ok := v.(T); ok {
				batch = append(batch, v)
			}
			if len(batch) == batchSize { // full
				if len(batch) > 0 {
					fn(batch)
					batch = make([]T, 0, batchSize)
					if timeout > 0 {
						ticker.Reset(timeout)
					}
				}
			}
		}
	}
}
