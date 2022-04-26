package ppu

import "github.com/thara/gorones/mapper"

const spriteCount = 64

type Emu struct {
	ppu PPU

	nt       [0x1000]uint8
	palletes [0x0020]uint8
	oam      [4 * spriteCount]uint8

	cpuDataBus uint8

	scan scan

	mapper    mapper.Mapper
	mirroring mapper.Mirroring
}

func NewEmu(mapper mapper.Mapper, port *Port) *Emu {
	e := &Emu{
		mapper:    mapper,
		mirroring: mapper.Mirroring(),
	}
	port.emu = e
	return e
}

func (e *Emu) read(addr vramAddr) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return e.mapper.Read(addr)
	case 0x2000 <= addr && addr <= 0x2FFF:
		return e.nt[toNTAddr(addr, e.mirroring)]
	case 0x3000 <= addr && addr <= 0x3EFF:
		return e.nt[toNTAddr(addr-0x1000, e.mirroring)]
	case 0x3F00 <= addr && addr <= 0x3FFF:
		return e.palletes[toPalleteAddr(addr)]
	default:
		return 0
	}
}

func (e *Emu) write(addr vramAddr, value uint8) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		e.mapper.Write(addr, value)
	case 0x2000 <= addr && addr <= 0x2FFF:
		e.nt[toNTAddr(addr, e.mirroring)] = value
	case 0x3000 <= addr && addr <= 0x3EFF:
		e.nt[toNTAddr(addr-0x1000, e.mirroring)] = value
	case 0x3F00 <= addr && addr <= 0x3FFF:
		e.palletes[toPalleteAddr(addr)] = value
	}
}

func toNTAddr(addr uint16, m mapper.Mirroring) uint16 {
	switch m {
	case mapper.Mirroring_Horizontal:
		if 0x2800 <= addr {
			return (0x0800 & addr) % 0x0400
		} else {
			return addr % 0x0400
		}
	case mapper.Mirroring_Vertical:
		return addr % 0x0800
	}
	return addr - 0x2000
}

func toPalleteAddr(addr uint16) uint16 {
	// http://wiki.nesdev.com/w/index.php/PPU_palettes#Memory_Map
	a := addr % 32
	if a%4 == 0 {
		return (a | 0x10)
	}
	return a
}
