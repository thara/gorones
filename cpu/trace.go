package cpu

// Trace is a snapshot of CPU state
type Trace struct {
	Pc            uint16
	Opcode        uint8
	Operand1      uint8
	Operand2      uint8
	A, X, Y, S, P uint8

	Mnemonic       mnemonic
	AddressingMode addressingMode
	Cycles         uint64
}

// Trace current CPU state and return snapshot
func (c *CPU) Trace() Trace {
	op := c.m.ReadCPU(c.pc)
	inst := Decode(op)

	len := inst.AddressingMode.instructionLength()
	var op1, op2 uint8
	switch len {
	case 3:
		op2 = c.m.ReadCPU(c.pc + 2)
		fallthrough
	case 2:
		op1 = c.m.ReadCPU(c.pc + 1)
	}
	return Trace{
		Pc:             c.pc,
		Opcode:         op,
		Operand1:       op1,
		Operand2:       op2,
		A:              c.a,
		X:              c.x,
		Y:              c.y,
		S:              c.s,
		P:              c.p.u8() | 0x20, // for nestest
		Mnemonic:       inst.Mnemonic,
		AddressingMode: inst.AddressingMode,
		Cycles:         c.cycles,
	}
}

func (m addressingMode) instructionLength() uint8 {
	switch m {
	case immediate, zeroPage, zeroPageX, zeroPageY,
		relative, indirectIndexed, indexedIndirect, indirectIndexedWithPenalty:
		return 2
	case indirect, absolute, absoluteX, absoluteXWithPenalty, absoluteY, absoluteYWithPenalty:
		return 3
	}
	return 1
}
