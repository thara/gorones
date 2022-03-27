package cpu

// https://wiki.nesdev.org/w/index.php?title=CPU_addressing_modes

// addressingMode for 6502
type addressingMode uint8

const (
	_ addressingMode = iota
	implicit
	accumulator
	immediate
	zeroPage
	zeroPageX
	zeroPageY
	absolute
	absoluteX
	absoluteXWithPenalty
	absoluteY
	absoluteYWithPenalty
	relative
	indirect
	indexedIndirect
	indirectIndexed
	indirectIndexedWithPenalty
)

func (c *CPU) getOperand(m addressingMode) uint16 {
	switch m {
	case implicit:
		return 0
	case accumulator:
		return uint16(c.a)
	case immediate:
		pc := c.pc
		c.pc++
		return pc
	case zeroPage:
		v := c.read(c.pc)
		c.pc++
		return uint16(v)
	case zeroPageX:
		c.tick()
		v := (uint16(c.read(c.pc)) + uint16(c.x)) & 0xFF
		c.pc++
		return uint16(v)
	case zeroPageY:
		c.tick()
		v := (uint16(c.read(c.pc)) + uint16(c.y)) & 0xFF
		c.pc++
		return uint16(v)
	case absolute:
		v := c.readWord(c.pc)
		c.pc += 2
		return v
	case absoluteX:
		v := c.readWord(c.pc)
		c.pc += 2
		c.tick()
		return v + uint16(c.x)
	case absoluteXWithPenalty:
		v := c.readWord(c.pc)
		c.pc += 2
		if pageCrossed(uint16(c.x), v) {
			c.tick()
		}
		return v + uint16(c.x)
	case absoluteY:
		v := c.readWord(c.pc)
		c.pc += 2
		c.tick()
		return v + uint16(c.y)
	case absoluteYWithPenalty:
		v := c.readWord(c.pc)
		c.pc += 2
		if pageCrossed(uint16(c.y), v) {
			c.tick()
		}
		return v + uint16(c.y)
	case relative:
		v := c.read(c.pc)
		c.pc++
		return uint16(v)
	case indirect:
		m := c.readWord(c.pc)
		v := c.readOnIndirect(m)
		c.pc += 2
		return v
	case indexedIndirect:
		m := c.read(c.pc)
		v := c.readOnIndirect(uint16(m + c.x))
		c.pc += 1
		c.tick()
		return v
	case indirectIndexed:
		m := c.read(c.pc)
		v := c.readOnIndirect(uint16(m))
		c.pc += 1
		c.tick()
		return v + uint16(c.y)
	case indirectIndexedWithPenalty:
		m := c.read(c.pc)
		v := c.readOnIndirect(uint16(m))
		c.pc += 1
		if pageCrossed(uint16(c.y), v) {
			c.tick()
		}
		return v + uint16(c.y)
	}

	panic("unrecognized addressing mode")
}
