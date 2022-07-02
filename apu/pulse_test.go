package apu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_pulse_write(t *testing.T) {
	t.Run("0x4000", func(t *testing.T) {
		var c pulseChannel
		c.write(0x4000, 0b10111111)

		assert.EqualValues(t, 0b10, c.dutyCycle)
		assert.True(t, c.lengthCounter.halt)
		assert.True(t, c.useConstantVolume)
		assert.EqualValues(t, 0b1111, c.envelopePeriod)
	})
	t.Run("0x4001", func(t *testing.T) {
		var c pulseChannel
		c.write(0x4001, 0b10101011)

		assert.True(t, c.sweepEnabled)
		assert.EqualValues(t, 0b010, c.sweepPeriod)
		assert.True(t, c.sweepNegate)
		assert.EqualValues(t, 0b011, c.sweepShift)
	})
	t.Run("0x4002 & 0x4003", func(t *testing.T) {
		var c pulseChannel
		c.enabled = true
		c.write(0x4002, 0b00000100)
		c.write(0x4003, 0b11101011)

		assert.EqualValues(t, 0b011_00000100, c.timerPeriod)
		assert.EqualValues(t, 28, c.lengthCounter.count)
	})
}

func Test_pulse_clockTimer(t *testing.T) {
	t.Run("counter is greater than zero", func(t *testing.T) {
		c := pulseChannel{timerCounter: 3}
		c.write(0x4002, 0b11) // timer low
		c.clockTimer()
		assert.EqualValues(t, 2, c.timerCounter)
	})

	setup := func(c *pulseChannel) {
		for i := 0; i < 3; i++ {
			c.clockTimer()
		}
	}

	t.Run("counter is zero", func(t *testing.T) {
		t.Run("reloads timerCounter and increments sequencer", func(t *testing.T) {
			c := pulseChannel{timerCounter: 3}
			c.write(0x4002, 0b11) // timer low
			setup(&c)
			before := c.timerSequencer
			c.clockTimer()

			assert.EqualValues(t, 0b11, c.timerCounter)
			assert.EqualValues(t, before+1, c.timerSequencer)
		})

		t.Run("if sequencer become over 8", func(t *testing.T) {
			c := pulseChannel{timerCounter: 3}
			c.write(0x4002, 0b11) // timer low
			setup(&c)
			c.timerSequencer = 7
			c.clockTimer()
			assert.EqualValues(t, 0, c.timerSequencer)
		})
	})

}

func Test_pulse_clockEnvelope(t *testing.T) {
	t.Run("start is on", func(t *testing.T) {
		c := pulseChannel{envelopeStart: true}
		c.write(0x4000, 0b111) // volume
		c.clockEnvelope()

		assert.EqualValues(t, 15, c.envelopeDecayLevelCounter)
		assert.EqualValues(t, c.envelopePeriod, c.envelopeCounter)
		assert.False(t, c.envelopeStart)
	})

	t.Run("start is off", func(t *testing.T) {
		t.Run("envelope's counter is greater than zero after clocked", func(t *testing.T) {
			c := pulseChannel{}
			c.write(0x4000, 0b111) // volume
			c.envelopeCounter = c.envelopePeriod

			before := c.envelopeCounter
			c.clockEnvelope()

			assert.EqualValues(t, before-1, c.envelopeCounter)
		})
		t.Run("envelope's counter is zero after clocked", func(t *testing.T) {
			t.Run("reloads envelope's counter", func(t *testing.T) {
				c := pulseChannel{envelopeCounter: 0}
				c.write(0x4000, 0b111) // volume

				c.clockEnvelope()
				assert.EqualValues(t, c.envelopePeriod, c.envelopeCounter)
			})
			t.Run("envelope's decayLevelCounter become to be greater than 0 after clocked", func(t *testing.T) {
				c := pulseChannel{envelopeDecayLevelCounter: 2}
				c.write(0x4000, 0b111) // volume

				before := c.envelopeDecayLevelCounter
				c.clockEnvelope()
				assert.EqualValues(t, before-1, c.envelopeDecayLevelCounter)
			})
			t.Run("envelope's decayLevelCounter become 0 after clocked", func(t *testing.T) {
				c := pulseChannel{}
				c.write(0x4000, 0b100000) // volume
				c.envelopeCounter = c.envelopePeriod
				c.clockEnvelope()
				assert.EqualValues(t, 15, c.envelopeDecayLevelCounter)
			})
		})
	})
}

func Test_pulse_sweepUnitMuted(t *testing.T) {
	tests := []struct {
		timerPeriod uint16
		want        bool
	}{
		{timerPeriod: 7, want: true},
		{timerPeriod: 0x800, want: true},
		{timerPeriod: 8, want: false},
		{timerPeriod: 0x7FE, want: false},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("timerPeriod:%d", tt.timerPeriod), func(t *testing.T) {
			c := pulseChannel{timerPeriod: tt.timerPeriod}
			got := c.sweepUnitMuted()
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_pulse_clockSweepUnit(t *testing.T) {
	t.Run("sweep unit counter is not 0", func(t *testing.T) {
		c := pulseChannel{carryMode: sweepOneComplement, sweepCounter: 3}
		c.write(0x4001, 0b10101000) // sweep

		before := c.sweepCounter
		c.clockSweepUnit()
		assert.Equal(t, before-1, c.sweepCounter)
	})
	t.Run("sweep unit counter is 0", func(t *testing.T) {
		c := pulseChannel{carryMode: sweepOneComplement, sweepCounter: 0, sweepReload: true}
		c.write(0x4001, 0b10101000) // sweep

		c.clockSweepUnit()
		assert.Equal(t, c.sweepPeriod, c.sweepCounter)
		assert.False(t, c.sweepReload)
	})
	t.Run("sweep unit reload is true", func(t *testing.T) {
		c := pulseChannel{carryMode: sweepOneComplement, sweepCounter: 1, sweepReload: true}
		c.write(0x4001, 0b10101000) // sweep

		c.clockSweepUnit()
		assert.Equal(t, c.sweepPeriod, c.sweepCounter)
		assert.False(t, c.sweepReload)
	})
	t.Run("sweep unit counter is zero and enabled and not muted", func(t *testing.T) {
		t.Run("reloads sweep unit counter and clear reload flag", func(t *testing.T) {
			// not muted
			c := pulseChannel{carryMode: sweepOneComplement, timerPeriod: 0b1000}
			c.write(0x4001, 0b10000001) // sweep

			before := c.timerPeriod
			c.clockSweepUnit()
			assert.Equalf(t, before+0b100, c.timerPeriod, "before:%x", before)
		})
		t.Run("if negated with one's complement", func(t *testing.T) {
			// negate = true, shift count = 2
			c := pulseChannel{carryMode: sweepOneComplement, timerPeriod: 0b101010010}
			c.write(0x4001, 0b10101010) // sweep
			c.clockSweepUnit()
			// 0b101010010 >> 2 -> 0b1010100
			// 0b1010100 * -1 - 1 -> -85
			// 0b101010010 - 85 -> 253
			assert.EqualValues(t, 253, c.timerPeriod)
		})
		t.Run("if negated with two's complement", func(t *testing.T) {
			// negate = true, shift count = 2
			c := pulseChannel{carryMode: sweepTwoComplement, timerPeriod: 0b101010010}
			c.write(0x4001, 0b10101010) // sweep
			c.clockSweepUnit()
			// 0b101010010 >> 2 -> 0b1010100
			// 0b1010100 * -1 -> -84
			// 0b101010010 - 84 -> 254
			assert.EqualValues(t, 254, c.timerPeriod)
		})
	})
}

func Test_pulse_enable(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		var c pulseChannel
		c.enabled = true
		c.write(0x4003, 0b10101000)
		// 1 0101 (21)
		assert.EqualValues(t, 0x14, c.lengthCounter.count)
	})
	t.Run("disabled", func(t *testing.T) {
		var c pulseChannel
		c.enabled = false
		before := c.lengthCounter.count
		c.write(0x4003, 0b11)
		// 1 0101 (21)
		assert.EqualValues(t, before, c.lengthCounter.count)
	})
}
