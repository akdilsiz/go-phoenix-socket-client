package gophoenixsocketclient

import (
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"time"
)

// Push ..
type Push interface {
	Resend(timeout time.Duration)
	Reset()
	Send()
}

type push struct {
	channel      string
	event        constants.ChannelEvent
	payload      []byte
	receivedResp *Received
	timeout      time.Duration
	timeoutTimer time.Ticker
	recHooks     []Hook
	sent         bool
}

// NewPush ..
func NewPush(channel string,
	event constants.ChannelEvent,
	payload []byte,
	timeout time.Duration) Push {
	p := new(push)
	p.channel = channel
	p.event = event
	p.payload = payload
	p.receivedResp = nil
	p.timeout = timeout
	p.recHooks = make([]Hook, 0)
	p.sent = false

	return p
}

// Resend ..
func (p *push) Resend(timeout time.Duration) {
	p.timeout = timeout
	p.Reset()
	p.Send()
}

// Send ..
func (p *push) Send() {

}

// Reset ..
func (p *push) Reset() {

}

func (p *push) hasReceived(status string) bool {
	if p.receivedResp != nil && p.receivedResp.Status == status {
		return true
	}

	return false
}
