package apu

type lengthCounter struct {
	count uint
	halt  bool
}

func (c *lengthCounter) reset() {
	c.count = 0
}

func (c *lengthCounter) reload(v uint8) {
	c.count = lengthTable[v]
}

func (c *lengthCounter) clock() {
	if 0 < c.count && !c.halt {
		c.count -= 1
	}
}

var lengthTable = []uint{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}
