package cpu

import "fmt"

// Interrupt Kinds of CPU interrupts
type Interrupt uint8

// currently supports NMI and IRQ only
const (
	NoInterrupt Interrupt = iota
	NMI
	IRQ
)

func (i Interrupt) vector() uint16 {
	switch i {
	case NMI:
		return 0xFFFA
	case IRQ:
		return 0xFFFE
	}
	panic(fmt.Sprintf("unsupported interrupt : %d", i))
}

func (c *CPU) handleInterrupt(intr *Interrupt) {
	if intr == nil {
		return
	}

	switch *intr {
	case NoInterrupt:
		return
	case NMI: // interrupted
	case IRQ:
		if !c.P[status_I] {
			return
		}
		// interrupted
	default:
		return
	}

	c.tick_n(2)
	c.pushStackWord(c.PC)
	// https://wiki.nesdev.com/w/index.php/Status_flags#The_B_flag
	// http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit
	c.pushStack(c.P.u8() | interruptB)
	c.P[status_I] = true
	c.PC = c.readWord(intr.vector())
	*intr = NoInterrupt
}
