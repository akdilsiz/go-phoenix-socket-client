package go_phoenix_socket_client

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewTimer(t *testing.T) {
	ctx := context.Background()
	timerCalc := func(_ int) time.Duration {
		return time.Second * 1
	}
	callback := func() {
		//
	}
	_ = NewTimer(ctx, callback, timerCalc)
}

func TestTimer_Tries(t *testing.T) {
	ctx := context.Background()
	timerCalc := func(_ int) time.Duration {
		return time.Second * 1
	}
	callback := func() {}
	tmr := NewTimer(ctx, callback, timerCalc)
	assert.Equal(t, 0, tmr.Tries())
}

func TestTimer_ScheduleTimeout(t *testing.T) {
	callbackEventCount := 0
	ctx := context.Background()
	timerCalc := func(_ int) time.Duration {
		return time.Second * 1
	}
	callback := func() {
		callbackEventCount += 1
	}
	tmr := NewTimer(ctx, callback, timerCalc)
	tmr.ScheduleTimeout()
	time.Sleep(time.Second * 2)
	assert.Equal(t, callbackEventCount, 1)
	assert.Equal(t, tmr.Tries(), 1)
}

func TestTimer_Reset(t *testing.T) {
	callbackEventCount := 0
	ctx := context.Background()
	timerCalc := func(_ int) time.Duration {
		return time.Second * 1
	}
	callback := func() {
		callbackEventCount += 1
	}
	tmr := NewTimer(ctx, callback, timerCalc)
	tmr.ScheduleTimeout()
	time.Sleep(time.Second * 2)
	assert.Equal(t, callbackEventCount, 1)
	assert.Equal(t, tmr.Tries(), 1)

	tmr.Reset()
	assert.Equal(t, tmr.Tries(), 0)
}
