package gophoenixsocketclient

import (
	"context"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"net/url"
	"strings"
	"sync"
)

// Longpoll phoenix longpoll transport interface
type Longpoll interface {
	EndpointURL() (*url.URL, error)
	SetOnClose(fn func(event CloseEvent))
	SetOnError(fn func(err string))
}

type longpoll struct {
	ctx           context.Context
	mut           sync.Mutex
	endpoint      string
	token         string
	skipHeartbeat bool
	pollEndpoint  string
	readyState    constants.SocketState
	onClose       func(event CloseEvent)
	onError       func(err string)
}

// NewLongpoll initialize phoenix longpoll connector
func NewLongpoll(ctx context.Context, endpoint string) Longpoll {
	lp := new(longpoll)
	lp.ctx = ctx
	lp.endpoint = endpoint
	lp.pollEndpoint = lp.normalizeEndpoint()
	return lp
}

func (lp *longpoll) normalizeEndpoint() string {
	e := strings.ReplaceAll(lp.endpoint, "ws://", "http://")
	e = strings.ReplaceAll(e, "wss://", "https://")
	e = strings.ReplaceAll(e, string(constants.Websocket)+"/", string(constants.Longpoll)+"/")
	return e
}

// EndpointURL ..
func (lp *longpoll) EndpointURL() (*url.URL, error) {
	u, err := url.Parse(lp.pollEndpoint)
	if err != nil {
		return nil, err
	}
	u.Query().Set("token", lp.token)

	return u, nil
}

// SetOnClose ..
func (lp *longpoll) SetOnClose(fn func(event CloseEvent)) {
	lp.mut.Lock()
	lp.onClose = fn
	lp.mut.Unlock()
}

// SetOnError ..
func (lp *longpoll) SetOnError(fn func(err string)) {
	lp.mut.Lock()
	lp.onError = fn
	lp.mut.Unlock()
}
