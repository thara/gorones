package apu

import (
	"context"
	"testing"
	"time"
)

func Test_envelope_noStart(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	e := runEnvelope(ctx)

	e.clock()
	assertRecvValue(t, e.output(), 0)
	e.clock()
	assertRecvValue(t, e.output(), 0)

	t.Run("reload", func(t *testing.T) {
		e := runEnvelope(ctx)
		e.update(0b1001) // v:9
		// v + 1
		for i := 0; i < 10; i++ {
			e.clock()
		}
		assertRecvValue(t, e.output(), 0)
	})

	t.Run("const volume", func(t *testing.T) {
		e := runEnvelope(ctx)
		e.update(0b11001) // const + v:9
		// v + 1
		for i := 0; i < 10; i++ {
			e.clock()
		}
		assertRecvValue(t, e.output(), 0b1001)
	})
}

func Test_envelope_start(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	t.Run("non const volume", func(t *testing.T) {
		e := runEnvelope(ctx)
		e.start()

		e.update(0b000011) // start on + v:3
		// start + v + 1
		envelopeClock(e, 5)
		assertRecvValue(t, e.output(), 14)

		for i := uint8(13); 0 < i; i-- {
			envelopeClock(e, 4)
			assertRecvValue(t, e.output(), i)
		}
		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 0)

		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 0, "no loop")
	})

	t.Run("non const volume + loop", func(t *testing.T) {
		e := runEnvelope(ctx)
		e.start()

		e.update(0b100011) // start on + v:3
		// start + v + 1
		envelopeClock(e, 5)
		assertRecvValue(t, e.output(), 14)

		for i := uint8(13); 0 < i; i-- {
			envelopeClock(e, 4)
			assertRecvValue(t, e.output(), i)
		}
		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 0)

		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 15, "loop")
	})

	t.Run("const volume", func(t *testing.T) {
		e := runEnvelope(ctx)
		e.start()

		e.update(0b010011) // start on + v:3
		// start + v + 1
		envelopeClock(e, 5)
		assertRecvValue(t, e.output(), 0b0011)

		for i := uint8(13); 0 < i; i-- {
			envelopeClock(e, 4)
			assertRecvValue(t, e.output(), 0b0011)
		}
		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 0b0011)

		envelopeClock(e, 4)
		assertRecvValue(t, e.output(), 0b0011)
	})
}

func envelopeClock(e *envelope, n int) {
	for i := 0; i < n; i++ {
		e.clock()
	}
}
