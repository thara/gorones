package main

import (
	"bytes"
	"encoding/binary"
	"log"

	"github.com/thara/gorones/ppu"
)

type renderer struct {
	px []byte
}

func newRenderer() *renderer {
	return &renderer{
		px: make([]byte, 4*ppu.WIDTH*ppu.HEIGHT),
	}
}

func (r *renderer) UpdateFrame(buf *[ppu.WIDTH * ppu.HEIGHT]uint8) {
	for i, v := range buf {
		r.px[i] = palettes[v]
		r.px[i+1] = palettes[v+1]
		r.px[i+2] = palettes[v+2]
		r.px[i+3] = palettes[v+3]
	}
}

func (r *renderer) pixels() []byte {
	return r.px
}

var palettes []byte = func() []byte {
	var bf bytes.Buffer
	if err := binary.Write(&bf, binary.BigEndian, rgba); err != nil {
		log.Fatalf("fail to write video buffer: %v", err)
	}
	return bf.Bytes()
}()

var rgba []uint32 = []uint32{
	0x7C7C7CFF, 0x0000FCFF, 0x0000BCFF, 0x4428BCFF, 0x940084FF, 0xA80020FF, 0xA81000FF, 0x881400FF,
	0x503000FF, 0x007800FF, 0x006800FF, 0x005800FF, 0x004058FF, 0x000000FF, 0x000000FF, 0x000000FF,
	0xBCBCBCFF, 0x0078F8FF, 0x0058F8FF, 0x6844FCFF, 0xD800CCFF, 0xE40058FF, 0xF83800FF, 0xE45C10FF,
	0xAC7C00FF, 0x00B800FF, 0x00A800FF, 0x00A844FF, 0x008888FF, 0x000000FF, 0x000000FF, 0x000000FF,
	0xF8F8F8FF, 0x3CBCFCFF, 0x6888FCFF, 0x9878F8FF, 0xF878F8FF, 0xF85898FF, 0xF87858FF, 0xFCA044FF,
	0xF8B800FF, 0xB8F818FF, 0x58D854FF, 0x58F898FF, 0x00E8D8FF, 0x787878FF, 0x000000FF, 0x000000FF,
	0xFCFCFCFF, 0xA4E4FCFF, 0xB8B8F8FF, 0xD8B8F8FF, 0xF8B8F8FF, 0xF8A4C0FF, 0xF0D0B0FF, 0xFCE0A8FF,
	0xF8D878FF, 0xD8F878FF, 0xB8F8B8FF, 0xB8F8D8FF, 0x00FCFCFF, 0xF8D8F8FF, 0x000000FF, 0x000000}
