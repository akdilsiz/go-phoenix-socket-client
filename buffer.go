package go_phoenix_socket_client

import "sync"

// Buffer Offset writable memory buffer.
// Reference by: https://github.com/aws/aws-sdk-go/blob/v1.44.33/aws/types.go#L163-L172
type Buffer struct {
	buf  []byte
	l    int64
	mut  sync.Mutex
	Grow float64
}

// NewBuffer make offset writable memory buffer.
func NewBuffer(l int64) *Buffer {
	return &Buffer{buf: make([]byte, l), l: l}
}

func NewBufferFromValues(buf []byte) *Buffer {
	return &Buffer{buf: buf, l: int64(cap(buf))}
}

// WriteAt write memory buffer with offset position
func (buf *Buffer) WriteAt(payload []byte, offset int64) (n int, err error) {
	n = len(payload)
	afterLength := offset + int64(n)
	buf.mut.Lock()
	defer buf.mut.Unlock()
	if int64(len(buf.buf)) < afterLength {
		if int64(cap(buf.buf)) < afterLength {
			if buf.Grow < 1 {
				buf.Grow = 1
			}
			newBuffer := make([]byte, afterLength, int64(buf.Grow*float64(afterLength)))
			copy(newBuffer, buf.buf)
			buf.buf = newBuffer
		}
		buf.buf = buf.buf[:afterLength]
	}
	copy(buf.buf[offset:], payload)
	return n, err
}

// Bytes returns all bytes in the memory buffer.
func (buf *Buffer) Bytes() []byte {
	buf.mut.Lock()
	defer buf.mut.Unlock()
	return buf.buf
}

// Reset clears the memory buffer.
func (buf *Buffer) Reset() {
	buf.buf = nil
}

// Length returns the size of the memory buffer.
func (buf *Buffer) Length() int64 {
	buf.mut.Lock()
	defer buf.mut.Unlock()
	return buf.l
}
