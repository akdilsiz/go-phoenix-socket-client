package constants

import "time"

const (
	DefaultVsn     = "2.0.0"
	DefaultTimeout = time.Second * 10
	WSCloseNormal  = time.Second * 1
)

type SocketState int

const (
	Connecting SocketState = 0
	Open       SocketState = 1
	Closing    SocketState = 2
	SClosed    SocketState = 3
)

type ChannelState string

const (
	Closed  ChannelState = "closed"
	Errored ChannelState = "errored"
	Joined  ChannelState = "joined"
	Joining ChannelState = "joining"
	Leaving ChannelState = "leaving"
)

type ChannelEvent string

const (
	Close ChannelEvent = "phx_close"
	Error ChannelEvent = "phx_error"
	Join  ChannelEvent = "phx_join"
	Reply ChannelEvent = "phx_reply"
	Leave ChannelEvent = "phx_leave"
)

type Transport string

const (
	Longpoll  Transport = "longpoll"
	Websocket Transport = "websocket"
)
