package gophoenixsocketclient

// StateChange ..
type StateChange struct {
	Open []struct {
		Ref      string
		Callback func()
	}
	Close []struct {
		Ref      string
		Callback func(event interface{})
	}
	Error []struct {
		Ref      string
		Callback func(err error, transportBefore Transport, establishedBefore int)
	}
	Message []struct {
		Ref      string
		Callback func(event interface{})
	}
}
