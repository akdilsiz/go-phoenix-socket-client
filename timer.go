package gophoenixsocketclient

import (
	"context"
	"time"
)

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

func NewTimer(ctx context.Context, callback func(), timerCalc TimerCalc) Timer {
	t := new(timer)
	t.ctx = ctx
	t.callback = callback
	t.timerCalc = timerCalc

	return t
}

func (t *timer) Tries() int {
	return t.tries
}

func (t *timer) Reset() {
	t.tries = 0
	if t.timer != nil {
		t.timerCtxCancel()
		t.timer.Stop()
	}
}

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
