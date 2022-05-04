package ppu

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thara/gorones/mapper"
)

type nopFrameRenderer struct{}

func (nopFrameRenderer) UpdateFrame(*[WIDTH * HEIGHT]uint8) {}

type portTestSuite struct {
	suite.Suite

	ppu *PPU
}

func (s *portTestSuite) SetupTest() {
	s.ppu = New(new(mapper.MapperMock), new(nopFrameRenderer))
}

func (s *portTestSuite) Test_PPUCTRL() {
	var addr uint16 = 0x2000

	s.ppu.Port.WriteRegister(addr, 0b01010000)
	s.True(s.ppu.ctrl.slave)
	s.True(s.ppu.ctrl.bgTable)

	s.ppu.Port.WriteRegister(addr, 0b00000111)
	s.True(s.ppu.ctrl.vramIncr)
	s.EqualValues(0b11, s.ppu.ctrl.nt)
}

func (s *portTestSuite) Test_PPUMASK() {
	var addr uint16 = 0x2001

	s.ppu.Port.WriteRegister(addr, 0b01010000)
	s.True(s.ppu.mask.green)
	s.True(s.ppu.mask.spr)

	s.ppu.Port.WriteRegister(addr, 0b00000111)
	s.True(s.ppu.mask.sprLeft)
	s.True(s.ppu.mask.bgLeft)
	s.True(s.ppu.mask.gray)
}

func (s *portTestSuite) Test_PPUSTATUS() {
	var addr uint16 = 0x2002

	s.ppu.status.vblank = true
	s.ppu.status.spr0Hit = true
	s.ppu.w = true

	s.EqualValues(0b11000000, s.ppu.Port.ReadRegister(addr))
	s.False(s.ppu.w)
	s.EqualValues(0b01000000, s.ppu.Port.ReadRegister(addr))
}

func Test_status(t *testing.T) {
	suite.Run(t, new(portTestSuite))
}
