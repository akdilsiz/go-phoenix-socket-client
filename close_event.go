package gophoenixsocketclient

type CloseEvent struct {
	Code     int
	Reason   string
	WasClean bool
}
