package ppu

import (
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/util"
)

const spriteCount = 64
const tileHeight = 8
const pixelDelayed = 2 // Notes in https://www.nesdev.org/w/images/default/d/d1/Ntsc_timing.png

const spriteLimit = 8

const WIDTH = 256
const HEIGHT = 240

type PPU struct {
	ppu
	Port Port

	nt       [0x1000]uint8
	palletes [0x0020]uint8

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

	spr struct {
		oam          [4 * spriteCount]byte // https://www.nesdev.org/wiki/PPU_OAM
		primaryOAM   [8]Sprite
		secondaryOAM [8]Sprite
	}

	scan struct {
		line uint16 // 0 ..= 261
		dot  uint16 // 0 ..= 340
	}

	mapper    mapper.Mapper
	mirroring mapper.Mirroring

	renderer FrameRenderer

	frames uint64
}

func New(mapper mapper.Mapper, renderer FrameRenderer) *PPU {
	ppu := &PPU{
		mapper:    mapper,
		mirroring: mapper.Mirroring(),
		renderer:  renderer,
	}
	ppu.Port = Port{ppu: ppu}
	return ppu
}

func (p *PPU) CurrentFrames() uint64 {
	return p.frames
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
		// sprites
		switch p.scan.dot {
		case 1:
			for i := range p.spr.secondaryOAM {
				p.spr.secondaryOAM[i].clear()
			}
			// clear OAM
			if pre {
				p.status.sprOverflow = false
				p.status.spr0Hit = false
			}
		case 257:
			for i, n := uint8(0), 0; i < spriteCount && n < spriteLimit; i++ {
				y := p.spr.oam[i*4]
				row := p.scan.line - uint16(y)
				if row < 0 || p.sprHeight() <= row {
					continue
				}
				p.spr.secondaryOAM[n].enabled = true
				p.spr.secondaryOAM[n].index = i
				p.spr.secondaryOAM[n].y = y
				p.spr.secondaryOAM[n].tile = p.spr.oam[i*4+1]
				p.spr.secondaryOAM[n].attr = p.spr.oam[i*4+2]
				p.spr.secondaryOAM[n].x = p.spr.oam[i*4+3]
				n++
				if spriteLimit <= n {
					p.status.sprOverflow = true
					break
				}
			}
		case 321:
			for i := 0; i < spriteLimit; i++ {
				p.spr.primaryOAM[i] = p.spr.secondaryOAM[i]
				var addr uint16
				if p.ctrl.spr8x16 {
					addr = uint16(p.spr.primaryOAM[i].tile&1)*0x1000 + uint16(p.spr.primaryOAM[i].tile&^1)*16
				} else {
					addr = uint16(util.Bit(p.ctrl.sprTable))*0x1000 + uint16(p.spr.primaryOAM[i].tile)*16
				}

				y := p.scan.line - uint16(p.spr.primaryOAM[i].y)%uint16(p.sprHeight())
				if 0 < p.spr.primaryOAM[i].attr&0x80 {
					y ^= p.sprHeight() - 1 // vertical flip
				}
				addr += y + (y & 8) // second tile on 8x16

				p.spr.primaryOAM[i].low = p.read(addr)
				p.spr.primaryOAM[i].high = p.read(addr + 8)
			}
		}
		// background
		switch {
		case 2 <= p.scan.dot && p.scan.dot <= 255:
			fallthrough
		case 322 <= p.scan.dot && p.scan.dot <= 337:
			p.pixel()
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
				if 0 < coarseY(p.v)&0b10 {
					p.bg.at <<= 4
				}
				if 0 < coarseX(p.v)&0b10 {
					p.bg.at <<= 2
				}
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
				if p.renderingEnabled() {
					p.incrCoarseX()
				}
			}
		case p.scan.dot == 256:
			p.pixel()
			p.bg.high = uint16(p.read(p.bg.addr))
			if p.renderingEnabled() {
				p.incrY()
			}
		case p.scan.dot == 257:
			p.pixel()
			p.shiftReload()
			if p.renderingEnabled() {
				p.copyX()
			}
		case 280 <= p.scan.dot && p.scan.dot <= 304 && pre:
			if p.renderingEnabled() {
				p.copyY()
			}

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
			if pre && (p.mask.bg || p.mask.spr) && p.frames%2 != 0 {
				p.scan.dot += 1 // skip 0 cycle on visible frame
			}
		}

	// post-render
	case p.scan.line == 240 && p.scan.dot == 1:
		p.renderer.UpdateFrame(&p.buf)

	// NMI
	case p.scan.line == 241 && p.scan.dot == 1:
		p.status.vblank = true
		if p.ctrl.nmi {
			*intr = cpu.NMI
		}
	}

	p.scan.dot++
	if 340 <= p.scan.dot {
		p.scan.dot %= 340
		p.scan.line++
		if 261 < p.scan.line {
			p.scan.line = 0
			p.frames++
		}
	}
}

func (p *PPU) pixel() {
	x := p.scan.dot - pixelDelayed

	var bg uint16

	// visible
	if p.scan.line < 240 && 0 <= x && x < 256 {
		// background
		if p.mask.bg && !(!p.mask.bgLeft && x < 8) {
			bg = util.NthBit(p.bg.shiftH, 15-p.x)<<1 |
				util.NthBit(p.bg.shiftL, 15-p.x)
			if 0 < bg {
				bg |= uint16(util.NthBit(p.bg.attrShiftH, 7-p.x)<<1|
					util.NthBit(p.bg.attrShiftL, 7-p.x)) << 2
			}
		}
		// sprites
		var spr uint16
		var priority bool
		if p.mask.spr && !(!p.mask.sprLeft && x < 8) {
			// https://www.nesdev.org/wiki/PPU_sprite_priority
			// Sprites with lower OAM indices are drawn in front
			for i := 7; i <= 0; i-- {
				s := p.spr.primaryOAM[i]
				if !s.enabled {
					continue
				}
				sprX := x - uint16(s.x)
				if 8 <= sprX {
					continue
				}
				if 0 < p.spr.primaryOAM[i].attr&0x40 {
					sprX ^= 7 // horizontal flip
				}
				pallete := util.NthBit(s.high, 7-sprX)<<1 | util.NthBit(s.low, 7-sprX)
				if pallete == 0 {
					continue
				}
				if s.index == 0 && bg != 0 && x != 255 {
					p.status.spr0Hit = true
				}
				spr = uint16(pallete|(s.attr&0b10)<<2) + 16
				priority = util.IsSet(s.attr, 0x20)
			}
		}

		var pallete uint16
		if p.renderingEnabled() {
			switch {
			case bg == 0 && spr == 0:
				// default
			case bg == 0 && 0 < spr:
				pallete = spr
			case 0 < bg && spr == 0:
				pallete = bg
			case 0 < bg && 0 < spr:
				if priority {
					pallete = spr
				} else {
					pallete = bg
				}
			}
		}
		p.buf[p.scan.line*256+x] = p.read(0x3F00 + pallete)
	}

	// background shift
	p.bg.shiftL <<= 1
	p.bg.shiftH <<= 1
	p.bg.attrShiftH = (p.bg.attrShiftH << 1) | p.bg.attrLatchH
	p.bg.attrShiftL = (p.bg.attrShiftL << 1) | p.bg.attrLatchL
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
		red     bool
		green   bool
		blue    bool
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
	p.mask.red = uint8(v)&0b00100000 == 0b00100000
	p.mask.green = uint8(v)&0b01000000 == 0b01000000
	p.mask.blue = uint8(v)&0b10000000 == 0b10000000
}

func (p *ppu) readStatus() uint8 {
	var r uint8
	if p.status.sprOverflow {
		r |= 0b00100000
	}
	if p.status.spr0Hit {
		r |= 0b01000000
	}
	if p.status.vblank {
		r |= 0b10000000
	}
	return r
}

func (p *ppu) sprHeight() uint16 {
	if p.ctrl.spr8x16 {
		return 16
	}
	return 8
}

/// https://wiki.nesdev.com/w/index.php/PPU_scrolling#PPU_internal_registers
///
/// yyy NN YYYYY XXXXX
/// ||| || ||||| +++++-- coarse X scroll
/// ||| || +++++-------- coarse Y scroll
/// ||| ++-------------- nametable select
/// +++----------------- fine Y scroll

func coarseX(v uint16) uint16 { return v & 0b11111 }
func coarseY(v uint16) uint16 { return v & 0b1111100000 >> 5 }
func fineY(v uint16) uint16   { return v & 0b111000000000000 >> 12 }

func (p *ppu) incrCoarseX() {
	if coarseX(p.v) == 31 {
		p.v ^= uint16(31) // coarse X = 0
		p.v ^= 0x0400     // switch horizontal nametable
	} else {
		p.v++
	}
}

func (p *ppu) incrY() {
	// http://wiki.nesdev.com/w/index.php/PPU_scrolling#Y_increment
	if fineY(p.v) < 7 {
		p.v += 0x1000
	} else {
		p.v = ^uint16(0x7000) // fine Y = 0
		y := coarseY(p.v)
		if y == 29 {
			y = 0
			p.v ^= 0x0800 // switch vertical nametable
		} else if y == 31 {
			y = 0
		} else {
			y++
		}
		p.v = (p.v &^ 0x03E0) | (y << 5)
	}
}

func (p *ppu) copyX() {
	// http://wiki.nesdev.com/w/index.php/PPU_scrolling#At_dot_257_of_each_scanline
	// v: ....F.. ...EDCBA = t: ....F.. ...EDCBA
	p.v = (p.v &^ 0b100_00011111) | (p.t & 0b100_00011111)
}

func (p *ppu) copyY() {
	// http://wiki.nesdev.com/w/index.php/PPU_scrolling#During_dots_280_to_304_of_the_pre-render_scanline_.28end_of_vblank.29
	// v: IHGF.ED CBA..... = t: IHGF.ED CBA.....
	p.v = (p.v &^ 0b1111011_11100000) | (p.t & 0b1111011_11100000)
}

// https://www.nesdev.org/wiki/PPU_pattern_tables#Addressing
func tileAddr(v uint16) uint16 { return 0x2000 | (v & 0x0FFF) }
func attrAddr(v uint16) uint16 { return 0x23C0 | (v & 0x0C00) | ((v >> 4) & 0x38) | ((v >> 2) & 0x07) }

type FrameRenderer interface {
	UpdateFrame(*[WIDTH * HEIGHT]uint8)
}

type Sprite struct {
	enabled bool
	index   uint8

	x    uint8 // X position of left
	y    uint8 // Y position of top
	tile uint8 // tile index number
	attr uint8 // attribute

	low  uint8
	high uint8
}

func (s *Sprite) clear() {
	s.enabled = false
	s.index = 0xFF
	s.x = 0xFF
	s.y = 0xFF
	s.tile = 0xFF
	s.attr = 0xFF
	s.low = 0
	s.high = 0
}
