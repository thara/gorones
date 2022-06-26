package apu

import "github.com/thara/gorones/util"

type noiseChannel struct {
	envelope uint8
	period   uint8

	envelopeCounter           uint8
	envelopeDecayLevelCounter uint8
	envelopeStart             bool

	shiftRegister uint16

	timerCounter uint16
	timerPeriod  uint16
	enabled      bool

	lengthCounter lengthCounter
}

func (c *noiseChannel) envelopeLoop() bool      { return util.IsSet(c.envelope, 5) }
func (c *noiseChannel) useConstantVolume() bool { return util.IsSet(c.envelope, 4) }
func (c *noiseChannel) envelopePeriod() uint8   { return c.envelope & 0b1111 }

func (c *noiseChannel) mode() bool        { return util.IsSet(c.period, 7) }
func (c *noiseChannel) timerEntry() uint8 { return c.period & 0b1111 }

func (c *noiseChannel) write(addr uint16, value uint8) {
	switch addr {
	case 0x400C:
		c.envelope = value
		c.lengthCounter.halt = util.IsSet(c.envelope, 5)
	case 0x400E:
		c.period = value
		c.timerPeriod = noiseTimerPeriodTable[c.timerEntry()]
	case 0x400F:
		c.lengthCounter.reload((value & 0b11111000) >> 3)
		c.envelopeStart = true
	}
}

func (c *noiseChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter -= 1
	} else {
		c.timerCounter = c.timerPeriod

		// LFSR
		var i int
		if c.mode() {
			i = 6
		} else {
			i = 1
		}
		feedback := c.shiftRegister ^ util.NthBit(c.shiftRegister, i)
		c.shiftRegister >>= 1
		c.shiftRegister |= (feedback << 14)
	}
}

func (c *noiseChannel) clockEnvelope() {
	if c.envelopeStart {
		c.envelopeDecayLevelCounter = 15
		c.envelopeCounter = c.envelopePeriod()
		c.envelopeStart = false
	} else {
		if 0 < c.envelopeCounter {
			c.envelopeCounter -= 1
		} else {
			c.envelopeCounter = c.envelopePeriod()
			if 0 < c.envelopeDecayLevelCounter {
				c.envelopeDecayLevelCounter -= 1
			} else if c.envelopeLoop() {
				c.envelopeDecayLevelCounter = 15
			}
		}
	}
}

func (c *noiseChannel) output() uint8 {
	if util.NthBit(c.shiftRegister, 0) == 0 || c.lengthCounter.count == 0 {
		return 0
	}
	var v uint8
	if c.useConstantVolume() {
		v = c.envelopePeriod()
	} else {
		v = c.envelopeDecayLevelCounter
	}
	return v & 0b1111
}

var noiseTimerPeriodTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}
