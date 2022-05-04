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

func (s *portTestSuite) Test_OAMADDR() {
	var addr uint16 = 0x2003

	s.ppu.Port.WriteRegister(addr, 255)
	s.EqualValues(255, s.ppu.oamAddr)
}

func (s *portTestSuite) Test_OAMDATA_read() {
	var addr uint16 = 0x2004

	s.ppu.spr.oam[0x09] = 0xA3
	s.ppu.Port.WriteRegister(0x2003, 0x09)

	s.EqualValues(0xA3, s.ppu.Port.ReadRegister(addr))
}

func (s *portTestSuite) Test_OAMDATA_write() {
	var addr uint16 = 0x2004

	s.ppu.Port.WriteRegister(0x2003, 0xAB)
	s.ppu.Port.WriteRegister(addr, 0x32)

	s.EqualValues(0x32, s.ppu.spr.oam[0xAB])
}

func (s *portTestSuite) Test_PPUSCROLL() {
	var addr uint16 = 0x2005

	s.ppu.Port.WriteRegister(addr, 0x1F)

	s.True(s.ppu.w)
	s.EqualValues(3, coarseX(s.ppu.t))
	s.EqualValues(0b111, s.ppu.x)

	s.ppu.Port.WriteRegister(addr, 0x0E)
	s.False(s.ppu.w)
	s.EqualValues(1, coarseY(s.ppu.t))
}

func (s *portTestSuite) Test_PPUADDR() {
	var addr uint16 = 0x2006

	s.ppu.Port.WriteRegister(addr, 0x3F)
	s.True(s.ppu.w)

	s.ppu.Port.WriteRegister(addr, 0x91)
	s.False(s.ppu.w)

	s.EqualValues(0x3F91, s.ppu.v)
	s.EqualValues(0x3F91, s.ppu.t)
}

func (s *portTestSuite) Test_PPUDATA() {
	var addr uint16 = 0x2007

	s.ppu.Port.WriteRegister(0x2006, 0x2F)
	s.ppu.Port.WriteRegister(0x2006, 0x11)

	s.ppu.Port.WriteRegister(addr, 0x83)

	s.EqualValues(0x83, s.ppu.read(0x2F11))
}

func Test_status(t *testing.T) {
	suite.Run(t, new(portTestSuite))
}
