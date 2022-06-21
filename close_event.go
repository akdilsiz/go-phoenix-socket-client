package go_phoenix_socket_client

type CloseEvent struct {
	Code     int
	Reason   string
	WasClean bool
}
