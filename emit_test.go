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

func TestCollect(t *testing.T) {
	ch := make(chan any, 10)
	var count atomic.Int64
	go collect[int](ch, 5, 0, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		count.Add(int64(len(batch)))
	})

	for i := 0; i < 10; i++ {
		ch <- i
	}

	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(10), count.Load())
}

func TestCollectBuffer(t *testing.T) {
	ch := make(chan any, 10)
	var count atomic.Int64
	go collect[int](ch, 5, 0, func(batch []int) {
		if len(batch) != 5 {
			assert.Fail(t, "expected batch size 5, got %d", len(batch))
		}
		count.Add(int64(len(batch)))
	})

	ch <- 0

	time.Sleep(time.Second)
	close(ch)
	assert.Equal(t, int64(0), count.Load())
}
