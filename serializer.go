package gophoenixsocketclient

import "fmt"

// Serializer ..
type Serializer interface {
	Encode(message Message, callback func([]byte))
	Decode(rawPayload []byte, callback func([]byte))
}

type serializer struct {
	headerLength int64
	metaLength   int64
	kinds        struct {
		push      byte
		reply     byte
		broadcast byte
	}
}

// NewSerializer initialize default serializer
func NewSerializer() Serializer {
	e := &serializer{
		headerLength: 1,
		metaLength:   4,
		kinds: struct {
			push      byte
			reply     byte
			broadcast byte
		}{push: 0, reply: 1, broadcast: 2},
	}
	return e
}

func (e serializer) Encode(message Message, callback func([]byte)) {

}

func (e serializer) Decode(rawPayload []byte, callback func([]byte)) {

}

func (e serializer) binaryEncode(message Message) *Buffer {
	var metaLength int64
	var offset int64 = 0

	metaLength = e.metaLength +
		int64(len(message.JoinRef)) +
		int64(len(message.Ref)) +
		int64(len(message.Topic)) +
		int64(len(message.Event))

	header := NewBuffer(e.headerLength + metaLength)

	offset++
	_, _ = header.WriteAt([]byte{e.kinds.push}, offset)
	offset++
	_, _ = header.WriteAt([]byte{uint8(len(message.JoinRef))}, offset)
	offset++
	_, _ = header.WriteAt([]byte{uint8(len(message.Ref))}, offset)
	offset++
	_, _ = header.WriteAt([]byte{uint8(len(message.Topic))}, offset)
	offset++
	_, _ = header.WriteAt([]byte{uint8(len(message.Event))}, offset)

	for i := 0; i < len(message.JoinRef); i++ {
		offset++
		_, _ = header.WriteAt([]byte(fmt.Sprintf("%#U", message.JoinRef[i])), offset)
	}

	for i := 0; i < len(message.Ref); i++ {
		offset++
		_, _ = header.WriteAt([]byte(fmt.Sprintf("%#U", message.JoinRef[i])), offset)
	}

	for i := 0; i < len(message.Topic); i++ {
		offset++
		_, _ = header.WriteAt([]byte(fmt.Sprintf("%#U", message.JoinRef[i])), offset)
	}

	for i := 0; i < len(message.Event); i++ {
		offset++
		_, _ = header.WriteAt([]byte(fmt.Sprintf("%#U", message.JoinRef[i])), offset)
	}

	combined := NewBuffer(header.Length() + int64(len(message.Payload)))
	_, _ = combined.WriteAt(header.Bytes(), 0)
	_, _ = combined.WriteAt(message.Payload, header.Length())

	header.Reset()

	return combined
}
