package ppu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/mapper"
)

func Test_coarseX(t *testing.T) {
	assert.EqualValues(t, 0b11101, coarseX(0b11011011_00111101))
}

func Test_coarseY(t *testing.T) {
	assert.EqualValues(t, 0b11001, coarseY(0b11011011_00111101))
}

func Test_fineY(t *testing.T) {
	assert.EqualValues(t, 0b101, fineY(0b11011011_00111101))
}

func Test_tileAddr(t *testing.T) {
	assert.EqualValues(t, 0b10101100111101, tileAddr(0b11011011_00111101))
}

func Test_attrAddr(t *testing.T) {
	assert.EqualValues(t, 0b10101111110111, attrAddr(0b11011011_00111101))
}

type mapperStub struct {
	s []byte
}

func (s *mapperStub) Read(addr uint16) uint8         { return s.s[addr] }
func (s *mapperStub) Write(addr uint16, value uint8) { s.s[addr] = value }
func (*mapperStub) Mirroring() mapper.Mirroring      { return mapper.Mirroring_Horizontal }
func (*mapperStub) PRG() []byte                      { return nil }
func (*mapperStub) CHR() []byte                      { return nil }

func Test_bg(t *testing.T) {
	var intr *cpu.Interrupt
	m := mapperStub{make([]byte, 65534)}
	ppu := New(&m, new(nopFrameRenderer))

	ppu.v = 0b101_10_11001_11101 // y: 5, nt: 2, y: 11001, x: 11101
	ppu.write(0x0035, 0x11)
	ppu.write(0x003D, 0x81)
	ppu.write(0x2B3D, 0x03)
	ppu.write(0x2BF7, 0x41)

	ppu.scan.dot = 1

	ppu.Step(intr)
	assert.EqualValues(t, 0x2B3D, ppu.bg.addr, "fetch name table: step 1")

	ppu.Step(intr)
	assert.EqualValues(t, 0x03, ppu.bg.nt, "fetch name table: step 2")

	ppu.Step(intr)
	assert.EqualValues(t, 0x2BF7, ppu.bg.addr, "fetch attribute table: step 1")

	ppu.Step(intr)
	assert.EqualValues(t, 0x41, ppu.bg.at, "fetch attribute table: step 2")

	ppu.Step(intr)
	assert.EqualValues(t, 0x0035, ppu.bg.addr, "fetch tile bitmap low byte: step 1")

	ppu.Step(intr)
	assert.EqualValues(t, 0x11, ppu.bg.low, "fetch tile bitmap low byte: step 2")

	ppu.Step(intr)
	assert.EqualValues(t, 0x003D, ppu.bg.addr, "Fetch tile bitmap high byte : step 1")

	ppu.Step(intr)
	assert.EqualValues(t, 0x81, ppu.bg.high, "Fetch tile bitmap high byte : step 2")
}

func Test_incrCoarseX(t *testing.T) {
	t.Run("increment coarse X", func(t *testing.T) {
		ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

		ppu.v = 0b000_10_11001_11101 // y: 0, nt: 2, y: 11001, x: 11101

		assert.EqualValues(t, 29, coarseX(ppu.v))

		ppu.incrCoarseX()
		assert.EqualValues(t, 30, coarseX(ppu.v))
	})

	t.Run("switch horizontal nametable", func(t *testing.T) {
		ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

		ppu.v = 0b000_11_11001_11111 // y: 0, nt: 3, y: 11001, x: 11111

		assert.EqualValues(t, 31, coarseX(ppu.v))
		assert.EqualValues(t, 3, nameTableSelect(ppu.v))

		ppu.incrCoarseX()
		assert.EqualValues(t, 0, coarseX(ppu.v))
		assert.EqualValues(t, 2, nameTableSelect(ppu.v))
	})
}

func Test_incrY(t *testing.T) {
	t.Run("increment fine Y", func(t *testing.T) {
		ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

		ppu.v = 0b101_10_10101_11101 // y: 101, nt: 2, y: 10101, x: 11101

		ppu.incrY()
		assert.EqualValues(t, 0b01101010_10111101, ppu.v)
	})

	t.Run("if fine Y == 7", func(t *testing.T) {
		t.Run("switch vertical nametable", func(t *testing.T) {
			ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

			ppu.v = 0b111_10_11101_11101 // y: 7, nt: 2, y: 29, x: 11101

			ppu.incrY()

			assert.EqualValues(t, 0, fineY(ppu.v))
			assert.EqualValues(t, 0, nameTableSelect(ppu.v))
			assert.EqualValues(t, 0, coarseY(ppu.v))
			assert.EqualValues(t, 0b11101, coarseX(ppu.v))
		})

		t.Run("clear coarse Y", func(t *testing.T) {
			ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

			ppu.v = 0b111_10_11111_11101 // y: 7, nt: 2, y: 31, x: 11101

			ppu.incrY()
			assert.EqualValues(t, 0, fineY(ppu.v))
			assert.EqualValues(t, 2, nameTableSelect(ppu.v))
			assert.EqualValues(t, 0, coarseY(ppu.v))
			assert.EqualValues(t, 0b11101, coarseX(ppu.v))
		})

		t.Run("increment coarse Y", func(t *testing.T) {
			ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

			ppu.v = 0b111_10_01011_11101 // y: 7, nt: 2, y: 31, x: 11101

			ppu.incrY()
			assert.EqualValues(t, 0, fineY(ppu.v))
			assert.EqualValues(t, 2, nameTableSelect(ppu.v))
			assert.EqualValues(t, 12, coarseY(ppu.v))
			assert.EqualValues(t, 0b11101, coarseX(ppu.v))
		})
	})
}

func Test_copyX(t *testing.T) {
	ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

	ppu.v = 0b000_10_01011_11101 // y: 0, nt: 2, y: 11, x: 29
	ppu.t = 0b000_11_01010_00101 // y: 0, nt: 3, y: 10, x: 5

	assert.EqualValues(t, 3, nameTableSelect(ppu.t))
	assert.EqualValues(t, 10, coarseY(ppu.t))
	assert.EqualValues(t, 5, coarseX(ppu.t))

	ppu.copyX()

	assert.EqualValues(t, 3, nameTableSelect(ppu.v))
	assert.EqualValues(t, 11, coarseY(ppu.v))
	assert.EqualValues(t, 5, coarseX(ppu.v))
}

func Test_copyY(t *testing.T) {
	ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

	ppu.v = 0b000_10_01011_11101 // y: 0, nt: 2, y: 11, x: 29
	ppu.t = 0b000_01_01010_00101 // y: 0, nt: 1, y: 10, x: 5

	assert.EqualValues(t, 1, nameTableSelect(ppu.t))
	assert.EqualValues(t, 10, coarseY(ppu.t))
	assert.EqualValues(t, 5, coarseX(ppu.t))

	ppu.copyY()

	assert.EqualValues(t, 0, nameTableSelect(ppu.v))
	assert.EqualValues(t, 10, coarseY(ppu.v))
	assert.EqualValues(t, 29, coarseX(ppu.v))
}

func Test_bgShift(t *testing.T) {
	ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

	ppu.bg.shiftL = 0b10101001
	ppu.bg.shiftH = 0b00101101
	ppu.bg.attrShiftL = 0b11001010
	ppu.bg.attrShiftH = 0b01111100
	ppu.bg.attrLatchL = 1
	ppu.bg.attrLatchH = 0

	ppu.bgShift()

	assert.EqualValues(t, 0b101010010, ppu.bg.shiftL)
	assert.EqualValues(t, 0b01011010, ppu.bg.shiftH)
	assert.EqualValues(t, 0b10010101, ppu.bg.attrShiftL)
	assert.EqualValues(t, 0b11111000, ppu.bg.attrShiftH)
}

func Test_bgShiftReload(t *testing.T) {
	ppu := New(new(mapper.MapperMock), new(nopFrameRenderer))

	ppu.bg.shiftL = 0b10101001_00101101
	ppu.bg.shiftH = 0b11101101_00011001
	ppu.bg.attrLatchL = 1
	ppu.bg.attrLatchH = 0

	ppu.bg.low = 0b11111111
	ppu.bg.high = 0b00000000
	ppu.bg.at = 0b1010110

	ppu.bgShiftReload()

	assert.EqualValues(t, 0b10101001_11111111, ppu.bg.shiftL)
	assert.EqualValues(t, 0b11101101_00000000, ppu.bg.shiftH)
	assert.EqualValues(t, 0, ppu.bg.attrLatchL)
	assert.EqualValues(t, 1, ppu.bg.attrLatchH)
}
