package main

import (
	"image"
	"image/color"

	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/util"
)

func loadCHRPatterns(m mapper.Mapper) []pattern {
	chr := m.CHR()

	p := make([]pattern, len(chr)/16)

	for i, c := range chr {
		if i%16 < 8 {
			p[i/16].low[i%16] |= c
		} else {
			p[i/16].high[i%16-8] |= c
		}
	}
	return p
}

type pattern struct {
	low  [8]byte
	high [8]byte
}

func (p *pattern) slice() [64]struct{ x, y, p uint } {
	var s [64]struct{ x, y, p uint }
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			s[y*8+x].x = uint(x)
			s[y*8+x].y = uint(y)
			col := 7 - x
			h := uint(p.high[y])
			l := uint(p.low[y])
			s[y*8+x].p = util.NthBit(h, col)<<1 | util.NthBit(l, col)
		}
	}
	return s
}

func (p *pattern) write(i int, img *image.RGBA) {
	offset := i * 8
	for _, s := range p.slice() {
		var c color.RGBA
		c.R = uint8(palletes[s.p] >> 4 & 0xFF)
		c.G = uint8(palletes[s.p] >> 2 & 0xFF)
		c.B = uint8(palletes[s.p] & 0xFF)
		c.A = 0xFF
		img.Set(offset+int(s.x), int(s.y), c)
	}
}

var palletes []uint32 = []uint32{
	0x000000, 0x858585, 0xAAAAAA, 0xFFFFFF,
}
