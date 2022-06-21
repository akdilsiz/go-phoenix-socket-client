package gophoenixsocketclient

// CloseEvent ..
type CloseEvent struct {
	Code     int
	Reason   string
	WasClean bool
}
