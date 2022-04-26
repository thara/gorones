package gorones

import (
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/input"
	"github.com/thara/gorones/mapper"
)

// NES components
type NES struct {
	cpu *cpu.Emu
}

func NewNES(m mapper.Mapper, ctrl1, ctrl2 input.Controller) *NES {
	t := cpuTicker{}
	ctx := cpuBus{mapper: m, ctrl1: ctrl1, ctrl2: ctrl2, t: &t}
	return &NES{
		cpu: cpu.NewEmu(&t, &ctx),
	}
}

func (n *NES) powerOn() {
	n.cpu.PowerOn()
}

func (n *NES) step() {
	n.cpu.Step()
}

type cpuTicker struct {
	cycles uint64
}

func (c *cpuTicker) Tick() {
	c.cycles += 1
}

// https://www.nesdev.org/wiki/CPU_memory_map
type cpuBus struct {
	wram   [0x0800]uint8
	mapper mapper.Mapper

	ctrl1, ctrl2 input.Controller

	t *cpuTicker
}

func (b *cpuBus) ReadCPU(addr uint16) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x07FF:
		return b.wram[addr]
	case 0x0800 <= addr && addr <= 0x0FFF:
		// mirrors of 0x0000-0x07FF
		return b.wram[addr-(addr%0x0800*0x800)]
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
	case 0x0000 <= addr && addr <= 0x07FF:
		b.wram[addr] = value
	case 0x0800 <= addr && addr <= 0x0FFF:
		// mirrors of 0x0000-0x07FF
		b.wram[addr-(addr%0x0800*0x800)] = value
	case addr == 0x4016:
		b.ctrl1.Write(value)
	case addr == 0x4017:
		b.ctrl2.Write(value)
		// TODO apu write
	case 0x4020 <= addr && addr <= 0xFFFF:
		b.mapper.Write(addr, value)
	}
}