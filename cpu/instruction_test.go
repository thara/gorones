package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type executeTestSuite struct {
	suite.Suite
	bus []uint8
	cpu CPU
}

func (s *executeTestSuite) SetupTest() {
	s.bus = newBusMock()
	s.cpu = New(tickMock, busMock(s.bus))
}

func (s *executeTestSuite) Test_LDA() {
	s.cpu.pc = 0x020F
	s.bus[0x020F] = 0xA9
	s.bus[0x0210] = 0x31

	s.cpu.Step()

	s.EqualValues(0x31, s.cpu.a)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_STA() {
	s.cpu.pc = 0x020F
	s.cpu.a = 0x91
	s.bus[0x020F] = 0x8D
	s.bus[0x0210] = 0x19
	s.bus[0x0211] = 0x04

	s.cpu.Step()

	s.EqualValues(0x91, s.bus[0x0419])
	s.EqualValues(4, s.cpu.cycles)
	s.EqualValues(0, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_TAX() {
	s.cpu.pc = 0x020F
	s.cpu.a = 0x83
	s.bus[0x020F] = 0xAA

	s.cpu.Step()

	s.EqualValues(0x83, s.cpu.x)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0x80, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_TYA() {
	s.cpu.pc = 0x020F
	s.cpu.y = 0xF0
	s.bus[0x020F] = 0x98

	s.cpu.Step()

	s.EqualValues(0xF0, s.cpu.a)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0x80, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_TSX() {
	s.cpu.pc = 0x020F
	s.cpu.s = 0xF3
	s.bus[0x020F] = 0xBA

	s.cpu.Step()

	s.EqualValues(0xF3, s.cpu.x)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0x80, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_PHA() {
	s.cpu.pc = 0x020F
	s.cpu.s = 0xFD
	s.cpu.a = 0x72
	s.bus[0x020F] = 0x48

	s.cpu.Step()

	s.EqualValues(0xFC, s.cpu.s)
	s.EqualValues(0x72, s.bus[0x01FD])
	s.EqualValues(3, s.cpu.cycles)
	s.EqualValues(0, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_PHP() {
	s.cpu.pc = 0x020F
	s.cpu.s = 0xFD
	s.cpu.a = 0x72
	s.cpu.p[status_N] = true
	s.cpu.p[status_D] = true
	s.cpu.p[status_C] = true
	s.bus[0x020F] = 0x08

	s.cpu.Step()

	s.EqualValues(0xFC, s.cpu.s)
	s.EqualValues(s.cpu.p.u8()|instructionB, s.bus[0x01FD])
	s.EqualValues(3, s.cpu.cycles)
}

func (s *executeTestSuite) Test_PLP() {
	s.cpu.pc = 0x020F
	s.cpu.s = 0xBF
	s.bus[0x020F] = 0x28
	s.bus[0x01C0] = 0x7A

	s.cpu.Step()

	s.EqualValues(0xC0, s.cpu.s)
	s.EqualValues(4, s.cpu.cycles)
	s.EqualValues(0b1001010, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_EOR() {
	s.cpu.pc = 0x020F
	s.cpu.a = 0x21
	s.bus[0x020F] = 0x49
	s.bus[0x0210] = 0x38

	s.cpu.Step()

	s.EqualValues(0x19, s.cpu.a)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_BIT() {
	s.cpu.pc = 0x020F
	s.cpu.a = 0x48
	s.bus[0x020F] = 0x2C
	s.bus[0x0210] = 0xB0
	s.bus[0x0211] = 0x03
	s.bus[0x03B0] = 0b11000000

	s.cpu.Step()

	s.EqualValues(4, s.cpu.cycles)
	s.EqualValues(0b11000000, s.cpu.p.u8())
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
			cpu := s.cpu

			cpu.pc = 0x020F
			cpu.a = tt.a

			s.bus[0x020F] = 0x6D
			s.bus[0x0210] = 0xD3
			s.bus[0x0211] = 0x04
			s.bus[0x04D3] = tt.m

			cpu.Step()

			s.EqualValues(tt.expectedA, cpu.a)
			s.EqualValues(tt.expectedP, cpu.p.u8())
		})
	}
}

func (s *executeTestSuite) Test_CPY() {
	s.cpu.pc = 0x020F
	s.cpu.y = 0x37
	s.bus[0x020F] = 0xCC
	s.bus[0x0210] = 0x36

	s.cpu.Step()

	s.EqualValues(4, s.cpu.cycles)
	s.EqualValues(0b00000001, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_INC() {
	s.cpu.pc = 0x020F
	s.bus[0x020F] = 0xEE
	s.bus[0x0210] = 0xD3
	s.bus[0x0211] = 0x04
	s.bus[0x04D3] = 0x7F

	s.cpu.Step()

	s.EqualValues(6, s.cpu.cycles)
	s.EqualValues(0x80, s.bus[0x04D3])
	s.EqualValues(0b10000000, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_DEC() {
	s.cpu.pc = 0x020F
	s.bus[0x020F] = 0xCE
	s.bus[0x0210] = 0xD3
	s.bus[0x0211] = 0x04
	s.bus[0x04D3] = 0xC0

	s.cpu.Step()

	s.EqualValues(6, s.cpu.cycles)
	s.EqualValues(0xBF, s.bus[0x04D3])
	s.EqualValues(0b10000000, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_ASL() {
	s.cpu.pc = 0x020F
	s.cpu.a = 0b10001010

	s.bus[0x020F] = 0x0A

	s.cpu.Step()

	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0b00010100, s.cpu.a)
	s.EqualValues(0b00000001, s.cpu.p.u8())
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
			cpu := s.cpu
			cpu.pc = 0x020F
			cpu.a = 0b10001010
			cpu.p.set(tt.p)

			cpu.Step()

			s.EqualValues(2, cpu.cycles)
			s.EqualValues(tt.expectedA, cpu.a)
			s.EqualValues(0b00000001, cpu.p.u8())
		})
	}
}

func (s *executeTestSuite) Test_JSR() {
	s.cpu.pc = 0x020F
	s.cpu.s = 0xBF

	s.bus[0x020F] = 0x20
	s.bus[0x0210] = 0x31
	s.bus[0x0211] = 0x40

	s.cpu.Step()

	s.EqualValues(0xBD, s.cpu.s)
	s.EqualValues(0x4031, s.cpu.pc)
	s.EqualValues(6, s.cpu.cycles)
	s.EqualValues(0x11, s.bus[0x01BE])
	s.EqualValues(0x02, s.bus[0x01BF])
}

func (s *executeTestSuite) Test_RTS() {
	s.cpu.pc = 0x0031
	s.cpu.s = 0xBD

	s.bus[0x0031] = 0x60
	s.bus[0x01BE] = 0x11
	s.bus[0x01BF] = 0x02

	s.cpu.Step()

	s.EqualValues(0xBF, s.cpu.s)
	s.EqualValues(0x0212, s.cpu.pc)
	s.EqualValues(6, s.cpu.cycles)
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
			cpu := s.cpu
			cpu.pc = 0x0031
			cpu.p.set(tt.p)

			s.bus[0x0031] = 0x90
			s.bus[0x0032] = tt.operand

			cpu.Step()

			s.EqualValues(tt.expectedPC, cpu.pc)
			s.EqualValues(tt.expectedCycles, cpu.cycles)
		})
	}
}

func (s *executeTestSuite) Test_CLD() {
	s.cpu.pc = 0x020F
	s.cpu.p.set(0b011001001)

	s.bus[0x020F] = 0xD8

	s.cpu.Step()

	s.EqualValues(0x0210, s.cpu.pc)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0b011000001, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_SEI() {
	s.cpu.pc = 0x020F
	s.cpu.p.set(0b011001001)

	s.bus[0x020F] = 0x78

	s.cpu.Step()

	s.EqualValues(0x0210, s.cpu.pc)
	s.EqualValues(2, s.cpu.cycles)
	s.EqualValues(0b011001101, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_BRK() {
	s.cpu.pc = 0x020F
	s.cpu.p.set(0b01100001)
	s.cpu.s = 0xBF

	s.bus[0x020F] = 0x00
	s.bus[0xFFFE] = 0x23
	s.bus[0xFFFF] = 0x40

	s.cpu.Step()

	s.EqualValues(0x4023, s.cpu.pc)
	s.EqualValues(7, s.cpu.cycles)
	s.EqualValues(0xBC, s.cpu.s)
	s.EqualValues(0b01000001, s.cpu.p.u8())
}

func (s *executeTestSuite) Test_RTI() {
	s.cpu.pc = 0x020F
	s.cpu.p.set(0b01100101)
	s.cpu.s = 0xBC

	s.bus[0x020F] = 0x40
	s.bus[0x01BD] = 0b10000010
	s.bus[0x01BE] = 0x11
	s.bus[0x01BF] = 0x02

	s.cpu.Step()

	s.EqualValues(0x0211, s.cpu.pc)
	s.EqualValues(6, s.cpu.cycles)
	s.EqualValues(0xBF, s.cpu.s)
	s.EqualValues(0b10000010, s.cpu.p.u8())
}

func Test_execute(t *testing.T) {
	suite.Run(t, new(executeTestSuite))
}
