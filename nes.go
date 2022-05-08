package gorones

import (
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/input"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/ppu"
)

// NES components
type NES struct {
	cpu *cpu.CPU
	ppu *ppu.PPU

	interrupt *cpu.Interrupt
}

func NewNES(m mapper.Mapper, ctrl1, ctrl2 input.Controller, renderer ppu.FrameRenderer) *NES {
	intr := cpu.NoInterrupt

	ppu := ppu.New(m, renderer)
	t := cpuTicker{ppu: ppu, interrupt: &intr}
	ctx := cpuBus{mapper: m, ppuPort: ppu.Port, ctrl1: ctrl1, ctrl2: ctrl2, t: &t}
	return &NES{
		cpu:       cpu.New(&t, &ctx),
		ppu:       ppu,
		interrupt: &intr,
	}
}

func (n *NES) PowerOn() {
	n.cpu.PowerOn()
}

func (n *NES) Reset() {
	n.cpu.Reset()
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

type cpuTicker struct {
	ppu *ppu.PPU

	cycles uint64

	interrupt *cpu.Interrupt
}

func (c *cpuTicker) Tick() {
	c.cycles += 1
	// 3 PPU cycles per 1 CPU cycle
	c.ppu.Step(c.interrupt)
	c.ppu.Step(c.interrupt)
	c.ppu.Step(c.interrupt)
}

// https://www.nesdev.org/wiki/CPU_memory_map
type cpuBus struct {
	wram   [0x0800]uint8
	mapper mapper.Mapper

	ppuPort ppu.Port

	ctrl1, ctrl2 input.Controller

	t *cpuTicker
}

func (b *cpuBus) ReadCPU(addr uint16) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return b.wram[addr%0x0800]
	case 0x2000 <= addr && addr <= 0x3FFF:
		return b.ppuPort.ReadRegister(ppuAddr(addr))
	case addr == 0x4016:
		return b.ctrl1.Read()
	case addr == 0x4017:
		return b.ctrl2.Read()
	case 0x4020 <= addr && addr <= 0xFFFF:
		return b.mapper.Read(addr)
	}
	return 0
}

func (b *cpuBus) WriteCPU(addr uint16, value uint8) {
	// OAMDMA https://wiki.nesdev.org/w/index.php?title=PPU_registers#OAM_DMA_.28.244014.29_.3E_write
	if addr == 0x4014 {
		start := uint16(value << 2)
		for addr := start; addr <= start+0xFF; addr++ {
			m := b.ReadCPU(addr)
			b.WriteCPU(0x2004, m)
		}
		// dummy cycles
		b.t.Tick()
		if b.t.cycles%2 == 1 {
			b.t.Tick()
		}
		return
	}

	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		b.wram[addr%0x0800] = value
	case 0x2000 <= addr && addr <= 0x3FFF:
		b.ppuPort.WriteRegister(ppuAddr(addr), value)
	case addr == 0x4016:
		b.ctrl1.Write(value)
	case addr == 0x4017:
		b.ctrl2.Write(value)
		// TODO apu write
	case 0x4020 <= addr && addr <= 0xFFFF:
		b.mapper.Write(addr, value)
	}
}

func ppuAddr(addr uint16) uint16 {
	// repears every 8 bytes
	return 0x2000 + addr%8

}
