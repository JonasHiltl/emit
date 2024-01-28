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
	maxTimeout   time.Duration
	maxBatchSize int
}

func WithTimeout(maxTimeout time.Duration) func(*OnOptions) {
	return func(s *OnOptions) {
		s.maxTimeout = maxTimeout
	}
}

func WithBatchSize(maxBatchSize int) func(*OnOptions) {
	return func(s *OnOptions) {
		s.maxBatchSize = maxBatchSize
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

	cn := make(chan any)
	e.subs[key[T]{}] = append(e.subs[key[T]{}], cn)

	go func() {
		for {
			var batch []T

			var expire <-chan time.Time
			if opts.maxTimeout > 0 {
				expire = time.After(opts.maxTimeout)
			}

			for {
				select {
				case value, ok := <-cn:
					if !ok {
						goto done
					}

					if v, ok := value.(T); ok {
						batch = append(batch, v)
					}
					if opts.maxBatchSize == 0 || len(batch) == opts.maxBatchSize {
						goto done
					}
				case <-expire:
					goto done
				}
			}
		done:
			if len(batch) > 0 {
				fn(batch)
			}
		}
	}()
}
