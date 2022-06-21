package gophoenixsocketclient

import "github.com/akdilsiz/go-phoenix-socket-client/constants"

// Message ..
type Message struct {
	JoinRef string                 `json:"join_ref"`
	Ref     string                 `json:"ref"`
	Event   constants.ChannelEvent `json:"event"`
	Topic   string                 `json:"topic"`
	Payload []byte                 `json:"payload"`
}
