package apu

type dmc struct {
	enabled bool

	flags      uint8
	irqEnabled bool
	loopFlag   bool
	rateIndex  uint8

	direct  uint8
	address uint8
	length  uint8

	timerCounter uint16
	timerPeriod  uint16

	remainingBytesCounter uint8

	sampleBuffer uint8

	// Memory reader
	addressCounter        uint16
	bytesRemainingCounter uint16

	outputLevel uint8

	silence           bool
	sampleBufferEmpty bool

	shiftRegister        uint8
	remainingBitsCounter uint8

	interrupted bool
}

func (c *dmc) directLoad() uint8 { return c.direct & 0b01111111 }

func (c *dmc) sampleAddress() uint16 { return 0xC000 + uint16(c.address)*64 }

func (c *dmc) sampleLength() uint16 { return uint16(c.length)*16 + 1 }

func (c *dmc) write(addr uint16, value uint8) {
	switch addr {
	case 0x4010:
		c.flags = value
		c.irqEnabled = (value>>7)&1 == 1
		c.loopFlag = (value>>6)&1 == 1
		c.rateIndex = value & 0b1111
		c.timerPeriod = DMC_TIMER_TABLE[value&0xF] >> 1
	case 0x4011:
		c.direct = value
		c.outputLevel = c.directLoad()
	case 0x4012:
		c.address = value
		c.addressCounter = c.sampleAddress()
	case 0x4013:
		c.length = value
		c.bytesRemainingCounter = c.sampleLength()
	default:
		break
	}
}

type DMCMemoryReader interface {
	Read(uint16) uint8
}

func (c *dmc) clockTimer(memoryReader DMCMemoryReader) bool {
	var cpuStall = false

	if 0 < c.timerCounter {
		c.timerCounter--
	} else {
		// the output cycle ends
		c.timerCounter = c.timerPeriod

		// Memory Reader
		if c.sampleBufferEmpty && c.bytesRemainingCounter != 0 {
			c.sampleBuffer = memoryReader.Read(c.addressCounter)
			c.addressCounter++
			if c.addressCounter == 0 {
				c.addressCounter = 0x8000
			}
			c.sampleBufferEmpty = false
			c.bytesRemainingCounter--

			if c.bytesRemainingCounter == 0 {
				if c.loopFlag {
					c.start()
				} else if c.irqEnabled {
					c.interrupted = true
				}
			}

			cpuStall = true
		}

		// Output unit
		if c.remainingBitsCounter == 0 {
			c.remainingBitsCounter = 8

			if c.sampleBufferEmpty {
				c.silence = true
			} else {
				c.silence = false
				c.shiftRegister = c.sampleBuffer
				c.sampleBufferEmpty = true
				c.sampleBuffer = 0
			}
		}

		if !c.silence {
			if c.shiftRegister&1 == 0 {
				if c.outputLevel < 1 {
					c.outputLevel -= 2
				}
			} else {
				if c.outputLevel < 126 {
					c.outputLevel += 2
				}
			}
		}
		c.remainingBitsCounter--
		c.shiftRegister >>= 1
	}
	return cpuStall
}

func (c *dmc) start() {
	c.outputLevel = c.directLoad()
	c.addressCounter = c.sampleAddress()
	c.bytesRemainingCounter = c.sampleLength()
}

func (c *dmc) output() uint8 {
	if c.silence {
		return 0
	}
	return c.outputLevel & 0x7F
}

var DMC_TIMER_TABLE = [16]uint16{
	0x1AC, 0x17C, 0x154, 0x140, 0x11E, 0x0FE, 0x0E2, 0x0D6,
	0x0BE, 0x0A0, 0x08E, 0x080, 0x06A, 0x054, 0x048, 0x036,
}
