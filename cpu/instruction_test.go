package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type busMock []uint8

func newBusMock() busMock {
	return make(busMock, 0x10000)
}

func (m busMock) ReadCPU(addr uint16) uint8 {
	return m[addr]
}

func (m busMock) WriteCPU(addr uint16, value uint8) {
	m[addr] = value
}

type tickFn func()

func (f tickFn) Tick() {
	f()
}

var tickMock = tickFn(func() {})

type getOperandTestSuite struct {
	suite.Suite
	bus []uint8
	emu *Emu
}

func (s *getOperandTestSuite) SetupTest() {
	s.bus = newBusMock()
	s.emu = NewEmu(tickMock, busMock(s.bus))
}

func (s *getOperandTestSuite) Test_implicit() {
	v := s.emu.getOperand(implicit)
	s.EqualValues(0, v)
	s.EqualValues(0, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_accumulator() {
	s.emu.cpu.A = 0xFB

	v := s.emu.getOperand(accumulator)
	s.EqualValues(0xFB, v)
	s.EqualValues(0, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_immediate() {
	s.emu.cpu.PC = 0x8234

	v := s.emu.getOperand(immediate)
	s.EqualValues(0x8234, v)
	s.EqualValues(0, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_zeroPage() {
	s.emu.cpu.PC = 0x0414
	s.bus[0x0414] = 0x91

	v := s.emu.getOperand(zeroPage)
	s.EqualValues(0x91, v)
	s.EqualValues(1, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_zeroPageX() {
	s.emu.cpu.PC = 0x0100
	s.emu.cpu.X = 0x93
	s.bus[0x0100] = 0x80

	v := s.emu.getOperand(zeroPageX)
	s.EqualValues(0x13, v)
	s.EqualValues(2, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_zeroPageY() {
	s.emu.cpu.PC = 0x0423
	s.emu.cpu.Y = 0xF1
	s.bus[0x0423] = 0x36

	v := s.emu.getOperand(zeroPageY)
	s.EqualValues(0x27, v)
	s.EqualValues(2, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_absolute() {
	s.emu.cpu.PC = 0x0423
	s.bus[0x0423] = 0x36
	s.bus[0x0424] = 0xF0

	v := s.emu.getOperand(absolute)
	s.EqualValues(0xF036, v)
	s.EqualValues(2, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_absoluteX() {
	s.bus[0x0423] = 0x36
	s.bus[0x0424] = 0xF0

	tests := []struct {
		name            string
		mode            addressingMode
		x               uint8
		expectedOperand uint16
		expectedCycles  int
	}{
		{"no oops", absoluteX, 0x31, 0xF067, 3},
		{"oops/not page crossed", absoluteXWithPenalty, 0x31, 0xF067, 2},
		{"oops/page crossed", absoluteXWithPenalty, 0xF0, 0xF126, 3},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x0423
			s.emu.cpu.X = tt.x

			v := s.emu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, s.emu.cpu.Cycles)
		})
	}
}

func (s *getOperandTestSuite) Test_absoluteY() {
	s.bus[0x0423] = 0x36
	s.bus[0x0424] = 0xF0

	tests := []struct {
		name            string
		mode            addressingMode
		y               uint8
		expectedOperand uint16
		expectedCycles  int
	}{
		{"no oops", absoluteY, 0x31, 0xF067, 3},
		{"oops/not page crossed", absoluteYWithPenalty, 0x31, 0xF067, 2},
		{"oops/page crossed", absoluteYWithPenalty, 0xF0, 0xF126, 3},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x0423
			s.emu.cpu.Y = tt.y

			v := s.emu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, s.emu.cpu.Cycles)
		})
	}
}

func (s *getOperandTestSuite) Test_relative() {
	s.emu.cpu.PC = 0x0414
	s.bus[0x0414] = 0x91

	v := s.emu.getOperand(relative)
	s.EqualValues(0x91, v)
	s.EqualValues(1, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_indirect() {
	s.emu.cpu.PC = 0x020F
	s.bus[0x020F] = 0x10
	s.bus[0x0210] = 0x03
	s.bus[0x0310] = 0x9F

	v := s.emu.getOperand(indirect)
	s.EqualValues(0x9F, v)
	s.EqualValues(4, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_indexedIndirect() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.X = 0x95
	s.bus[0x020F] = 0xF0
	s.bus[0x0085] = 0x12
	s.bus[0x0086] = 0x90

	v := s.emu.getOperand(indexedIndirect)
	s.EqualValues(0x9012, v)
	s.EqualValues(4, s.emu.cpu.Cycles)
}

func (s *getOperandTestSuite) Test_indirectIndexed() {
	s.bus[0x020F] = 0xF0
	s.bus[0x00F0] = 0x12
	s.bus[0x00F1] = 0x90

	tests := []struct {
		name            string
		mode            addressingMode
		y               uint8
		expectedOperand uint16
		expectedCycles  int
	}{
		{"no oops", indirectIndexed, 0xF3, 0x9105, 4},
		{"not page crossed", indirectIndexedWithPenalty, 0x83, 0x9095, 3},
		{"page crossed", indirectIndexedWithPenalty, 0xF3, 0x9105, 4},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x020F
			s.emu.cpu.Y = tt.y

			v := s.emu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, s.emu.cpu.Cycles)
		})
	}
}

func Test_getOperand(t *testing.T) {
	suite.Run(t, new(getOperandTestSuite))
}

type executeTestSuite struct {
	suite.Suite
	bus []uint8
	emu *Emu
}

func (s *executeTestSuite) SetupTest() {
	s.bus = newBusMock()
	s.emu = NewEmu(tickMock, busMock(s.bus))
}

func (s *executeTestSuite) Test_LDA() {
	s.emu.cpu.PC = 0x020F
	s.bus[0x020F] = 0xA9
	s.bus[0x0210] = 0x31

	s.emu.Step()

	s.EqualValues(0x31, s.emu.cpu.A)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_STA() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.A = 0x91
	s.bus[0x020F] = 0x8D
	s.bus[0x0210] = 0x19
	s.bus[0x0211] = 0x04

	s.emu.Step()

	s.EqualValues(0x91, s.bus[0x0419])
	s.EqualValues(4, s.emu.cpu.Cycles)
	s.EqualValues(0, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_TAX() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.A = 0x83
	s.bus[0x020F] = 0xAA

	s.emu.Step()

	s.EqualValues(0x83, s.emu.cpu.X)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0x80, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_TYA() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.Y = 0xF0
	s.bus[0x020F] = 0x98

	s.emu.Step()

	s.EqualValues(0xF0, s.emu.cpu.A)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0x80, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_TSX() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.S = 0xF3
	s.bus[0x020F] = 0xBA

	s.emu.Step()

	s.EqualValues(0xF3, s.emu.cpu.X)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0x80, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_PHA() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.S = 0xFD
	s.emu.cpu.A = 0x72
	s.bus[0x020F] = 0x48

	s.emu.Step()

	s.EqualValues(0xFC, s.emu.cpu.S)
	s.EqualValues(0x72, s.bus[0x01FD])
	s.EqualValues(3, s.emu.cpu.Cycles)
	s.EqualValues(0, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_PHP() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.S = 0xFD
	s.emu.cpu.A = 0x72
	s.emu.cpu.P[status_N] = true
	s.emu.cpu.P[status_D] = true
	s.emu.cpu.P[status_C] = true
	s.bus[0x020F] = 0x08

	s.emu.Step()

	s.EqualValues(0xFC, s.emu.cpu.S)
	s.EqualValues(s.emu.cpu.P.u8()|instructionB, s.bus[0x01FD])
	s.EqualValues(3, s.emu.cpu.Cycles)
}

func (s *executeTestSuite) Test_PLP() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.S = 0xBF
	s.bus[0x020F] = 0x28
	s.bus[0x01C0] = 0x7A

	s.emu.Step()

	s.EqualValues(0xC0, s.emu.cpu.S)
	s.EqualValues(4, s.emu.cpu.Cycles)
	s.EqualValues(0b1001010, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_EOR() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.A = 0x21
	s.bus[0x020F] = 0x49
	s.bus[0x0210] = 0x38

	s.emu.Step()

	s.EqualValues(0x19, s.emu.cpu.A)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_BIT() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.A = 0x48
	s.bus[0x020F] = 0x2C
	s.bus[0x0210] = 0xB0
	s.bus[0x0211] = 0x03
	s.bus[0x03B0] = 0b11000000

	s.emu.Step()

	s.EqualValues(4, s.emu.cpu.Cycles)
	s.EqualValues(0b11000000, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_ADC() {
	tests := []struct {
		a         uint8
		m         uint8
		expectedA uint8
		expectedP uint8
	}{
		{0x50, 0x10, 0x60, 0b00000000},
		{0x50, 0x50, 0xA0, 0b11000000},
		{0x50, 0x90, 0xE0, 0b10000000},
		{0x50, 0xD0, 0x20, 0b00000001},
		{0xD0, 0x10, 0xE0, 0b10000000},
		{0xD0, 0x50, 0x20, 0b00000001},
		{0xD0, 0x90, 0x60, 0b01000001},
		{0xD0, 0xD0, 0xA0, 0b10000001},
	}
	for i, tt := range tests {
		s.Run(fmt.Sprintf("pattern:%d", i), func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x020F
			s.emu.cpu.A = tt.a

			s.bus[0x020F] = 0x6D
			s.bus[0x0210] = 0xD3
			s.bus[0x0211] = 0x04
			s.bus[0x04D3] = tt.m

			s.emu.Step()

			s.EqualValues(tt.expectedA, s.emu.cpu.A)
			s.EqualValues(tt.expectedP, s.emu.cpu.P.u8())
		})
	}
}

func (s *executeTestSuite) Test_CPY() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.Y = 0x37
	s.bus[0x020F] = 0xCC
	s.bus[0x0210] = 0x36

	s.emu.Step()

	s.EqualValues(4, s.emu.cpu.Cycles)
	s.EqualValues(0b00000001, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_INC() {
	s.emu.cpu.PC = 0x020F
	s.bus[0x020F] = 0xEE
	s.bus[0x0210] = 0xD3
	s.bus[0x0211] = 0x04
	s.bus[0x04D3] = 0x7F

	s.emu.Step()

	s.EqualValues(6, s.emu.cpu.Cycles)
	s.EqualValues(0x80, s.bus[0x04D3])
	s.EqualValues(0b10000000, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_DEC() {
	s.emu.cpu.PC = 0x020F
	s.bus[0x020F] = 0xCE
	s.bus[0x0210] = 0xD3
	s.bus[0x0211] = 0x04
	s.bus[0x04D3] = 0xC0

	s.emu.Step()

	s.EqualValues(6, s.emu.cpu.Cycles)
	s.EqualValues(0xBF, s.bus[0x04D3])
	s.EqualValues(0b10000000, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_ASL() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.A = 0b10001010

	s.bus[0x020F] = 0x0A

	s.emu.Step()

	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0b00010100, s.emu.cpu.A)
	s.EqualValues(0b00000001, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_ROL() {
	s.bus[0x020F] = 0x2A

	tests := []struct {
		name      string
		p         uint8
		expectedA uint8
	}{
		{"no carry", 0b00000001, 0b00010101},
		{"carry", 0b10000000, 0b00010100},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x020F
			s.emu.cpu.A = 0b10001010
			s.emu.cpu.P.set(tt.p)

			s.emu.Step()

			s.EqualValues(2, s.emu.cpu.Cycles)
			s.EqualValues(tt.expectedA, s.emu.cpu.A)
			s.EqualValues(0b00000001, s.emu.cpu.P.u8())
		})
	}
}

func (s *executeTestSuite) Test_JSR() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.S = 0xBF

	s.bus[0x020F] = 0x20
	s.bus[0x0210] = 0x31
	s.bus[0x0211] = 0x40

	s.emu.Step()

	s.EqualValues(0xBD, s.emu.cpu.S)
	s.EqualValues(0x4031, s.emu.cpu.PC)
	s.EqualValues(6, s.emu.cpu.Cycles)
	s.EqualValues(0x11, s.bus[0x01BE])
	s.EqualValues(0x02, s.bus[0x01BF])
}

func (s *executeTestSuite) Test_RTS() {
	s.emu.cpu.PC = 0x0031
	s.emu.cpu.S = 0xBD

	s.bus[0x0031] = 0x60
	s.bus[0x01BE] = 0x11
	s.bus[0x01BF] = 0x02

	s.emu.Step()

	s.EqualValues(0xBF, s.emu.cpu.S)
	s.EqualValues(0x0212, s.emu.cpu.PC)
	s.EqualValues(6, s.emu.cpu.Cycles)
}

func (s *executeTestSuite) Test_BCC() {
	tests := []struct {
		name           string
		operand        uint8
		p              uint8
		expectedPC     uint16
		expectedCycles uint
	}{
		{"branch failed", 0x03, 0b10000001, 0x33, 2},
		{"branch succeed", 0x03, 0b11000000, 0x36, 3},
		{"branch succeed & new page", 0xD0, 0b11000000, 0x03, 3},
	}
	for _, tt := range tests {
		s.Run(tt.name, func() {
			s.emu.cpu.reset()
			s.emu.cpu.PC = 0x0031
			s.emu.cpu.P.set(tt.p)

			s.bus[0x0031] = 0x90
			s.bus[0x0032] = tt.operand

			s.emu.Step()

			s.EqualValues(tt.expectedPC, s.emu.cpu.PC)
			s.EqualValues(tt.expectedCycles, s.emu.cpu.Cycles)
		})
	}
}

func (s *executeTestSuite) Test_CLD() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.P.set(0b011001001)

	s.bus[0x020F] = 0xD8

	s.emu.Step()

	s.EqualValues(0x0210, s.emu.cpu.PC)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0b011000001, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_SEI() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.P.set(0b011001001)

	s.bus[0x020F] = 0x78

	s.emu.Step()

	s.EqualValues(0x0210, s.emu.cpu.PC)
	s.EqualValues(2, s.emu.cpu.Cycles)
	s.EqualValues(0b011001101, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_BRK() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.P.set(0b01100001)
	s.emu.cpu.S = 0xBF

	s.bus[0x020F] = 0x00
	s.bus[0xFFFE] = 0x23
	s.bus[0xFFFF] = 0x40

	s.emu.Step()

	s.EqualValues(0x4023, s.emu.cpu.PC)
	s.EqualValues(7, s.emu.cpu.Cycles)
	s.EqualValues(0xBC, s.emu.cpu.S)
	s.EqualValues(0b01000001, s.emu.cpu.P.u8())
}

func (s *executeTestSuite) Test_RTI() {
	s.emu.cpu.PC = 0x020F
	s.emu.cpu.P.set(0b01100101)
	s.emu.cpu.S = 0xBC

	s.bus[0x020F] = 0x40
	s.bus[0x01BD] = 0b10000010
	s.bus[0x01BE] = 0x11
	s.bus[0x01BF] = 0x02

	s.emu.Step()

	s.EqualValues(0x0211, s.emu.cpu.PC)
	s.EqualValues(6, s.emu.cpu.Cycles)
	s.EqualValues(0xBF, s.emu.cpu.S)
	s.EqualValues(0b10000010, s.emu.cpu.P.u8())
}

func Test_execute(t *testing.T) {
	suite.Run(t, new(executeTestSuite))
}
