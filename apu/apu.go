package apu

import "github.com/thara/gorones/util"

type APU struct {
	sampleRate  uint
	framePeriod uint

	pulse1   pulseChannel
	pulse2   pulseChannel
	triangle triangleChannel
	noise    noiseChannel
	dmc      dmc

	cycles uint

	frameCounterControl uint8
	frameSequenceStep   int
	frameInterrupted    bool

	audio           Audio
	dmcMemoryReader DMCMemoryReader
}

func New(audio Audio, dmcMemoryReader DMCMemoryReader) *APU {
	const (
		samplingFrequency uint = 1_789_772
		downSamplingRate  uint = 44100
	)
	return &APU{
		sampleRate:      samplingFrequency / downSamplingRate,
		framePeriod:     7458,
		audio:           audio,
		dmcMemoryReader: dmcMemoryReader}
}

type Audio interface {
	write(float32)
}

func (a *APU) frameSequenceMode() frameSequenceMode {
	if util.IsSet(a.frameCounterControl, 7) {
		return frameSequenceMode4Step
	}
	return frameSequenceMode5Step
}

func (a *APU) frameInterruptInhibit() bool { return util.IsSet(a.frameCounterControl, 6) }

func (a *APU) Reset() {
	a.write(0x4017, 0) // frame irq enabled
	a.write(0x4015, 0) // all channels disabled

	for addr := uint16(0x4000); addr <= 0x400F; addr++ {
		a.write(addr, 0)
	}
	for addr := uint16(0x4010); addr <= 0x4013; addr++ {
		a.write(addr, 0)
	}
}

func (a *APU) Step() bool {
	a.cycles += 1

	if a.cycles%a.sampleRate == 0 {
		a.audio.write(a.sample())
	}

	var cpuStall = false
	if a.cycles%2 == 0 {
		a.pulse1.clockTimer()
		a.pulse2.clockTimer()
		a.noise.clockTimer()
		cpuStall = a.dmc.clockTimer(a.dmcMemoryReader)
	}

	a.triangle.clockTimer()

	if a.cycles%a.framePeriod == 0 {
		switch a.frameSequenceMode() {
		case frameSequenceMode4Step:
			a.pulse1.clockEnvelope()
			a.pulse2.clockEnvelope()
			a.triangle.clockLinearCounter()
			a.noise.clockEnvelope()

			if a.frameSequenceStep == 1 || a.frameSequenceStep == 3 {
				a.pulse1.lengthCounter.clock()
				a.pulse1.clockSweepUnit()
				a.pulse2.lengthCounter.clock()
				a.pulse2.clockSweepUnit()
				a.triangle.lengthCounter.clock()
				a.noise.lengthCounter.clock()
			}

			if a.frameSequenceStep == 3 && !a.frameInterruptInhibit() {
				a.frameInterrupted = true
			}

			a.frameSequenceStep = (a.frameSequenceStep + 1) % 4
		case frameSequenceMode5Step:
			if a.frameSequenceStep < 4 || a.frameSequenceStep == 5 {
				a.pulse1.clockEnvelope()
				a.pulse2.clockEnvelope()
				a.triangle.clockLinearCounter()
				a.noise.clockEnvelope()
			}

			if a.frameSequenceStep == 1 || a.frameSequenceStep == 4 {
				a.pulse1.lengthCounter.clock()
				a.pulse1.clockSweepUnit()
				a.pulse2.lengthCounter.clock()
				a.pulse2.clockSweepUnit()
				a.triangle.lengthCounter.clock()
				a.noise.lengthCounter.clock()
			}

			a.frameSequenceStep = (a.frameSequenceStep + 1) % 5
		}

		if a.dmc.interrupted {
			a.frameInterrupted = true
		}
	}

	return cpuStall
}

func (a *APU) sample() float32 {
	p1 := float32(a.pulse1.output())
	p2 := float32(a.pulse2.output())
	triangle := float32(a.triangle.output())
	noise := float32(a.noise.output())
	dmc := float32(a.dmc.output())

	var pulseOut float32
	if p1 != 0.0 || p2 != 0.0 {
		pulseOut = 95.88 / ((8128.0 / (p1 + p2)) + 100.0)
	} else {
		pulseOut = 0.0
	}

	var tndOut float32
	if triangle != 0.0 || noise != 0.0 || dmc != 0.0 {
		tndOut = 159.79 / (1/((triangle/8227)+(noise/12241)+(dmc/22638)) + 100)
		// tndOut = 0.0
	} else {
		tndOut = 0.0
	}

	return pulseOut + tndOut
}

func (a *APU) read(addr uint16) uint8 {
	switch addr {
	case 0x4015:
		var value uint8
		if a.dmc.interrupted {
			value |= 0x80
		}
		if a.frameInterrupted && !a.frameInterruptInhibit() {
			value |= 0x40
		}
		if 0 < a.dmc.bytesRemainingCounter {
			value |= 0x20
		}
		if 0 < a.noise.lengthCounter.count {
			value |= 0x08
		}
		if 0 < a.triangle.lengthCounter.count {
			value |= 0x04
		}
		if 0 < a.pulse2.lengthCounter.count {
			value |= 0x02
		}
		if 0 < a.pulse1.lengthCounter.count {
			value |= 0x01
		}

		a.frameInterrupted = false

		return value
	default:
		return 0x00
	}
}

func (a *APU) write(addr uint16, value uint8) {
	switch {
	case addr <= 0x4000 && addr <= 0x4003:
		a.pulse1.write(addr, value)

	case addr <= 0x4004 && addr <= 0x4007:
		a.pulse2.write(addr, value)

	case addr <= 0x4008 && addr <= 0x400B:
		a.triangle.write(addr, value)

	case addr <= 0x400C && addr <= 0x400F:
		a.noise.write(addr, value)

	case addr <= 0x4010 && addr <= 0x4013:
		a.dmc.write(addr, value)

	case addr == 0x4015:
		a.pulse1.lengthCounter.enable(util.IsSet(value, 0))
		a.pulse2.lengthCounter.enable(util.IsSet(value, 1))
		a.triangle.lengthCounter.enable(util.IsSet(value, 2))
		a.noise.lengthCounter.enable(util.IsSet(value, 3))
		a.dmc.enabled = util.IsSet(value, 4)
	case addr == 0x4017:
		a.frameCounterControl = value
	default:
		break
	}
}

type frameSequenceMode int

const (
	_ frameSequenceMode = iota
	frameSequenceMode4Step
	frameSequenceMode5Step
)
