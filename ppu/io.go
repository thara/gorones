package ppu

import "github.com/thara/gorones/mapper"

func (p *PPU) read(addr uint16) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return p.mapper.Read(addr)
	case 0x2000 <= addr && addr <= 0x2FFF:
		return p.nt[toNTAddr(addr, p.mirroring)]
	case 0x3000 <= addr && addr <= 0x3EFF:
		return p.nt[toNTAddr(addr-0x1000, p.mirroring)]
	case 0x3F00 <= addr && addr <= 0x3FFF:
		return p.palletes[toPalleteAddr(addr)]
	default:
		return 0
	}
}

func (p *PPU) write(addr uint16, value uint8) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		p.mapper.Write(addr, value)
	case 0x2000 <= addr && addr <= 0x2FFF:
		p.nt[toNTAddr(addr, p.mirroring)] = value
	case 0x3000 <= addr && addr <= 0x3EFF:
		p.nt[toNTAddr(addr-0x1000, p.mirroring)] = value
	case 0x3F00 <= addr && addr <= 0x3FFF:
		p.palletes[toPalleteAddr(addr)] = value
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
