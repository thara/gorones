package ppu

type Port struct {
	emu *Emu
}

func (p *Port) ReadRegister(addr uint16) uint8 {
	e := p.emu

	var result uint8
	switch addr {
	case 0x2002: // PPUSTATUS
		result = e.ppu.readStatus() | (e.cpuDataBus & 0b11111)
		e.ppu.status.vblank = false
		e.ppu.w = false
	case 0x2004: // OAMDATA
		result = e.oam[e.ppu.oamAddr]
	case 0x2007:
		// https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		if e.ppu.v <= 0x3EFF {
			result = e.ppu.data
			e.ppu.data = e.read(e.ppu.v)
		} else {
			result = e.read(e.ppu.v)
		}
		if e.ppu.ctrl.vramIncr {
			e.ppu.v += 1
		} else {
			e.ppu.v += 32
		}
	default:
		result = e.cpuDataBus
	}
	e.cpuDataBus = result
	return result
}

func (p *Port) WriteRegister(addr uint16, value uint8) {
	e := p.emu

	switch addr {
	case 0x2000:
		e.ppu.setController(value)
		// t: ...BA.. ........ = d: ......BA
		e.ppu.t = e.ppu.t&^uint16(0b110000000000) | uint16(e.ppu.ctrl.nt)<<10
	case 0x2001:
		e.ppu.setMask(value)
	case 0x2003:
		e.ppu.oamAddr = value
	case 0x2004:
		e.oam[e.ppu.oamAddr] = value
		e.ppu.oamAddr += 1
	case 0x2005:
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242005_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242005_second_write_.28w_is_1.29
		//TODO
	case 0x2006:
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242006_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242006_second_write_.28w_is_1.29
		//TODO
	case 0x2007:
		e.write(e.ppu.v, value)
	}
}
