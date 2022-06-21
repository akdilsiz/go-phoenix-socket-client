package go_phoenix_socket_client

type Hook struct {
	Status   int
	Callback func(received Received)
}
