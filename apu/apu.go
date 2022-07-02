package apu

import (
	"fmt"

	"github.com/thara/gorones/util"
)

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

	audio AudioRenderer

	Port Port
}

func New(audio AudioRenderer) *APU {
	const (
		samplingFrequency uint = 1_789_772
		downSamplingRate  uint = 44100
	)
	apu := &APU{
		sampleRate:  samplingFrequency / downSamplingRate,
		framePeriod: 7458,
		audio:       audio,
		pulse1:      pulseChannel{carryMode: sweepOneComplement},
		pulse2:      pulseChannel{carryMode: sweepTwoComplement},
		noise:       noiseChannel{shiftRegister: 1},
	}
	apu.Port = Port{apu}
	return apu
}

type AudioRenderer interface {
	Write(float32)
}

func (a *APU) frameSequenceMode() frameSequenceMode {
	if util.IsSet(a.frameCounterControl, 7) {
		return frameSequenceMode5Step
	}
	return frameSequenceMode4Step
}

func (a *APU) frameInterruptInhibit() bool { return util.IsSet(a.frameCounterControl, 6) }

func (a *APU) Step(dmcMemoryReader DMCMemoryReader) bool {
	a.cycles += 1

	if a.cycles%a.sampleRate == 0 {
		a.audio.Write(a.sample())
	}

	var cpuStall = false
	if a.cycles%2 == 0 {
		a.pulse1.clockTimer()
		a.pulse2.clockTimer()
		a.noise.clockTimer()
		cpuStall = a.dmc.clockTimer(dmcMemoryReader)
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
				a.pulse1.clockLengthCounter()
				a.pulse1.clockSweepUnit()
				a.pulse2.clockLengthCounter()
				a.pulse2.clockSweepUnit()
				a.triangle.clockLengthCounter()
				a.noise.clockLengthCounter()
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
				a.pulse1.clockLengthCounter()
				a.pulse1.clockSweepUnit()
				a.pulse2.clockLengthCounter()
				a.pulse2.clockSweepUnit()
				a.triangle.clockLengthCounter()
				a.noise.clockLengthCounter()
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
	// if 0 < p1 {
	// 	fmt.Print(p1)
	// }
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
	} else {
		tndOut = 0.0
	}

	return pulseOut + tndOut
}

type Port struct {
	apu *APU
}

func (a *Port) Reset() {
	a.Write(0x4017, 0) // frame irq enabled
	a.Write(0x4015, 0) // all channels disabled

	for addr := uint16(0x4000); addr <= 0x400F; addr++ {
		a.Write(addr, 0)
	}
	for addr := uint16(0x4010); addr <= 0x4013; addr++ {
		a.Write(addr, 0)
	}
}

func (a *Port) Read(addr uint16) uint8 {
	switch addr {
	case 0x4015:
		var value uint8
		if a.apu.dmc.interrupted {
			value |= 0x80
		}
		if a.apu.frameInterrupted && !a.apu.frameInterruptInhibit() {
			value |= 0x40
		}
		if 0 < a.apu.dmc.bytesRemainingCounter {
			value |= 0x20
		}
		if 0 < a.apu.noise.lengthCounter {
			value |= 0x08
		}
		if 0 < a.apu.triangle.lengthCounter {
			value |= 0x04
		}
		if 0 < a.apu.pulse2.lengthCounter {
			value |= 0x02
		}
		if 0 < a.apu.pulse1.lengthCounter {
			value |= 0x01
		}

		a.apu.frameInterrupted = false

		return value
	default:
		return 0x00
	}
}

func (a *Port) Write(addr uint16, value uint8) {
	fmt.Printf("write %04x, %08b\n", addr, value)
	switch {
	case 0x4000 <= addr && addr <= 0x4003:
		a.apu.pulse1.write(addr, value)

	case 0x4004 <= addr && addr <= 0x4007:
		a.apu.pulse2.write(addr, value)

	case 0x4008 <= addr && addr <= 0x400B:
		a.apu.triangle.write(addr, value)

	case 0x400C <= addr && addr <= 0x400F:
		a.apu.noise.write(addr, value)

	case 0x4010 <= addr && addr <= 0x4013:
		a.apu.dmc.write(addr, value)

	case addr == 0x4015:
		a.apu.pulse1.setEnabled(value&1 == 1)
		a.apu.pulse2.setEnabled(value&2 == 2)
		a.apu.triangle.setEnabled(value&4 == 4)
		a.apu.noise.setEnabled(value&8 == 8)

		a.apu.dmc.enabled = value&16 == 16
	case addr == 0x4017:
		a.apu.frameCounterControl = value
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

var lengthTable = [32]uint8{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}
