package util

import (
	"golang.org/x/exp/constraints"
)

func Bit(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func NthBit[T constraints.Unsigned](b T, n int) T {
	return (b >> n) & 1
}
