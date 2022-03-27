package cpu

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type statusTestSuite struct {
	suite.Suite

	p status
}

func (s *statusTestSuite) SetupTest() {
	var p status
	s.p = p
}

func (s *statusTestSuite) Test_u8() {
	s.EqualValues(0, s.p.u8())

	s.p[status_C] = true
	s.EqualValues(0b00000001, s.p.u8())
	s.p[status_Z] = true
	s.EqualValues(0b00000011, s.p.u8())
	s.p[status_I] = true
	s.EqualValues(0b00000111, s.p.u8())
	s.p[status_D] = true
	s.EqualValues(0b00001111, s.p.u8())
	s.p[status_V] = true
	s.EqualValues(0b01001111, s.p.u8())
	s.p[status_N] = true
	s.EqualValues(0b11001111, s.p.u8())
}

func (s *statusTestSuite) Test_set() {
	s.p.set(0)
	s.EqualValues(0, s.p.u8())

	s.p.set(0b00000001)
	s.EqualValues(0b00000001, s.p.u8())
	s.p.set(0b00000011)
	s.EqualValues(0b00000011, s.p.u8())
	s.p.set(0b00000111)
	s.EqualValues(0b00000111, s.p.u8())
	s.p.set(0b00001111)
	s.EqualValues(0b00001111, s.p.u8())
	s.p.set(0b00011111)
	s.EqualValues(0b00001111, s.p.u8())
	s.p.set(0b00111111)
	s.EqualValues(0b00001111, s.p.u8())
	s.p.set(0b01111111)
	s.EqualValues(0b01001111, s.p.u8())
	s.p.set(0b11111111)
	s.EqualValues(0b11001111, s.p.u8())
}

func (s *statusTestSuite) Test_insert() {
	s.p.set(0b00000001)

	s.p.insert(0b11000000)
	s.EqualValues(0b11000001, s.p.u8())
}

func (s *statusTestSuite) Test_setZN() {
	s.p.set(0b01000000)

	s.p.setZN(0)
	s.EqualValues(0b01000010, s.p.u8())

	s.p.setZN(0b10000000)
	s.EqualValues(0b11000000, s.p.u8())
}

func Test_status(t *testing.T) {
	suite.Run(t, new(statusTestSuite))
}
