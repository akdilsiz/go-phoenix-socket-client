package gophoenixsocketclient

type Received struct {
	Status   string
	Response []byte
	Ref      string
}
