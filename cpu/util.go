package cpu

func bit(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func pageCrossed[T ~uint16 | ~int16](a, b T) bool {
	var p int = 0xFF00
	return (a+b)&T(p) != (b & T(p))
}

// func nthBit(b, n uint8) uint8 {
// 	return (b >> n) & 1
// }
