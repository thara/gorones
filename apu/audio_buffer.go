package apu

import (
	"encoding/binary"
	"math"
)

type AudioBuffer struct {
	buf []byte

	readPos, writePos int
}

func NewAudioBuffer(cap int) AudioBuffer {
	return AudioBuffer{
		buf: make([]byte, cap),
	}
}

func (b *AudioBuffer) Write(v float32) {
	n := math.Float32bits(v)

	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, n)
	for _, v := range buf {
		if len(b.buf) <= b.writePos {
			b.writePos = 0
		}
		b.buf[b.writePos] = v
		b.writePos++
	}
}

func (b *AudioBuffer) Read(p []byte) (n int, err error) {
	end := b.readPos + len(p)

	var over bool
	if len(b.buf) <= end {
		end = len(b.buf)
		over = true
	}

	n = copy(p, b.buf[b.readPos:end])

	b.readPos += len(p)

	if over {
		rem := end - (len(b.buf) - 1)
		n += copy(p[n:], b.buf[:rem])
		b.readPos = rem
	}
	return
}
