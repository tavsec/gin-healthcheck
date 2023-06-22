package checks

import (
	"context"
	"runtime"
	"sync/atomic"
)

type contextCheck struct {
	name       string
	terminated atomic.Int32 // TODO: When the minimal supported base go version is 1.19, use atomic.Bool
	ctx        context.Context
}

func NewContextCheck(ctx context.Context, name ...string) Check {
	if len(name) > 1 {
		panic("context check does only accept one name")
	}
	if ctx == nil {
		panic("context check needs a context")
	}

	contextName := "Unknown"
	if len(name) == 1 {
		contextName = name[0]
	} else {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			contextName = details.Name()
		}
	}

	c := contextCheck{
		name: contextName,
		ctx:  ctx,
	}

	go func() {
		<-ctx.Done()
		c.terminated.Store(1)
	}()

	return &c
}

func (c *contextCheck) Pass() bool {
	return c.terminated.Load() == 0
}

func (c *contextCheck) Name() string {
	return c.name
}
