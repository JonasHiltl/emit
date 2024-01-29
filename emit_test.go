package emit

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRegisterSub(t *testing.T) {
	type event string

	e := &Emitter{}
	On(e, func(t []event) {}, WithBatchSize(10), WithTimeout(10*time.Second))

	assert.Len(t, e.subs, 1)
	assert.Len(t, e.subs[key[event]{}], 1)
}

func TestBatch(t *testing.T) {
	ch := make(chan any, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	var count atomic.Int64
	go batch[int](ch, 5, 0, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(10), count.Load())

	ch = make(chan any, 10)
	for i := 0; i < 10; i++ {
		ch <- i
	}

	count.Store(0)
	i := 0
	go batch[int](ch, 3, 0, func(batch []int) {
		if i < 3 && len(batch) != 3 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		if i == 3 && len(batch) != 1 {
			assert.Fail(t, "expected batch size 1, got %d", len(batch))
		}
		i++
		count.Add(int64(len(batch)))
	})
	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(10), count.Load())
}
