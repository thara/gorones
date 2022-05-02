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

func NthBit[T constraints.Unsigned, N constraints.Integer](b T, n N) T {
	return (b >> n) & 1
}

func IsSet[T constraints.Unsigned, N constraints.Integer](b T, n N) bool {
	return (b>>n)&1 == 1
}
