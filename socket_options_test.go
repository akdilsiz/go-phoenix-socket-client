package gophoenixsocketclient

import (
	"github.com/akdilsiz/go-phoenix-socket-client/constants"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestDefaultSocketOptions(t *testing.T) {
	defaultSocketOptions := DefaultSocketOptions()
	assert.Equal(t, constants.DefaultTimeout, defaultSocketOptions.Timeout)
	assert.Equal(t, constants.Websocket, defaultSocketOptions.Transport)
	assert.Equal(t, "*gophoenixsocketclient.serializer", reflect.TypeOf(defaultSocketOptions.Encode).String())
	assert.Equal(t, "*gophoenixsocketclient.serializer", reflect.TypeOf(defaultSocketOptions.Decode).String())
	assert.Equal(t, time.Second*30, defaultSocketOptions.HeartbeatIntervalMS)
	assert.Equal(t, time.Second*10, defaultSocketOptions.RejoinAfterMS(0))
	assert.Equal(t, time.Millisecond*5000, defaultSocketOptions.ReconnectAfterMS(0))
	assert.Len(t, defaultSocketOptions.Params, 0)
}
