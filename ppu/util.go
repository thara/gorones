package ppu

import (
	"golang.org/x/exp/constraints"

	"github.com/thara/gorones/util"
)

func nth[T constraints.Unsigned, N constraints.Integer](b T, n N) uint8 {
	return uint8(util.NthBit(b, n))
}
