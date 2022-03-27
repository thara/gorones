package cpu

func bit(b bool) uint8 {
	if b {
		return 1
	}
	return 0
}

func nthBit(b, n uint8) uint8 {
	return (b >> n) & 1
}
