package gophoenixsocketclient

type Hook struct {
	Status   int
	Callback func(received Received)
}
