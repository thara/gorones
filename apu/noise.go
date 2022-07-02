package apu

type noiseChannel struct {
	enabled bool

	envelopeLoop      bool
	useConstantVolume bool
	envelopePeriod    uint8

	modeFlag bool

	envelopeCounter           uint8
	envelopeDecayLevelCounter uint8
	envelopeStart             bool

	shiftRegister uint16

	timerCounter uint16
	timerPeriod  uint16

	lengthCounter     uint8
	lengthCounterHalt bool
}

func (c *noiseChannel) write(addr uint16, value uint8) {
	switch addr {
	case 0x400C:
		c.lengthCounterHalt = (value>>5)&1 == 1
		c.envelopeLoop = (value>>5)&1 == 1
		c.useConstantVolume = (value>>4)&1 == 1
		c.envelopePeriod = value & 0b1111
		c.envelopeStart = true
	case 0x400E:
		c.modeFlag = (value>>7)&1 == 1
		c.timerPeriod = noiseTimerPeriodTable[value&0b1111]
	case 0x400F:
		if c.enabled {
			c.lengthCounter = lengthTable[value>>3]
		}
		c.envelopeStart = true
	}
}

func (c *noiseChannel) setEnabled(v bool) {
	c.enabled = v
	if !v {
		c.lengthCounter = 0
	}
}

func (c *noiseChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter -= 1
	} else {
		c.timerCounter = c.timerPeriod

		// LFSR
		var shift uint8
		if c.modeFlag {
			shift = 6
		} else {
			shift = 1
		}
		feedback := (c.shiftRegister & 1) ^ ((c.shiftRegister >> shift) & 1)
		c.shiftRegister >>= 1
		c.shiftRegister |= (feedback << 14)
	}
}

func (c *noiseChannel) clockEnvelope() {
	if c.envelopeStart {
		c.envelopeDecayLevelCounter = 15
		c.envelopeCounter = c.envelopePeriod
		c.envelopeStart = false
		return
	}

	if 0 < c.envelopeCounter {
		c.envelopeCounter--
		return
	}

	c.envelopeCounter = c.envelopePeriod
	if 0 < c.envelopeDecayLevelCounter {
		c.envelopeDecayLevelCounter--
	} else if c.envelopeLoop {
		c.envelopeDecayLevelCounter = 15
	}
}

func (c *noiseChannel) clockLengthCounter() {
	if !c.lengthCounterHalt && 0 < c.lengthCounter {
		c.lengthCounter--
	}
}

func (c *noiseChannel) output() uint8 {
	if !c.enabled || c.shiftRegister&1 == 1 || c.lengthCounter == 0 {
		return 0
	}
	if c.useConstantVolume {
		return c.envelopePeriod & 0xF
	} else {
		return c.envelopeDecayLevelCounter & 0xF
	}
}

var noiseTimerPeriodTable = []uint16{
	4, 8, 16, 32, 64, 96, 128, 160, 202, 254, 380, 508, 762, 1016, 2034, 4068,
}
