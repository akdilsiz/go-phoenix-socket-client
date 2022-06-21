package gophoenixsocketclient

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	buffer := NewBuffer(1024)
	assert.Len(t, buffer.buf, 1024)
	assert.Equal(t, int64(1024), buffer.l)
	assert.Equal(t, float64(0), buffer.Grow)
}

func TestNewBufferFromValue(t *testing.T) {
	buf := make([]byte, 1024)
	buffer := NewBufferFromValues(buf)
	assert.Len(t, buffer.buf, 1024)
	assert.Equal(t, int64(1024), buffer.l)
	assert.Equal(t, float64(0), buffer.Grow)
}

func TestBuffer_WriteAt(t *testing.T) {
	buf := make([]byte, 3)
	buf[0] = uint8(0)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	buffer := NewBuffer(1024)
	n, err := buffer.WriteAt(buf, 1021)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, uint8(0), buffer.buf[1021])
	assert.Equal(t, uint8(1), buffer.buf[1022])
	assert.Equal(t, uint8(2), buffer.buf[1023])
}

func TestBuffer_WriteAtWithBigValue(t *testing.T) {
	buf := make([]byte, 3)
	buf[0] = uint8(0)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	buffer := NewBuffer(1024)
	n, err := buffer.WriteAt(buf, 2048)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, uint8(0), buffer.buf[2048])
	assert.Equal(t, uint8(1), buffer.buf[2049])
	assert.Equal(t, uint8(2), buffer.buf[2050])
	assert.Equal(t, float64(1), buffer.Grow)
}

func TestBuffer_Bytes(t *testing.T) {
	buf := make([]byte, 3)
	buf[0] = uint8(0)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	buffer := NewBuffer(3)
	n, err := buffer.WriteAt(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte{0, 1, 2}, buffer.Bytes())
}

func TestBuffer_Length(t *testing.T) {
	buf := make([]byte, 3)
	buf[0] = uint8(0)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	buffer := NewBuffer(4)
	n, err := buffer.WriteAt(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, int64(4), buffer.Length())
}

func TestBuffer_Reset(t *testing.T) {
	buf := make([]byte, 3)
	buf[0] = uint8(0)
	buf[1] = uint8(1)
	buf[2] = uint8(2)
	buffer := NewBuffer(3)
	n, err := buffer.WriteAt(buf, 0)
	assert.Nil(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, []byte{0, 1, 2}, buffer.Bytes())

	buffer.Reset()
	assert.Equal(t, int64(0), buffer.Length())
	assert.Equal(t, []byte(nil), buffer.buf)
}
