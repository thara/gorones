package gorones

import (
	"github.com/thara/gorones/apu"
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/input"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/ppu"
)

// NES components
type NES struct {
	cpu *cpu.CPU
	ppu *ppu.PPU
	apu *apu.APU

	interrupt *cpu.Interrupt

	cycles uint64

	wram   [0x0800]uint8
	mapper mapper.Mapper

	ctrl1, ctrl2 input.Controller
}

func NewNES(m mapper.Mapper, ctrl1, ctrl2 input.Controller, frameRenderer ppu.FrameRenderer, audioRenderer apu.AudioRenderer) *NES {
	intr := cpu.NoInterrupt

	nes := &NES{interrupt: &intr, mapper: m, ctrl1: ctrl1, ctrl2: ctrl2}
	nes.cpu = cpu.New(nes, nes)
	nes.ppu = ppu.New(m, frameRenderer)
	nes.apu = apu.New(audioRenderer)
	return nes
}

func (n *NES) PowerOn() {
	n.cpu.PowerOn()
}

func (n *NES) Reset() {
	n.cpu.Reset()
	n.apu.Reset()
}

func (n *NES) InitNEStest() {
	n.cpu.PC = 0xC000
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_ref-1
	n.cpu.P.Set(0x24)
	n.cpu.Cycles = 7
}

func (n *NES) RunFrame() {
	before := n.ppu.CurrentFrames()
	for before == n.ppu.CurrentFrames() {
		n.step()
	}
}

func (n *NES) step() {
	n.cpu.Step(n.interrupt)
}

func (n *NES) Tick() {
	n.cycles += 1

	cpuSteel := n.apu.Step(n)
	if cpuSteel {
		n.cycles += 4
	}

	// 3 PPU cycles per 1 CPU cycle
	n.ppu.Step(n.interrupt)
	n.ppu.Step(n.interrupt)
	n.ppu.Step(n.interrupt)
}

// https://www.nesdev.org/wiki/CPU_memory_map

func (b *NES) ReadCPU(addr uint16) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return b.wram[addr%0x0800]
	case 0x2000 <= addr && addr <= 0x3FFF:
		return b.ppu.ReadRegister(ppuAddr(addr))

	case 0x4000 <= addr && addr <= 0x4013:
		fallthrough
	case addr == 0x4015:
		return b.apu.Read(addr)

	case addr == 0x4016:
		return b.ctrl1.Read()
	case addr == 0x4017:
		return b.ctrl2.Read()
	case 0x4020 <= addr && addr <= 0xFFFF:
		return b.mapper.Read(addr)
	}
	return 0
}

func (b *NES) WriteCPU(addr uint16, value uint8) {
	// OAMDMA https://wiki.nesdev.org/w/index.php?title=PPU_registers#OAM_DMA_.28.244014.29_.3E_write
	if addr == 0x4014 {
		start := uint16(value) * 0x100
		for i := uint16(0); i <= 0xFF; i++ {
			m := b.ReadCPU(start + i)
			b.WriteCPU(0x2004, m)
		}
		// dummy cycles
		b.Tick()
		if b.cycles%2 == 1 {
			b.Tick()
		}
		return
	}

	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		b.wram[addr%0x0800] = value
	case 0x2000 <= addr && addr <= 0x3FFF:
		b.ppu.WriteRegister(ppuAddr(addr), value)

	case 0x4000 <= addr && addr <= 0x4013:
		fallthrough
	case addr == 0x4015:
		b.apu.Write(addr, value)

	case addr == 0x4016:
		b.ctrl1.Write(value)
	case addr == 0x4017:
		b.ctrl2.Write(value)
		b.apu.Write(addr, value)
	case 0x4020 <= addr && addr <= 0xFFFF:
		b.mapper.Write(addr, value)
	}
}

// adapt to apu.DMCMemoryReader
func (b *NES) Read(addr uint16) uint8 {
	return b.ReadCPU(addr)
}

func ppuAddr(addr uint16) uint16 {
	// repeats every 8 bytes
	return 0x2000 + addr%8

}
