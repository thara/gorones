package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_audioBuffer_write(t *testing.T) {
	b := NewAudioBuffer(5)

	b.Write(12.3)
	assert.Equal(t, 4, b.writePos)
	assert.EqualValues(t, 205, b.buf[0])
	assert.EqualValues(t, 204, b.buf[1])
	assert.EqualValues(t, 68, b.buf[2])
	assert.EqualValues(t, 65, b.buf[3])
	assert.EqualValues(t, 0, b.buf[4])

	b.Write(12.3)
	assert.Equal(t, 3, b.writePos)
	assert.EqualValues(t, 205, b.buf[4])
	assert.EqualValues(t, 204, b.buf[0])
	assert.EqualValues(t, 68, b.buf[1])
	assert.EqualValues(t, 65, b.buf[2])
}

func Test_audioBuffer_read(t *testing.T) {
	b := NewAudioBuffer(5)
	b.Write(12.3)

	buf := make([]byte, 2)
	n, _ := b.Read(buf)
	assert.Equal(t, 2, n)
	assert.EqualValues(t, 205, buf[0])
	assert.EqualValues(t, 204, buf[1])

	n, _ = b.Read(buf)
	assert.Equal(t, 2, n)
	assert.EqualValues(t, 68, buf[0])
	assert.EqualValues(t, 65, buf[1])

	n, _ = b.Read(buf)
	assert.Equal(t, 2, n)
	assert.EqualValues(t, 0, buf[0])
	assert.EqualValues(t, 205, buf[1])

	buf = make([]byte, 3)
	n, _ = b.Read(buf)
	assert.Equal(t, 3, n)
	assert.EqualValues(t, 204, buf[0])
	assert.EqualValues(t, 68, buf[1])
	assert.EqualValues(t, 65, buf[2])
}
