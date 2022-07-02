package apu

import "github.com/thara/gorones/util"

type dmc struct {
	enabled bool

	flags   uint8
	direct  uint8
	address uint8
	length  uint8

	timerCounter uint8

	bitsRemainingCounter uint8

	sampleBuffer uint8

	// Memory reader
	addressCounter        uint16
	bytesRemainingCounter uint16

	outputLevel uint8

	silence           bool
	sampleBufferEmpty bool

	shiftRegister uint8

	interrupted bool
}

func (c *dmc) irqEnabled() bool { return util.IsSet(c.flags, 7) }
func (c *dmc) loopFlag() bool   { return util.IsSet(c.flags, 6) }
func (c *dmc) rateIndex() uint8 { return c.flags & 0b1111 }

func (c *dmc) directLoad() uint8 { return c.direct & 0b01111111 }

func (c *dmc) sampleAddress() uint16 { return 0xC000 + uint16(c.address)*64 }
func (c *dmc) sampleLength() uint16  { return uint16(c.length)*16 + 1 }

func (c *dmc) write(addr uint16, value uint8) {
	switch addr {
	case 0x4010:
		c.flags = value
	case 0x4011:
		c.direct = value
		c.outputLevel = c.directLoad()
	case 0x4012:
		c.address = value
	case 0x4013:
		c.length = value
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
		c.timerCounter -= 1
	} else {
		// the output cycle ends
		c.timerCounter = 8

		// Memory Reader
		if c.sampleBufferEmpty && c.bytesRemainingCounter != 0 {
			c.sampleBuffer = memoryReader.Read(c.addressCounter)
			c.addressCounter += 1
			if c.addressCounter == 0 {
				c.addressCounter = 0x8000
			}
			c.bytesRemainingCounter -= 1

			if c.bytesRemainingCounter == 0 {
				if c.loopFlag() {
					c.start()
				}
				if c.irqEnabled() {
					c.interrupted = true
				}
			}

			cpuStall = true
		}

		// Output unit
		if c.sampleBufferEmpty {
			c.silence = true
		} else {
			c.silence = false
			c.shiftRegister = c.sampleBuffer
			c.sampleBufferEmpty = true
			c.sampleBuffer = 0
		}

		if !c.silence {
			if util.NthBit(c.shiftRegister, 0) == 1 {
				if c.outputLevel < c.outputLevel&+2 {
					c.outputLevel += 2
				}
			} else {
				if c.outputLevel-2 < c.outputLevel {
					c.outputLevel -= 2
				}
			}
		}
		c.shiftRegister >>= 1
		c.bitsRemainingCounter -= 1
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
	return c.outputLevel & 0b01111111
}
