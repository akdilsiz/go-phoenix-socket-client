package go_phoenix_socket_client

import (
	"context"
	"fmt"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"sync"
	"time"
)

// Socket ..
type Socket interface {
	Protocol() string
	EndpointURL() string
	ReplaceTransport(newTransport constants.Transport)
	Disconnect(callback func(), code int, reason string)
	Teardown(callback func(), code int, reason string)
	MakeRef() string
}

type socket struct {
	ctx                    context.Context
	mut                    sync.Mutex
	options                *SocketOptions
	stateChangeCallbacks   StateChange
	channels               []interface{}
	sendBuffer             chan []byte
	ref                    int64
	timeout                time.Duration
	transport              constants.Transport
	establishedConnections int64
	defaultEncoder         Serializer
	defaultDecoder         Serializer
	closeWasClean          bool
	connectClock           int
	encode                 Serializer
	decode                 Serializer
	heartbeatIntervalMS    time.Duration
	rejoinAfterMS          func(tries int) time.Duration
	reconnectAfterMS       func(tries int) time.Duration
	logger                 Logger
	longpollerTimeout      time.Duration
	params                 url.Values
	endpoint               string
	vsn                    string
	heartbeatTimer         *time.Ticker
	pendingHeartbeatRef    interface{}
	reconnectTimer         Timer
	wsConn                 *websocket.Conn
	dialer                 *websocket.Dialer
	listenCtx              context.Context
	listenCtxCancel        context.CancelFunc
	onOpen                 func(msg Message)
}

// NewSocket ..
func NewSocket(ctx context.Context, endpoint string, options ...*SocketOptions) Socket {
	s := new(socket)
	s.ctx = ctx
	s.stateChangeCallbacks = StateChange{}
	s.channels = make([]interface{}, 0)
	s.sendBuffer = make(chan []byte)
	s.defaultEncoder = NewSerializer()
	s.defaultDecoder = NewSerializer()
	s.establishedConnections = 0
	s.closeWasClean = false
	s.connectClock = 1
	s.ref = 0
	s.parseOptions(makeOptions(options...))
	s.endpoint = fmt.Sprintf("%s/%s", endpoint, constants.Websocket)
	s.heartbeatTimer = nil
	s.pendingHeartbeatRef = nil
	s.reconnectTimer = NewTimer(s.ctx, func() {
		// TODO: TearDown method implementation
	}, s.reconnectAfterMS)

	return s
}

func (s *socket) parseOptions(options *SocketOptions) {
	s.options = options
	s.timeout = options.Timeout
	s.transport = options.Transport
	if s.transport != constants.Longpoll {
		if options.Encode == nil {
			s.encode = s.defaultEncoder
		} else {
			s.encode = options.Encode
		}
		if options.Decode == nil {
			s.decode = s.defaultDecoder
		} else {
			s.decode = options.Decode
		}
	} else {
		s.encode = s.defaultEncoder
		s.decode = s.defaultDecoder
	}
	s.heartbeatIntervalMS = options.HeartbeatIntervalMS
	s.rejoinAfterMS = options.RejoinAfterMS
	s.reconnectAfterMS = options.ReconnectAfterMS
	s.logger = options.Logger
	s.longpollerTimeout = options.LongpollerTimeout
	s.params = options.Params
	s.vsn = options.VSN
}

func (s *socket) hasLogger() bool {
	return s.logger != nil
}

func (s *socket) MakeRef() string {
	newRef := s.ref + 1
	if newRef == s.ref {
		s.ref = 0
	} else {
		s.ref = newRef
	}

	return strconv.FormatInt(s.ref, 10)
}

// ReplaceTransport ...
func (s *socket) ReplaceTransport(newTransport constants.Transport) {
	s.Disconnect(nil, 0, "")
	s.transport = newTransport
}

// Protocol ...
func (s *socket) Protocol() string {
	matched, _ := regexp.MatchString(`(?s)^https`, s.endpoint)
	if matched {
		return "wss"
	}
	return "ws"
}

// EndpointURL ...
func (s *socket) EndpointURL() string {
	u, err := url.Parse(s.endpoint)
	if err != nil {
		return ""
	}
	for k, v := range s.params {
		for i := 0; i < len(v); i++ {
			if i == 0 {
				u.Query().Set(k, v[i])
			} else {
				u.Query().Add(k, v[i])
			}
		}
	}
	u.Query().Set("vsn", s.vsn)
	uri := u.String()
	if string(uri[0]) != "/" {
		return uri
	}
	if string(uri[1]) == "/" {
		return fmt.Sprintf("%s:%s", s.Protocol(), uri)
	}

	return fmt.Sprintf("%s://%s", s.Protocol(), uri)
}

// Disconnect ..
func (s *socket) Disconnect(callback func(), code int, reason string) {
	s.connectClock++
	s.closeWasClean = true
	s.reconnectTimer.Reset()
	s.listenCtxCancel()
}

// Teardown ...
// TODO: implement me!!
func (s *socket) Teardown(callback func(), code int, reason string) {

}

func (s *socket) Connect() {
	s.connectClock++
	if s.wsConn != nil {
		return
	}
	s.dialer = &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: constants.DefaultTimeout,
		ReadBufferSize:   512,
		WriteBufferSize:  512,
		Jar:              nil,
	}
}
