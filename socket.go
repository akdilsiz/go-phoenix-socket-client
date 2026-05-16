package gophoenixsocketclient

import (
	"context"
	"fmt"
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"sync"
	"time"
)

// Socket phoenix websocket transport interface
type Socket struct {
	ctx                    context.Context
	mut                    sync.Mutex
	options                *SocketOptions
	stateChangeCallbacks   StateChange
	channels               []interface{}
	sendBuffer             chan []byte
	ref                    int64
	timeout                time.Duration
	transport              Transport
	establishedConnections int
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

// NewSocket initialize phoenix websocket connector
func NewSocket(ctx context.Context, endpoint string, options ...*SocketOptions) (Transport, error) {
	s := new(Socket)
	s.ctx = ctx

	if _, err := url.Parse(endpoint); err != nil {
		return nil, err
	}

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

	return s, nil
}

func (s *Socket) parseOptions(options *SocketOptions) {
	s.options = options
	s.timeout = options.Timeout
	s.transport = options.Transport
	if reflect.ValueOf(s.transport).Elem() == reflect.ValueOf(Longpoll{}).Elem() {
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

func (s *Socket) hasLogger() bool {
	return s.logger != nil
}

func (s *Socket) makeRef() string {
	s.mut.Lock()
	newRef := s.ref + 1
	if newRef == s.ref {
		s.ref = 0
	} else {
		s.ref = newRef
	}
	s.mut.Lock()

	return strconv.FormatInt(s.ref, 10)
}

func (s *Socket) OnOpen() {
	s.mut.Lock()
	s.closeWasClean = false
	s.establishedConnections++
	//TODO: flush send buffer
	s.reconnectTimer.Reset()
	s.mut.Unlock()
	for i := 0; i < len(s.stateChangeCallbacks.Open); i++ {
		s.stateChangeCallbacks.Open[i].Callback()
	}
}

func (s *Socket) SetOnOpen(callback func()) string {
	ref := s.makeRef()
	s.mut.Lock()
	s.stateChangeCallbacks.Open = append(s.stateChangeCallbacks.Open, struct {
		Ref      string
		Callback func()
	}{
		Ref:      ref,
		Callback: callback,
	})
	s.mut.Unlock()
	return ref
}

func (s *Socket) OnError(err error) {
	s.mut.Lock()
	transportBefore := s.transport
	establishedBefore := s.establishedConnections
	s.mut.Unlock()
	for i := 0; i < len(s.stateChangeCallbacks.Error); i++ {
		s.stateChangeCallbacks.Error[i].Callback(err, transportBefore, establishedBefore)
	}
	if transportBefore == s.transport || establishedBefore > 0 {
		//TODO: Trigger chan error
	}
}

func (s *Socket) SetOnError(callback func(err error, transportBefore Transport, establishedBefore int)) string {
	ref := s.makeRef()
	s.mut.Lock()
	s.stateChangeCallbacks.Error = append(s.stateChangeCallbacks.Error, struct {
		Ref      string
		Callback func(err error, transportBefore Transport, establishedBefore int)
	}{
		Ref:      ref,
		Callback: callback,
	})
	s.mut.Unlock()
	return ref
}

func (s *Socket) OnMessage(event interface{}) {

}

func (s *Socket) SetOnMessage(callback func(event interface{})) string {
	ref := s.makeRef()
	s.mut.Lock()
	s.stateChangeCallbacks.Message = append(s.stateChangeCallbacks.Message, struct {
		Ref      string
		Callback func(event interface{})
	}{
		Ref:      ref,
		Callback: callback,
	})
	s.mut.Unlock()
	return ref
}

func (s *Socket) OnClose(event interface{}) {

}

func (s *Socket) SetOnClose(callback func(event interface{})) string {
	ref := s.makeRef()
	s.mut.Lock()
	s.stateChangeCallbacks.Close = append(s.stateChangeCallbacks.Close, struct {
		Ref      string
		Callback func(event interface{})
	}{
		Ref:      ref,
		Callback: callback,
	})
	s.mut.Unlock()
	return ref
}

// ReplaceTransport ...
func (s *Socket) ReplaceTransport(newTransport Transport) {
	s.Disconnect(nil, 0, "")
	s.transport = newTransport
}

// EndpointURL ...
func (s *Socket) EndpointURL() *url.URL {
	u, _ := url.Parse(s.endpoint)
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
	return u
}

// Disconnect ..
func (s *Socket) Disconnect(callback func(), code int, reason string) {
	s.connectClock++
	s.closeWasClean = true
	s.reconnectTimer.Reset()
	s.listenCtxCancel()
}

// Teardown ...
// TODO: implement me!!
func (s *Socket) Teardown(callback func(), code int, reason string) {

}

// Connect ..
func (s *Socket) Connect() {
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
