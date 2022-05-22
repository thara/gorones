package apu

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_lengthCounter(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	t.Run("not enabled", func(t *testing.T) {
		c := runLengthCounter(ctx)

		c.load(0b11111000)
		c.clock()
		assert.EqualValues(t, 0, c.counter())
	})

	t.Run("enabled & not halt", func(t *testing.T) {
		c := runLengthCounter(ctx)

		c.setEnabled(true)
		c.load(0b00101000) // 5
		assert.EqualValues(t, 4, c.counter())
		c.clock()
		assert.EqualValues(t, 3, c.counter())
		c.clock()
		assert.EqualValues(t, 2, c.counter())
		c.clock()
		assert.EqualValues(t, 1, c.counter())
		c.clock()
		assert.EqualValues(t, 0, c.counter())
		c.clock()
		assert.EqualValues(t, 0, c.counter())
	})

	t.Run("enabled & halt", func(t *testing.T) {
		c := runLengthCounter(ctx)

		c.setEnabled(true)
		c.load(0b00101000) // 5
		c.setHalt(true)
		assert.EqualValues(t, 4, c.counter())
		c.clock()
		assert.EqualValues(t, 4, c.counter())
	})

	t.Run("clear counter by enabled on", func(t *testing.T) {
		c := runLengthCounter(ctx)

		c.setEnabled(true)
		c.load(0b00101000) // 5
		assert.EqualValues(t, 4, c.counter())
		c.setEnabled(false)
		assert.EqualValues(t, 0, c.counter())
	})
}
