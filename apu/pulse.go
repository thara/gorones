package apu

import "github.com/thara/gorones/util"

type pulseChannel struct {
	volume uint8
	sweep  uint8
	low    uint8
	high   uint8

	lengthCounter uint

	enabled bool

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

func (c *pulseChannel) enable(v bool) {
	c.enabled = v
	if !v {
		c.lengthCounter = 0
	}
}

func (c *pulseChannel) dutyCycle() int { return int(c.volume >> 6) }

func (c *pulseChannel) envelopeLoop() bool      { return util.IsSet(c.volume, 5) }
func (c *pulseChannel) lengthCounterHalt() bool { return util.IsSet(c.volume, 5) }
func (c *pulseChannel) useConstantVolume() bool { return util.IsSet(c.volume, 4) }
func (c *pulseChannel) envelopePeriod() uint8   { return c.volume & 0b1111 }

func (c *pulseChannel) sweepEnabled() bool   { return util.IsSet(c.sweep, 7) }
func (c *pulseChannel) sweepPeriod() uint8   { return (c.sweep & 0b01110000) >> 4 }
func (c *pulseChannel) sweepNegate() bool    { return util.IsSet(c.sweep, 3) }
func (c *pulseChannel) sweepShift() uint8    { return c.sweep & 0b111 }
func (c *pulseChannel) sweepUnitMuted() bool { return c.timerPeriod < 8 || 0x7FF < c.timerPeriod }

func (c *pulseChannel) timerHigh() uint8         { return c.high & 0b111 }
func (c *pulseChannel) lengthCounterLoad() uint8 { return (c.high & 0b11111000) >> 3 }

func (c *pulseChannel) timerReload() uint16 { return uint16(c.low) | (uint16(c.timerHigh()) << 8) }

func (c *pulseChannel) write(addr uint16, value uint8) {
	switch addr {

	case 0x4000:
		c.volume = value
	case 0x4001:
		c.sweep = value
		c.sweepReload = true
	case 0x4002:
		c.low = value
		c.timerPeriod = c.timerReload()
	case 0x4003:
		c.high = value
		if c.enabled {
			c.lengthCounter = lengthTable[c.lengthCounterLoad()]
		}
		c.timerPeriod = c.timerReload()
		c.timerSequencer = 0
		c.envelopeStart = true

	}
}

func (c *pulseChannel) clockTimer() {
	if 0 < c.timerCounter {
		c.timerCounter -= 1
	} else {
		c.timerCounter = c.timerReload()
		c.timerSequencer += 1
		if c.timerSequencer == 8 {
			c.timerSequencer = 0
		}
	}
}

func (c *pulseChannel) clockEnvelope() {
	if c.envelopeStart {
		c.envelopeDecayLevelCounter = 15
		c.envelopeCounter = c.envelopePeriod()
		c.envelopeStart = false
		return
	}

	if 0 < c.envelopeCounter {
		c.envelopeCounter -= 1
		return
	}

	c.envelopeCounter = c.envelopePeriod()
	if 0 < c.envelopeDecayLevelCounter {
		c.envelopeDecayLevelCounter -= 1
	} else if c.envelopeDecayLevelCounter == 0 && c.envelopeLoop() {
		c.envelopeDecayLevelCounter = 15
	}
}

func (c *pulseChannel) clockSweepUnit() {
	// Updating the period
	if c.sweepCounter == 0 && c.sweepEnabled() && c.sweepShift() != 0 && !c.sweepUnitMuted() {
		var changeAmount = c.timerPeriod >> c.sweepShift()
		if c.sweepNegate() {
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
		c.sweepCounter = c.sweepPeriod()
		c.sweepReload = false
	} else {
		c.sweepCounter -= 1
	}
}

func (c *pulseChannel) clockLengthCounter() {
	if 0 < c.lengthCounter && !c.lengthCounterHalt() {
		c.lengthCounter -= 1
	}
}

func (c *pulseChannel) output() uint8 {
	if c.lengthCounter == 0 || c.sweepUnitMuted() || dutyTable[c.dutyCycle()][c.timerSequencer] == 0 {
		return 0
	}
	var volume uint8
	if c.useConstantVolume() {
		volume = c.envelopePeriod()
	} else {
		volume = c.envelopeDecayLevelCounter
	}
	return volume & 0b1111
}

var lengthTable = []uint{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}

var dutyTable [4][8]uint8 = [4][8]uint8{
	{0, 1, 0, 0, 0, 0, 0, 0}, // 12.5%
	{0, 1, 1, 0, 0, 0, 0, 0}, // 25%
	{0, 1, 1, 1, 1, 0, 0, 0}, // 50%
	{1, 0, 0, 1, 1, 1, 1, 1}, // 25% negated
}
