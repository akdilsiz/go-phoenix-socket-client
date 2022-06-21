package gophoenixsocketclient

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
