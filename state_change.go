package gophoenixsocketclient

// State ..
type State struct {
	Ref      int64
	Callback func()
}

// StateChange ..
type StateChange struct {
	Open    []State
	Close   []State
	Error   []State
	Message []State
}
