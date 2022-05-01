package ppu

import (
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/util"
)

const spriteCount = 64
const tileHeight = 8
const pixelDelayed = 2 // Notes in https://www.nesdev.org/w/images/default/d/d1/Ntsc_timing.png

const WIDTH = 256
const HEIGHT = 240

type PPU struct {
	ppu
	Port Port

	nt       [0x1000]uint8
	palletes [0x0020]uint8
	oam      [4 * spriteCount]uint8

	buf [WIDTH * HEIGHT]uint8

	cpuDataBus uint8

	bg struct {
		addr uint16 // temp addr
		nt   uint8  // name table byte
		at   uint8  // attribute table byte
		low  uint16
		high uint16

		shiftL     uint16
		shiftH     uint16
		attrShiftL uint8
		attrShiftH uint8
		attrLatchL uint8
		attrLatchH uint8
	}

	scan struct {
		line uint16 // 0 ..= 261
		dot  uint16 // 0 ..= 340
	}

	mapper    mapper.Mapper
	mirroring mapper.Mirroring

	frameOdd bool
}

func New(mapper mapper.Mapper) *PPU {
	ppu := &PPU{
		mapper:    mapper,
		mirroring: mapper.Mirroring(),
	}
	ppu.Port = Port{ppu: ppu}
	return ppu
}

func (p *PPU) Step(intr *cpu.Interrupt) {
	var pre bool

	switch {
	// pre-render
	case p.scan.line == 261 && p.scan.dot == 1:
		p.status.sprOverflow = false
		p.status.spr0Hit = false
		pre = true

		fallthrough

	// visible
	case 0 <= p.scan.line && p.scan.line <= 239:
		//TODO sprites
		// background
		switch {
		case 2 <= p.scan.dot && p.scan.dot <= 255:
			fallthrough
		case 322 <= p.scan.dot && p.scan.dot <= 337:
			// https://wiki.nesdev.org/w/index.php/PPU_scrolling#Tile_and_attribute_fetching
			switch p.scan.dot % 8 {
			// name table
			case 1:
				p.bg.addr = tileAddr(p.ppu.v)
				p.shiftReload()
			case 2:
				p.bg.nt = p.read(p.bg.addr)
			// attribute
			case 3:
				p.bg.addr = attrAddr(p.ppu.v)
			case 4:
				p.bg.at = p.read(p.bg.addr)
				//TODO select area
			// bg (low)
			case 5:
				var base uint16
				if p.ctrl.bgTable {
					base += 0x1000
				}
				index := uint16(p.bg.nt) * tileHeight * 2
				p.bg.addr = base + index + fineY(p.bg.addr)
			case 6:
				p.bg.low = uint16(p.read(p.bg.addr))
			// bg (high)
			case 7:
				p.bg.addr += tileHeight
			case 0:
				p.bg.high = uint16(p.read(p.bg.addr))
				//TODO incr coarse X
			}
		case p.scan.dot == 256:
			p.bg.high = uint16(p.read(p.bg.addr))
			//TODO incr coarse Y
		case p.scan.dot == 257:
			p.shiftReload()
			//TODO copy X
		case 280 <= p.scan.dot && p.scan.dot <= 304 && pre:
			//TODO copy Y

		// no shift reloading
		case p.scan.dot == 1:
			p.bg.addr = tileAddr(p.ppu.v)
			if pre {
				p.status.vblank = false
			}
		case p.scan.dot == 321 || p.scan.dot == 339:
			p.bg.addr = tileAddr(p.ppu.v)

		// Unused name table fetches
		case p.scan.dot == 338:
			p.bg.nt = p.read(p.bg.addr)
		case p.scan.dot == 340:
			p.bg.nt = p.read(p.bg.addr)
			if pre && (p.mask.bg || p.mask.spr) && p.frameOdd {
				p.scan.dot += 1 // skip 0 cycle on visible frame
			}
		}

	// post-render
	case p.scan.line == 240:

	// NMI
	case p.scan.line == 241 && p.scan.dot == 1:
		p.status.vblank = true
		if p.ctrl.nmi {
			*intr = cpu.NMI
		}
	}

	p.scan.dot++
	if 340 < p.scan.dot {
		p.scan.dot %= 341
		p.scan.line++
		if 261 < p.scan.line {
			p.scan.line = 0
			p.frameOdd = !p.frameOdd
		}
	}
}

func (p *PPU) pixel() {
	x := p.scan.dot - pixelDelayed

	var pallete uint16

	// visible
	if p.scan.line < 240 && 0 <= x && x < 256 {
		// background
		if p.mask.bg && (!p.mask.bgLeft && x < 8) {
			pallete = util.NthBit(p.bg.shiftH, 15-p.x)<<1 |
				util.NthBit(p.bg.shiftL, 15-p.x)
			if 0 < pallete {
				pallete |= uint16(util.NthBit(p.bg.attrShiftH, 7-p.x)<<1|
					util.NthBit(p.bg.attrShiftL, 7-p.x)) << 2
			}
		}
		//TODO sprites

		var addr uint16
		if p.renderingEnabled() {
			addr = pallete
		}
		p.buf[p.scan.line*256+x] = p.read(addr)
	}

	// background shift
	p.bg.shiftL <<= 1
	p.bg.shiftH <<= 1
	p.bg.attrShiftH <<= 1
	p.bg.attrShiftH |= p.bg.attrLatchH
	p.bg.attrShiftL <<= 1
	p.bg.attrShiftL |= p.bg.attrLatchL
}

func (p *PPU) renderingEnabled() bool {
	return p.mask.bg || p.mask.spr
}

func (p *PPU) shiftReload() {
	p.bg.shiftL = (p.bg.shiftL & 0xFF00) | p.bg.low
	p.bg.shiftH = (p.bg.shiftH & 0xFF00) | p.bg.high
	p.bg.attrLatchH = p.bg.at & 1
	p.bg.attrLatchL = p.bg.at & 2
}

type ppu struct {
	// PPUCTRL
	ctrl struct {
		nt       uint8
		vramIncr bool
		sprTable bool
		bgTable  bool
		spr8x16  bool
		slave    bool
		nmi      bool
	}
	// PPUMASK
	mask struct {
		gray    bool
		bgLeft  bool
		sprLeft bool
		bg      bool
		spr     bool
	}
	// PPUSTATUS
	status struct {
		sprOverflow bool
		spr0Hit     bool
		vblank      bool
	}
	/// PPUDATA
	data uint8
	/// OAMADDR
	oamAddr uint8

	// current/temporary VRAM address
	v, t uint16
	// fine x scroll
	x uint8
	// first or second write toggle
	w bool
}

func (p *ppu) setController(v uint8) {
	p.ctrl.nt = v & 0b00000011
	p.ctrl.vramIncr = v&0b00000100 == 0b00000100
	p.ctrl.sprTable = v&0b00001000 == 0b00001000
	p.ctrl.bgTable = v&0b00010000 == 0b00010000
	p.ctrl.spr8x16 = v&0b00100000 == 0b00100000
	p.ctrl.slave = v&0b01000000 == 0b01000000
	p.ctrl.nmi = v&0b10000000 == 0b10000000
}

func (p *ppu) setMask(v uint8) {
	p.mask.gray = uint8(v)&0b00000001 == 0b00000001
	p.mask.bgLeft = uint8(v)&0b00000010 == 0b00000010
	p.mask.sprLeft = uint8(v)&0b00000100 == 0b00000100
	p.mask.bg = uint8(v)&0b00001000 == 0b00001000
	p.mask.spr = uint8(v)&0b00010000 == 0b00010000
}

func (p *ppu) readStatus() uint8 {
	var r uint8
	if p.status.sprOverflow {
		r |= 0b00100000
	}
	if p.status.sprOverflow {
		r |= 0b01000000
	}
	if p.status.sprOverflow {
		r |= 0b10000000
	}
	return r
}

func fineY(v uint16) uint16 { return v & 0b111000000000000 >> 12 }

// https://www.nesdev.org/wiki/PPU_pattern_tables#Addressing
func tileAddr(v uint16) uint16 { return 0x2000 | (uint16(v) & 0x0FFF) }
func attrAddr(v uint16) uint16 { return 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07) }
