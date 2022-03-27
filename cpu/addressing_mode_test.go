package cpu

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type getOperandTestSuite struct {
	suite.Suite
	bus []uint8
	cpu CPU
}

func (s *getOperandTestSuite) SetupTest() {
	s.bus = newBusMock()
	s.cpu = New(tickMock, busMock(s.bus))
}

func (s *getOperandTestSuite) Test_implicit() {
	v := s.cpu.getOperand(implicit)
	s.EqualValues(0, v)
	s.EqualValues(0, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_accumulator() {
	s.cpu.a = 0xFB

	v := s.cpu.getOperand(accumulator)
	s.EqualValues(0xFB, v)
	s.EqualValues(0, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_immediate() {
	s.cpu.pc = 0x8234

	v := s.cpu.getOperand(immediate)
	s.EqualValues(0x8234, v)
	s.EqualValues(0, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_zeroPage() {
	s.cpu.pc = 0x0414
	s.bus[0x0414] = 0x91

	v := s.cpu.getOperand(zeroPage)
	s.EqualValues(0x91, v)
	s.EqualValues(1, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_zeroPageX() {
	s.cpu.pc = 0x0100
	s.cpu.x = 0x93
	s.bus[0x0100] = 0x80

	v := s.cpu.getOperand(zeroPageX)
	s.EqualValues(0x13, v)
	s.EqualValues(2, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_zeroPageY() {
	s.cpu.pc = 0x0423
	s.cpu.y = 0xF1
	s.bus[0x0423] = 0x36

	v := s.cpu.getOperand(zeroPageY)
	s.EqualValues(0x27, v)
	s.EqualValues(2, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_absolute() {
	s.cpu.pc = 0x0423
	s.bus[0x0423] = 0x36
	s.bus[0x0424] = 0xF0

	v := s.cpu.getOperand(absolute)
	s.EqualValues(0xF036, v)
	s.EqualValues(2, s.cpu.cycles)
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
			cpu := s.cpu
			cpu.pc = 0x0423
			cpu.x = tt.x

			v := cpu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, cpu.cycles)
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
			cpu := s.cpu
			cpu.pc = 0x0423
			cpu.y = tt.y

			v := cpu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, cpu.cycles)
		})
	}
}

func (s *getOperandTestSuite) Test_relative() {
	s.cpu.pc = 0x0414
	s.bus[0x0414] = 0x91

	v := s.cpu.getOperand(relative)
	s.EqualValues(0x91, v)
	s.EqualValues(1, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_indirect() {
	s.cpu.pc = 0x020F
	s.bus[0x020F] = 0x10
	s.bus[0x0210] = 0x03
	s.bus[0x0310] = 0x9F

	v := s.cpu.getOperand(indirect)
	s.EqualValues(0x9F, v)
	s.EqualValues(4, s.cpu.cycles)
}

func (s *getOperandTestSuite) Test_indexedIndirect() {
	s.cpu.pc = 0x020F
	s.cpu.x = 0x95
	s.bus[0x020F] = 0xF0
	s.bus[0x0085] = 0x12
	s.bus[0x0086] = 0x90

	v := s.cpu.getOperand(indexedIndirect)
	s.EqualValues(0x9012, v)
	s.EqualValues(4, s.cpu.cycles)
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
			cpu := s.cpu
			cpu.pc = 0x020F
			cpu.y = tt.y

			v := cpu.getOperand(tt.mode)
			s.EqualValues(tt.expectedOperand, v)
			s.EqualValues(tt.expectedCycles, cpu.cycles)
		})
	}
}

func Test_getOperand(t *testing.T) {
	suite.Run(t, new(getOperandTestSuite))
}
