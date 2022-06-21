package go_phoenix_socket_client

import (
	"context"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"net/url"
	"strings"
	"sync"
)

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

func (lp *longpoll) EndpointURL() (*url.URL, error) {
	u, err := url.Parse(lp.pollEndpoint)
	if err != nil {
		return nil, err
	}
	u.Query().Set("token", lp.token)

	return u, nil
}

func (lp *longpoll) SetOnClose(fn func(event CloseEvent)) {
	lp.mut.Lock()
	lp.onClose = fn
	lp.mut.Unlock()
}

func (lp *longpoll) SetOnError(fn func(err string)) {
	lp.mut.Lock()
	lp.onError = fn
	lp.mut.Unlock()
}
