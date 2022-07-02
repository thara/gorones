package apu

type lengthCounter struct {
	count uint8
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
