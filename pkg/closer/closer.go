package closer

import (
	"context"
	"fmt"
	"log"
	"sync"
)

type Callback func(ctx context.Context) error
type EasyCallback func()

type Closer struct {
	mutex     sync.Mutex
	callbacks []Callback
}

func New() *Closer {
	return &Closer{}
}

func (receiver *Closer) Add(callback Callback) {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	receiver.callbacks = append(receiver.callbacks, callback)
}

func (receiver *Closer) Close(ctx context.Context) error {
	receiver.mutex.Lock()
	defer receiver.mutex.Unlock()

	complete := make(chan struct{}, 1)

	go func() {
		for index, callback := range receiver.callbacks {
			if err := callback(ctx); err != nil {
				log.Printf("[! (%d)] %v", index, err)
			}
		}

		complete <- struct{}{}
	}()

	select {
	case <-complete:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	}
}

func Wrap(callback EasyCallback) Callback {
	return func(_ context.Context) error {
		callback()

		return nil
	}
}
