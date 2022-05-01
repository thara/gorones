package cpu

import "github.com/thara/gorones/util"

var bit = util.Bit

func pageCrossed[T ~uint16 | ~int16](a, b T) bool {
	var p int = 0xFF00
	return (a+b)&T(p) != (b & T(p))
}
