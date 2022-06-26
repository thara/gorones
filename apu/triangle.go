package apu

import (
	"github.com/thara/gorones/util"
)

type triangleChannel struct {
	linearCounterSetup uint8
	low                uint8
	high               uint8

	linearCounterReloadFlag bool

	timerCounter uint16
	sequencer    uint8

	linearCounter uint8

	lengthCounter lengthCounter
}

func (c *triangleChannel) controlFlag() bool          { return util.IsSet(c.linearCounterSetup, 7) }
func (c *triangleChannel) linearCounterReload() uint8 { return c.linearCounterSetup & 0b01111111 }

func (c *triangleChannel) timerLow() uint8          { return c.low }
func (c *triangleChannel) timerHigh() uint8         { return c.high & 0b111 }
func (c *triangleChannel) lengthCounterLoad() uint8 { return (c.high & 0b11111000) >> 3 }

func (c *triangleChannel) timerReload() uint16 { return uint16(c.low) | (uint16(c.timerHigh()) << 8) }

func (c *triangleChannel) write(addr uint16, value uint8) {
	switch addr {
	case 0x4008:
		c.linearCounterSetup = value
		c.lengthCounter.halt = util.IsSet(c.linearCounterSetup, 7)
	case 0x400A:
		c.low = value
	case 0x400B:
		c.high = value
		c.linearCounterReloadFlag = true
		c.lengthCounter.reload(c.lengthCounterLoad())
	default:
		break
	}
}

func (c *triangleChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter -= 1
	} else {
		c.timerCounter = c.timerReload()
		if 0 < c.linearCounter && 0 < c.lengthCounter.count {
			c.sequencer += 1
			if c.sequencer == 32 {
				c.sequencer = 0
			}
		}
	}
}

func (c *triangleChannel) clockLinearCounter() {
	if c.linearCounterReloadFlag {
		c.linearCounter = c.linearCounterReload()
	} else {
		c.linearCounter -= 1
	}

	if c.controlFlag() {
		c.linearCounterReloadFlag = false
	}
}

func (c *triangleChannel) output() uint8 {
	if c.controlFlag() || !c.lengthCounter.enabled || c.lengthCounter.count == 0 || c.linearCounter == 0 {
		return 0
	}
	// 15, 14, 13, 12, 11, 10,  9,  8,  7,  6,  5,  4,  3,  2,  1,  0
	//  0,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10, 11, 12, 13, 14, 15
	s := int(c.sequencer)
	v := int(s) - 15 - (s / 16)
	if v < 0 {
		return uint8(-v)
	}
	return uint8(v)
}
