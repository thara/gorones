package mapper

import (
	"fmt"

	"github.com/pkg/errors"
)

// https://www.nesdev.org/wiki/Mapper

// Mapper emulates circuits, hardware, and the configuration and capabilities of cartridges.
type Mapper interface {
	// Read reads a byte from the mapper.
	Read(addr uint16) uint8

	// Write writes a byte into the mapper.
	Write(addr uint16, value uint8)

	Mirroring() Mirroring

	PRG() []byte
	CHR() []byte
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

	mirroring Mirroring
	mirrored  bool
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
		prg:       prg,
		chr:       chr,
		mirroring: rom.header.mirroring,
		mirrored:  prgSize == 0x4000,
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

func (m *mapper0) Mirroring() Mirroring {
	return m.mirroring
}

func (m *mapper0) PRG() []byte { return append([]byte(nil), m.prg...) }
func (m *mapper0) CHR() []byte { return append([]byte(nil), m.chr...) }

func (m mapper0) String() string {
	return fmt.Sprintf(`mapper 0:
	PRG: 0x%x byte
	CHR: 0x%x byte
	mirroring: %s
	mirrored: %t
`, len(m.prg), len(m.chr), m.mirroring, m.mirrored)
}
