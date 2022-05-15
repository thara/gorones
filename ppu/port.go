package ppu

type Port struct {
	ppu *PPU
}

func (p *Port) ReadRegister(addr uint16) uint8 {
	var result uint8
	switch addr {
	case 0x2002: // PPUSTATUS
		result = p.ppu.readStatus() | (p.ppu.cpuDataBus & 0b11111)
		p.ppu.status.vblank = false
		p.ppu.w = false
		// race condition
		if p.ppu.scan.line == 241 && p.ppu.scan.dot < 2 {
			result &^= 0x80
		}
	case 0x2004: // OAMDATA
		// https://wiki.nesdev.com/w/index.php/PPU_sprite_evaluation
		if p.ppu.scan.line < 240 && 1 <= p.ppu.scan.dot && p.ppu.scan.dot <= 64 {
			// during sprite evaluation
			result = 0xFF
		} else {
			result = p.ppu.spr.oam[p.ppu.oamAddr]
		}
	case 0x2007: // PPUDATA
		// https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		if p.ppu.v <= 0x3EFF {
			result = p.ppu.data
			p.ppu.data = p.ppu.read(p.ppu.v)
		} else {
			result = p.ppu.read(p.ppu.v)
		}
		incrV(p.ppu)
	default:
		result = p.ppu.cpuDataBus
	}
	p.ppu.cpuDataBus = result
	return result
}

func (p *Port) WriteRegister(addr uint16, value uint8) {
	switch addr {
	case 0x2000: // PPUCTRL
		p.ppu.setController(value)
		// t: ...BA.. ........ = d: ......BA
		p.ppu.t = p.ppu.t&^uint16(0b110000000000) | uint16(p.ppu.ctrl.nt)<<10
	case 0x2001: // PPUMASK
		p.ppu.setMask(value)
	case 0x2003: // OAMADDR
		p.ppu.oamAddr = value
	case 0x2004: // OAMDATA
		p.ppu.spr.oam[p.ppu.oamAddr] = value
		p.ppu.oamAddr++
	case 0x2005: // PPUSCROLL
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242005_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242005_second_write_.28w_is_1.29
		d := uint16(value)
		if !p.ppu.w {
			// first write
			// t: ....... ...HGFED = d: HGFED...
			// x:              CBA = d: .....CBA
			p.ppu.t = (p.ppu.t &^ 0b11111) | (d&0b11111000)>>3
			p.ppu.x = value & 0b111
		} else {
			// second write
			// t: CBA..HG FED..... = d: HGFEDCBA
			p.ppu.t = (p.ppu.t &^ 0b111001111100000) | ((d & 0b111) << 12) | ((d & 0b11111000) << 2)
		}
		p.ppu.w = !p.ppu.w
	case 0x2006: // PPUADDR
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242006_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php/PPU_scrolling#.242006_second_write_.28w_is_1.29
		d := uint16(value)
		if !p.ppu.w {
			// first write
			// t: .FEDCBA ........ = d: ..FEDCBA
			// t: X...... ........ = 0
			p.ppu.t = (p.ppu.t &^ 0b011111100000000) | ((d & 0b111111) << 8)
		} else {
			// second write
			// t: ....... HGFEDCBA = d: HGFEDCBA
			// v                   = t
			p.ppu.t = (p.ppu.t &^ 0b11111111) | d
			p.ppu.v = p.ppu.t
		}
		p.ppu.w = !p.ppu.w
	case 0x2007: // PPUDATA
		p.ppu.write(p.ppu.v, value)
		incrV(p.ppu)
	}
}

func incrV(p *PPU) {
	if p.ctrl.vramIncr {
		p.v += 32
	} else {
		p.v += 1
	}
}
