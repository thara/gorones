package ppu

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
