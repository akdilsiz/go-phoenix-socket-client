package gophoenixsocketclient

import (
	"context"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"net/url"
	"strings"
	"sync"
)

// Longpoll phoenix longpoll transport interface
type Longpoll struct {
	ctx           context.Context
	mut           sync.Mutex
	endpoint      string
	token         string
	skipHeartbeat bool
	pollEndpoint  string
	readyState    constants.SocketState
}

// NewLongpoll initialize phoenix longpoll connector
func NewLongpoll(ctx context.Context, endpoint string) (Transport, error) {
	lp := new(Longpoll)
	lp.ctx = ctx
	lp.endpoint = endpoint
	lp.pollEndpoint = lp.normalizeEndpoint()
	return lp, nil
}

func (lp *Longpoll) normalizeEndpoint() string {
	e := strings.ReplaceAll(lp.endpoint, "ws://", "http://")
	e = strings.ReplaceAll(e, "wss://", "https://")
	e = strings.ReplaceAll(e, string(constants.Websocket)+"/", string(constants.Longpoll)+"/")
	return e
}

// EndpointURL ..
func (lp *Longpoll) EndpointURL() *url.URL {
	u, _ := url.Parse(lp.pollEndpoint)
	u.Query().Set("token", lp.token)

	return u
}

func (lp *Longpoll) ReplaceTransport(newTransport Transport) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) Disconnect(callback func(), code int, reason string) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) Teardown(callback func(), code int, reason string) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) OnOpen() {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) SetOnOpen(callback func()) string {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) OnError(err error) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) SetOnError(callback func(err error, transportBefore Transport, establishedBefore int)) string {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) OnMessage(event interface{}) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) SetOnMessage(callback func(event interface{})) string {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) OnClose(event interface{}) {
	//TODO implement me
	panic("implement me")
}

func (lp *Longpoll) SetOnClose(callback func(event interface{})) string {
	//TODO implement me
	panic("implement me")
}
