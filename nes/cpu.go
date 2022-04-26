package nes

func (n *NES) fetch() uint8 {
	op := n.read(n.cpu.PC)
	n.cpu.PC += 1
	return op
}

func (n *NES) handleInterrupt() {
	if n.interrupt == nil {
		return
	}

	intr := *n.interrupt
	switch intr {
	case NMI: // interrupted
	case IRQ:
		if !n.cpu.P[status_I] {
			return
		}
		// interrupted
	default:
		return
	}

	n.tick_n(2)
	n.pushStackWord(n.cpu.PC)
	// https://wiki.nesdev.com/w/index.php/Status_flags#The_B_flag
	// http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit
	n.pushStack(n.cpu.P.u8() | interruptB)
	n.cpu.P[status_I] = true
	n.cpu.PC = n.readWord(intr.vector())
	n.interrupt = nil
}

func (n *NES) pushStack(v uint8) {
	n.write(uint16(n.cpu.S)+0x0100, v)
	n.cpu.S -= 1
}

func (n *NES) pushStackWord(v uint16) {
	n.pushStack(uint8(v >> 8))
	n.pushStack(uint8(v & 0xFF))
}

func (n *NES) pullStack() uint8 {
	n.cpu.S += 1
	return n.read(uint16(n.cpu.S) + 0x0100)
}

func (n *NES) pullStackWord() uint16 {
	return uint16(n.pullStack()) | uint16(n.pullStack())<<8
}
