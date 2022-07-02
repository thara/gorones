package apu

type pulseChannel struct {
	enabled bool

	volume            uint8
	dutyCycle         uint8
	envelopeLoop      bool
	useConstantVolume bool
	envelopePeriod    uint8

	sweep        uint8
	sweepEnabled bool
	sweepPeriod  uint8
	sweepNegate  bool
	sweepShift   uint8

	high uint8

	lengthCounter lengthCounter

	timerCounter   uint16
	timerSequencer int
	timerPeriod    uint16

	envelopeCounter           uint8
	envelopeDecayLevelCounter uint8
	envelopeStart             bool

	sweepCounter uint8
	sweepReload  bool

	carryMode sweepComplement
}

type sweepComplement uint16

const (
	sweepOneComplement sweepComplement = 0
	sweepTwoComplement                 = 1
)

func (c *pulseChannel) sweepUnitMuted() bool { return c.timerPeriod < 8 || 0x7FF < c.timerPeriod }

func (c *pulseChannel) write(addr uint16, value uint8) {
	switch addr {
	case 0x4000, 0x4004:
		c.volume = value
		c.dutyCycle = c.volume >> 6
		c.lengthCounter.halt = (value>>5)&1 == 0
		c.envelopeLoop = (value>>5)&1 == 1
		c.useConstantVolume = (value>>4)&1 == 1
		c.envelopePeriod = c.volume & 0b1111
	case 0x4001, 0x4005:
		c.sweep = value
		c.sweepEnabled = (value>>7)&1 == 1
		c.sweepPeriod = (c.sweep & 0b01110000) >> 4
		c.sweepNegate = (value>>3)&1 == 1
		c.sweepShift = c.sweep & 0b111
		c.sweepReload = true
	case 0x4002, 0x4006:
		c.timerPeriod = uint16(value) | (uint16(c.timerPeriod&0xFF00) << 8)
	case 0x4003, 0x4007:
		c.high = value
		if c.enabled {
			c.lengthCounter.reload(value >> 3)
		}
		c.timerPeriod = c.timerPeriod&0x00FF | (uint16(value&7) << 8)
		c.timerSequencer = 0
		c.envelopeStart = true
	}
}

func (c *pulseChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter -= 1
	} else {
		c.timerCounter = c.timerPeriod
		c.timerSequencer += 1
		if c.timerSequencer == 8 {
			c.timerSequencer = 0
		}
	}
}

func (c *pulseChannel) clockEnvelope() {
	if c.envelopeStart {
		c.envelopeDecayLevelCounter = 15
		c.envelopeCounter = c.envelopePeriod
		c.envelopeStart = false
		return
	}

	if 0 < c.envelopeCounter {
		c.envelopeCounter -= 1
		return
	}

	c.envelopeCounter = c.envelopePeriod
	if 0 < c.envelopeDecayLevelCounter {
		c.envelopeDecayLevelCounter -= 1
	} else if c.envelopeDecayLevelCounter == 0 && c.envelopeLoop {
		c.envelopeDecayLevelCounter = 15
	}
}

func (c *pulseChannel) clockSweepUnit() {
	// Updating the period
	if c.sweepCounter == 0 && c.sweepEnabled && c.sweepShift != 0 && !c.sweepUnitMuted() {
		var changeAmount = c.timerPeriod >> c.sweepShift
		if c.sweepNegate {
			switch c.carryMode {
			case sweepOneComplement:
				changeAmount = ^changeAmount
			case sweepTwoComplement:
				changeAmount = ^changeAmount + 1
			}
		}
		c.timerPeriod += changeAmount
	}

	if c.sweepCounter == 0 || c.sweepReload {
		c.sweepCounter = c.sweepPeriod
		c.sweepReload = false
	} else {
		c.sweepCounter -= 1
	}
}

func (c *pulseChannel) output() uint8 {
	if !c.enabled || c.lengthCounter.count == 0 || c.sweepUnitMuted() || dutyTable[c.dutyCycle][c.timerSequencer] == 0 {
		return 0
	}
	if c.useConstantVolume {
		return c.envelopePeriod
	} else {
		return c.envelopeDecayLevelCounter
	}
}

var dutyTable [4][8]uint8 = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0}, // 12.5%
	{0, 1, 1, 0, 0, 0, 0, 0}, // 25%
	{0, 1, 1, 1, 1, 0, 0, 0}, // 50%
	{1, 0, 0, 1, 1, 1, 1, 1}, // 25% negated
}
