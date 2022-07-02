package apu

type triangleChannel struct {
	enabled bool

	controlFlag         bool
	linearCounterReload uint8

	timerPeriod uint16

	linearCounterReloadFlag bool

	timerCounter uint16
	sequencer    uint8

	linearCounter uint8

	lengthCounter     uint8
	lengthCounterHalt bool
}

func (c *triangleChannel) write(addr uint16, value uint8) {
	switch addr {
	case 0x4008:
		c.controlFlag = (value>>7)&1 == 1
		c.lengthCounterHalt = (value>>7)&1 == 1
		c.linearCounterReload = value & 0b01111111
	case 0x400A:
		c.timerPeriod = uint16(value) | (c.timerPeriod & 0xFF00)
	case 0x400B:
		c.timerPeriod = (c.timerPeriod & 0x00FF) | (uint16(value&7) << 8)
		c.linearCounterReloadFlag = true
		if c.enabled {
			c.lengthCounter = lengthTable[(value&0b11111000)>>3]
		}
	}
}

func (c *triangleChannel) setEnabled(v bool) {
	c.enabled = v
	if !v {
		c.lengthCounter = 0
	}
}

func (c *triangleChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter--
	} else {
		c.timerCounter = c.timerPeriod
		if 0 < c.linearCounter && 0 < c.lengthCounter {
			c.sequencer++
			if c.sequencer == 32 {
				c.sequencer = 0
			}
		}
	}
}

func (c *triangleChannel) clockLinearCounter() {
	if c.linearCounterReloadFlag {
		c.linearCounter = c.linearCounterReload
	} else {
		c.linearCounter--
	}

	if !c.controlFlag {
		c.linearCounterReloadFlag = false
	}
}

func (c *triangleChannel) clockLengthCounter() {
	if !c.lengthCounterHalt && 0 < c.lengthCounter {
		c.lengthCounter--
	}
}

func (c *triangleChannel) output() uint8 {
	if !c.enabled || c.lengthCounter == 0 || c.linearCounter == 0 || c.timerPeriod < 2 {
		return 0
	}
	return sequencerTable[c.sequencer] & 0xF
}

var sequencerTable = [32]uint8{
	15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1, 0,
	0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15,
}
