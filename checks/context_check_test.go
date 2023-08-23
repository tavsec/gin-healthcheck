package checks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestContextCheck(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)

	c := NewContextCheck(ctx)

	assert.NotNil(t, c)
	assert.Equal(t, "github.com/tavsec/gin-healthcheck/checks.TestContextCheck", c.Name())
	assert.True(t, c.Pass())

	cancel()
	<-ctx.Done()

	// We need to give time to the goroutine to get scheduled before checking the status
	time.Sleep(1 * time.Millisecond)

	assert.False(t, c.Pass())
}

func TestContextCheckWithName(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	assert.NotNil(t, ctx)
	assert.NotNil(t, cancel)

	c := NewContextCheck(ctx, "test")

	assert.NotNil(t, c)
	assert.Equal(t, "test", c.Name())
	assert.True(t, c.Pass())

	cancel()
	<-ctx.Done()

	// We need to give time to the goroutine to get scheduled before checking the status
	time.Sleep(1 * time.Millisecond)

	assert.False(t, c.Pass())
}

func TestWrongContext(t *testing.T) {
	assertPanic(t, func() {
		NewContextCheck(nil)
	})
}

func TestWrongName(t *testing.T) {
	assertPanic(t, func() {
		NewContextCheck(context.Background(), "test", "test2")
	})
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}
