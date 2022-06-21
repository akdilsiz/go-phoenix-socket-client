package gophoenixsocketclient

import (
	"context"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"net/url"
	"time"
)

type TimerCalc func(tries int) time.Duration

var (
	rejoinAfterMS = []time.Duration{
		time.Second * 1,
		time.Second * 2,
		time.Second * 5,
	}
	reconnectAfterMS = []time.Duration{
		time.Millisecond * 10,
		time.Millisecond * 50,
		time.Millisecond * 100,
		time.Millisecond * 150,
		time.Millisecond * 200,
		time.Millisecond * 250,
		time.Millisecond * 500,
		time.Millisecond * 1000,
		time.Millisecond * 2000,
	}
)

type SocketOptions struct {
	Timeout             time.Duration
	Transport           constants.Transport
	Encode              Serializer
	Decode              Serializer
	HeartbeatIntervalMS time.Duration
	RejoinAfterMS       TimerCalc
	ReconnectAfterMS    TimerCalc
	Logger              Logger
	LongpollerTimeout   time.Duration
	Params              url.Values
	VSN                 string
}

func DefaultSocketOptions() *SocketOptions {
	so := new(SocketOptions)
	so.Timeout = constants.DefaultTimeout
	so.Transport = constants.Websocket
	so.Encode = NewSerializer()
	so.Decode = NewSerializer()
	so.HeartbeatIntervalMS = time.Second * 30
	so.RejoinAfterMS = func(tries int) time.Duration {
		if len(rejoinAfterMS) < tries {
			return rejoinAfterMS[tries-1]
		}
		return time.Second * 10
	}
	so.ReconnectAfterMS = func(tries int) time.Duration {
		if len(reconnectAfterMS) < tries {
			return reconnectAfterMS[tries-1]
		}
		return time.Millisecond * 5000
	}
	so.Logger = NewLogger(context.Background(), constants.DefaultVsn)
	so.Params = url.Values{}

	return so
}

func makeOptions(options ...*SocketOptions) (opts *SocketOptions) {
	if len(options) > 0 {
		opts = options[0]
		return opts
	}
	opts = DefaultSocketOptions()
	return opts
}
