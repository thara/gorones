package mapper

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
)

// ROM wraps byte array of iNES format binary.
type ROM struct {
	header header
	raw    []byte
}

// Kind of Nametable Mirroring
// https://wiki.nesdev.org/w/index.php?title=Mirroring#Nametable_Mirroring
type Mirroring int

const (
	_ Mirroring = iota
	Mirroring_Horizontal
	Mirroring_Vertical
)

func (m Mirroring) String() string {
	switch m {
	case Mirroring_Horizontal:
		return "H"
	case Mirroring_Vertical:
		return "V"
	}
	return "Unknown"
}

type header struct {
	mapperNO   uint8
	prgROMSize uint
	chrROMSize uint
	mirroring  Mirroring
}

var (
	magicNumber = []byte{0x4E, 0x45, 0x53, 0x1A}
	padding     = []byte{0, 0, 0, 0, 0}
)

// ParseROM load NES binary program in iNES file format
func ParseROM(r io.Reader) (*ROM, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, errors.Wrap(err, "failed to parse for magic number")
	}
	if !bytes.Equal(buf, magicNumber) {
		return nil, errors.New("invalid magic number")
	}

	if n, err := io.ReadAtLeast(r, buf, 4); err != nil {
		return nil, errors.Wrap(err, "failed to parse for prg/chr rom size and mapper info`")
	} else if n != 4 {
		return nil, errors.Errorf("invalid prg/chr rom size and mapper info reading: n=%d", n)
	}

	prgROMSize := buf[0]
	chrROMSize := buf[1]
	flag6 := buf[2]
	flag7 := buf[3]

	var mirroring Mirroring
	if flag6&1 == 0 {
		mirroring = Mirroring_Horizontal
	} else {
		mirroring = Mirroring_Vertical
	}

	mapperNo := (flag6 & 0b1111000) | (flag7&0b11110000)<<4

	// skip flag 8, 9, 10
	buf = make([]byte, 3)
	if n, err := io.ReadAtLeast(r, buf, 3); err != nil {
		return nil, errors.Wrap(err, "failed to parse for flag 8..10")
	} else if n != 3 {
		return nil, errors.Errorf("invalid skip flags reading: n=%d", n)
	}

	// validate unused padding
	buf = make([]byte, 5)
	if n, err := io.ReadAtLeast(r, buf, 5); err != nil {
		return nil, errors.Wrap(err, "failed to parse for padding")
	} else if n != 5 {
		return nil, errors.Errorf("invalid padding reading: n=%d", n)
	} else if !bytes.Equal(buf, padding) {
		return nil, errors.New("invalid padding")
	}

	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read raw data after header")
	}
	return &ROM{
		header: header{
			mapperNO:   mapperNo,
			prgROMSize: uint(prgROMSize),
			chrROMSize: uint(chrROMSize),
			mirroring:  mirroring,
		},
		raw: raw,
	}, nil
}
