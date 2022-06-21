package gophoenixsocketclient

// Hook ..
type Hook struct {
	Status   int
	Callback func(received Received)
}
