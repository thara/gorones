package input

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandard_Write(t *testing.T) {
	var ctrl StandardController

	ctrl.Write(0b10101010)
	assert.Equal(t, ctrl.strobe, false)
	assert.EqualValues(t, ctrl.cur, 1)

	ctrl.Write(0b10101011)
	assert.Equal(t, ctrl.strobe, true)
	assert.EqualValues(t, ctrl.cur, 1)
}

func TestStandard_Read(t *testing.T) {
	t.Run("not strobe", func(t *testing.T) {
		var ctrl StandardController
		ctrl.Write(0b10101010)

		ctrl.Update(0b11010101)

		assert.EqualValues(t, 0x41, ctrl.Read())
		assert.EqualValues(t, 0b10, ctrl.cur)

		assert.EqualValues(t, 0x40, ctrl.Read())
		assert.EqualValues(t, 0b100, ctrl.cur)

		assert.EqualValues(t, 0x41, ctrl.Read())
		assert.EqualValues(t, 0b1000, ctrl.cur)

		assert.EqualValues(t, 0x40, ctrl.Read())
		assert.EqualValues(t, 0b10000, ctrl.cur)

		assert.EqualValues(t, 0x41, ctrl.Read())
		assert.EqualValues(t, 0b100000, ctrl.cur)

		assert.EqualValues(t, 0x40, ctrl.Read())
		assert.EqualValues(t, 0b1000000, ctrl.cur)

		assert.EqualValues(t, 0x41, ctrl.Read())
		assert.EqualValues(t, 0b10000000, ctrl.cur)

		assert.EqualValues(t, 0x41, ctrl.Read())
		assert.EqualValues(t, 0, ctrl.cur)

		// over reading
		assert.EqualValues(t, 0x040, ctrl.Read())
		assert.EqualValues(t, 0, ctrl.cur)
		assert.EqualValues(t, 0x040, ctrl.Read())
		assert.EqualValues(t, 0, ctrl.cur)
		assert.EqualValues(t, 0x040, ctrl.Read())
		assert.EqualValues(t, 0, ctrl.cur)
	})

	t.Run("strobe", func(t *testing.T) {
		var ctrl StandardController
		ctrl.Write(0b01010101)

		ctrl.Update(0b10101010)
		assert.EqualValues(t, 0x040, ctrl.Read())
		assert.EqualValues(t, 1, ctrl.cur)

		ctrl.Update(0b11101011)
		assert.EqualValues(t, 0x041, ctrl.Read())
		assert.EqualValues(t, 1, ctrl.cur)
	})
}
