package apu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lengthCounter_clock(t *testing.T) {
	t.Run("count is over 0", func(t *testing.T) {
		t.Run("halt", func(t *testing.T) {
			lc := lengthCounter{count: 3, halt: true}
			before := lc.count
			lc.clock()
			assert.Equal(t, before, lc.count)
		})
		t.Run("not halt", func(t *testing.T) {
			lc := lengthCounter{count: 3, halt: false}
			before := lc.count
			lc.clock()
			assert.Equal(t, before-1, lc.count)
		})
	})
	t.Run("count is 0", func(t *testing.T) {
		t.Run("halt", func(t *testing.T) {
			lc := lengthCounter{count: 0, halt: true}
			before := lc.count
			lc.clock()
			assert.Equal(t, before, lc.count)
		})
		t.Run("not halt", func(t *testing.T) {
			lc := lengthCounter{count: 0, halt: false}
			before := lc.count
			lc.clock()
			assert.Equal(t, before, lc.count)
		})
	})
}
