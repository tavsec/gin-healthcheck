package checks

import (
	"context"
	"runtime"
	"sync/atomic"
)

type contextCheck struct {
	name       string
	terminated uint32 // TODO: When the minimal supported base go version is 1.19, use atomic.Bool
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
		atomic.StoreUint32(&c.terminated, 1)
	}()

	return &c
}

func (c *contextCheck) Pass() bool {
	v := atomic.LoadUint32(&c.terminated)
	return v == 0
}

func (c *contextCheck) Name() string {
	return c.name
}
