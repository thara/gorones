package cpu

import "fmt"

// interrupt Kinds of CPU interrupts
type interrupt uint8

// currently supports NMI and IRQ only
const (
	_ interrupt = iota
	NMI
	IRQ
)

func (i interrupt) vector() uint16 {
	switch i {
	case NMI:
		return 0xFFFA
	case IRQ:
		return 0xFFFE
	}
	panic(fmt.Sprintf("unsupported interrupt : %d", i))
}

func (c *CPU) handleInterrupt() {
	if c.interrupt == nil {
		return
	}

	intr := *c.interrupt
	switch intr {
	case NMI: // interrupted
	case IRQ:
		if !c.p[status_I] {
			return
		}
		// interrupted
	default:
		return
	}

	c.tick_n(2)
	c.pushStackWord(c.pc)
	// https://wiki.nesdev.com/w/index.php/Status_flags#The_B_flag
	// http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit
	c.pushStack(c.p.u8() | interruptB)
	c.p[status_I] = true
	c.pc = c.readWord(intr.vector())
	c.interrupt = nil
}
