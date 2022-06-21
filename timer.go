package gophoenixsocketclient

import (
	"context"
	"time"
)

// Timer scheduling system
type Timer interface {
	Tries() int
	Reset()
	ScheduleTimeout()
}

type timer struct {
	ctx            context.Context
	timerCtx       context.Context
	timerCtxCancel context.CancelFunc
	callback       func()
	timerCalc      TimerCalc
	timer          *time.Timer
	tries          int
}

// NewTimer initialize timer
func NewTimer(ctx context.Context, callback func(), timerCalc TimerCalc) Timer {
	t := new(timer)
	t.ctx = ctx
	t.callback = callback
	t.timerCalc = timerCalc

	return t
}

// Tries how many times it has been scheduled
func (t *timer) Tries() int {
	return t.tries
}

// Reset resetting to the timer
func (t *timer) Reset() {
	t.tries = 0
	if t.timer != nil {
		t.timerCtxCancel()
		t.timer.Stop()
	}
}

// ScheduleTimeout timer is applied repeatedly with the number of tries
func (t *timer) ScheduleTimeout() {
	if t.timer != nil {
		t.timerCtxCancel()
		t.timer.Stop()
	}
	t.timerCtx, t.timerCtxCancel = context.WithCancel(t.ctx)
	t.timer = time.NewTimer(t.timerCalc(t.tries + 1))
	go t.run()
}

func (t *timer) run() {
	select {
	case <-t.ctx.Done():
		return
	case <-t.timer.C:
		t.tries++
		t.callback()
	}
}
