package ppu

func (p *PPU) ReadRegister(addr uint16) uint8 {
	var result uint8
	switch addr {
	case 0x2002: // PPUSTATUS
		result = p.readStatus() | (p.cpuDataBus & 0b11111)
		p.status.vblank = false
		p.w = false
		// race condition
		if p.scan.line == 241 && p.scan.dot < 2 {
			result &^= 0x80
		}
	case 0x2004: // OAMDATA
		// https://wiki.nesdev.com/w/index.php_sprite_evaluation
		if p.scan.line < 240 && 1 <= p.scan.dot && p.scan.dot <= 64 {
			// during sprite evaluation
			result = 0xFF
		} else {
			result = p.spr.oam[p.oamAddr]
		}
	case 0x2007: // PPUDATA
		// https://www.nesdev.org/wiki/PPU_registers#The_PPUDATA_read_buffer_(post-fetch)
		if p.v <= 0x3EFF {
			result = p.data
			p.data = p.read(p.v)
		} else {
			result = p.read(p.v)
		}
		incrV(p)
	default:
		result = p.cpuDataBus
	}
	p.cpuDataBus = result
	return result
}

func (p *PPU) WriteRegister(addr uint16, value uint8) {
	switch addr {
	case 0x2000: // PPUCTRL
		p.setController(value)
		// t: ...BA.. ........ = d: ......BA
		p.t = p.t&^uint16(0b110000000000) | uint16(p.ctrl.nt)<<10
	case 0x2001: // PPUMASK
		p.setMask(value)
	case 0x2003: // OAMADDR
		p.oamAddr = value
	case 0x2004: // OAMDATA
		p.spr.oam[p.oamAddr] = value
		p.oamAddr++
	case 0x2005: // PPUSCROLL
		// http://wiki.nesdev.org/w/index.php_scrolling#.242005_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php_scrolling#.242005_second_write_.28w_is_1.29
		d := uint16(value)
		if !p.w {
			// first write
			// t: ....... ...HGFED = d: HGFED...
			// x:              CBA = d: .....CBA
			p.t = (p.t &^ 0b11111) | (d&0b11111000)>>3
			p.x = value & 0b111
		} else {
			// second write
			// t: CBA..HG FED..... = d: HGFEDCBA
			p.t = (p.t &^ 0b111001111100000) | ((d & 0b111) << 12) | ((d & 0b11111000) << 2)
		}
		p.w = !p.w
	case 0x2006: // PPUADDR
		// http://wiki.nesdev.org/w/index.php_scrolling#.242006_first_write_.28w_is_0.29
		// http://wiki.nesdev.org/w/index.php_scrolling#.242006_second_write_.28w_is_1.29
		d := uint16(value)
		if !p.w {
			// first write
			// t: .FEDCBA ........ = d: ..FEDCBA
			// t: X...... ........ = 0
			p.t = (p.t &^ 0b011111100000000) | ((d & 0b111111) << 8)
		} else {
			// second write
			// t: ....... HGFEDCBA = d: HGFEDCBA
			// v                   = t
			p.t = (p.t &^ 0b11111111) | d
			p.v = p.t
		}
		p.w = !p.w
	case 0x2007: // PPUDATA
		p.write(p.v, value)
		incrV(p)
	}
}

func incrV(p *PPU) {
	if p.ctrl.vramIncr {
		p.v += 32
	} else {
		p.v += 1
	}
}
