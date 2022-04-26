package ppu

type PPU struct {
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
	v, t vramAddr
	// fine x scroll
	x uint8
	// first or second write toggle
	w bool
}

func (p *PPU) setController(v uint8) {
	p.ctrl.nt = v & 0b00000011
	p.ctrl.vramIncr = v&0b00000100 == 0b00000100
	p.ctrl.sprTable = v&0b00001000 == 0b00001000
	p.ctrl.bgTable = v&0b00010000 == 0b00010000
	p.ctrl.spr8x16 = v&0b00100000 == 0b00100000
	p.ctrl.slave = v&0b01000000 == 0b01000000
	p.ctrl.nmi = v&0b10000000 == 0b10000000
}

func (p *PPU) setMask(v uint8) {
	p.mask.gray = uint8(v)&0b00000001 == 0b00000001
	p.mask.bgLeft = uint8(v)&0b00000010 == 0b00000010
	p.mask.sprLeft = uint8(v)&0b00000100 == 0b00000100
	p.mask.bg = uint8(v)&0b00001000 == 0b00001000
	p.mask.spr = uint8(v)&0b00010000 == 0b00010000
}

func (p *PPU) readStatus() uint8 {
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

// type controller uint8

// func (c controller) nt() uint16     { return 0x2000 + 0x400*(uint16(c)&0b00000011) } // name table base index
// func (c controller) vramIncr() bool { return uint8(c)&0b00000100 == 0b00000100 }     // vram address increment
// func (c controller) sprTable() bool { return uint8(c)&0b00001000 == 0b00001000 }     // sprite pattern table
// func (c controller) bgTable() bool  { return uint8(c)&0b00010000 == 0b00010000 }     // background pattern table
// func (c controller) spr8x16() bool  { return uint8(c)&0b00100000 == 0b00100000 }     // sprite size
// func (c controller) slave() bool    { return uint8(c)&0b01000000 == 0b01000000 }     // PPU master/slave
// func (c controller) nmi() bool      { return uint8(c)&0b10000000 == 0b10000000 }     // enabled NMI

// type mask uint8

// func (m mask) gray() bool    { return uint8(m)&0b00000001 == 0b00000001 } // grayscale
// func (m mask) bgLeft() bool  { return uint8(m)&0b00000010 == 0b00000010 } // show background in leftmost 8 pixels
// func (m mask) sprLeft() bool { return uint8(m)&0b00000100 == 0b00000100 } // show sprite in leftmost 8 pixels
// func (m mask) bg() bool      { return uint8(m)&0b00001000 == 0b00001000 } // show background
// func (m mask) spr() bool     { return uint8(m)&0b00010000 == 0b00010000 } // show sprites

// type status uint8

// func (s status) sprOverflow() bool { return uint8(s)&0b00100000 == 0b00100000 } // sprite overflow
// func (s status) spr0Hit() bool     { return uint8(s)&0b01000000 == 0b01000000 } // sprite 0 hit
// func (s status) vblank() bool      { return uint8(s)&0b10000000 == 0b10000000 } // is vblank

// vramAddr is VRAM address (15bits)
type vramAddr = uint16

// https://www.nesdev.org/wiki/PPU_scrolling#PPU_internal_registers
type scroll vramAddr

func (v scroll) coarseX() uint8 { return uint8(uint16(v) & 0b000000000011111) }
func (v scroll) coarseY() uint8 { return uint8(uint16(v) & 0b000001111100000 >> 5) }
func (v scroll) nt() uint8      { return uint8(uint16(v) & 0b000110000000000 >> 10) }
func (v scroll) fineY() uint8   { return uint8(uint16(v) & 0b111000000000000 >> 12) }

// https://www.nesdev.org/wiki/PPU_pattern_tables#Addressing
type inPatternTable vramAddr

func (v inPatternTable) addr() uint16 { return uint16(v) & 0b00111111111111 }

type scan struct {
	line  uint8 // 0 ..= 261
	cycle uint8 // 0 ..= 340
}
