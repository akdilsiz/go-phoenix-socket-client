package go_phoenix_socket_client

type State struct {
	Ref      int64
	Callback func()
}

type StateChange struct {
	Open    []State
	Close   []State
	Error   []State
	Message []State
}
