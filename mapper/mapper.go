package mapper

import (
	"github.com/pkg/errors"
)

// https://www.nesdev.org/wiki/Mapper

// Mapper emulates circuits, hardware, and the configuration and capabilities of cartridges.
type Mapper interface {
	// Read reads a byte from the mapper.
	Read(addr uint16) uint8

	// Write writes a byte into the mapper.
	Write(addr uint16, value uint8)
}

// Mapper creates a mapper object from this rom's data
func (r *ROM) Mapper() (Mapper, error) {
	switch r.header.mapperNO {
	case 0:
		return newMapper0(r), nil
	}
	return nil, errors.Errorf("unsupported mapper no: %d", r.header.mapperNO)
}

type mapper0 struct {
	prg []byte
	chr []byte

	mirrored bool
}

func newMapper0(rom *ROM) Mapper {
	prgSize := rom.header.prgROMSize * 0x4000
	prg := rom.raw[:prgSize]

	var chr []byte
	if rom.header.chrROMSize == 0 {
		chr = make([]byte, 0x2000)
	} else {
		chrSize := rom.header.chrROMSize * 0x2000
		chr = rom.raw[prgSize : prgSize+chrSize]
	}
	return &mapper0{
		prg:      prg,
		chr:      chr,
		mirrored: prgSize == 0x4000,
	}
}

func (m *mapper0) Read(addr uint16) uint8 {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		return m.chr[addr]
	case 0x8000 <= addr && addr <= 0xFFFF:
		if m.mirrored {
			addr %= 0x4000
		} else {
			addr -= 0x8000
		}
		return m.prg[addr]
	}
	return 0
}

func (m *mapper0) Write(addr uint16, value uint8) {
	switch {
	case 0x0000 <= addr && addr <= 0x1FFF:
		m.chr[addr] = value
	}
}
