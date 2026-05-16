package gophoenixsocketclient

import (
	"net/url"
)

type Transport interface {
	EndpointURL() *url.URL
	ReplaceTransport(newTransport Transport)
	Disconnect(callback func(), code int, reason string)
	Teardown(callback func(), code int, reason string)
	OnOpen()
	SetOnOpen(callback func()) string
	OnError(err error)
	SetOnError(callback func(err error, transportBefore Transport, establishedBefore int)) string
	OnMessage(event interface{})
	SetOnMessage(callback func(event interface{})) string
	OnClose(event interface{})
	SetOnClose(callback func(event interface{})) string
}
