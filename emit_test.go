package emit

import (
	"fmt"
	"testing"
	"time"
)

func TestRegisterSub(t *testing.T) {
	type event string

	e := &Emitter{}
	On(e, func(t []event) {
		fmt.Println(t)
	}, WithBatchSize(10), WithTimeout(10*time.Second))

	if len(e.subs) != 1 || len(e.subs[key[event]{}]) != 1 {
		t.Fatal("Failed to register subscriber")
	}
}
