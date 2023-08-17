package client

import (
	"context"
	"log"
)

type Executor struct {
	ctx context.Context
	C   <-chan struct{}
}

func NewExecutor(ctx context.Context) *Executor {
	return &Executor{ctx: ctx}
}

func (e *Executor) Run() {
	done := make(chan struct{})
	e.C = done

	go func() {
		for {
			select {
			case <-e.ctx.Done():
				close(taskCh)
				return
			}
		}
	}()

	go func() {
		defer func() {
			done <- struct{}{}
		}()
		for objectName := range taskCh {
			if err := bucket.DeleteObject(objectName); err != nil {
				log.Printf("[warn] failed to remove object: %s", objectName)
			} else {
				log.Printf("[info] successful remove fail object: %s", objectName)
			}
		}
	}()
}
