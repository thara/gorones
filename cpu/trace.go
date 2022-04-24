package cpu

// Trace is a snapshot of CPU state
type Trace struct {
	CPU

	Opcode   uint8
	Operand1 uint8
	Operand2 uint8

	Mnemonic       mnemonic
	AddressingMode addressingMode
}

// Trace current CPU state and return snapshot
func (e *Emu) Trace() Trace {
	op := e.m.ReadCPU(e.cpu.PC)
	inst := Decode(op)

	len := inst.AddressingMode.instructionLength()
	var op1, op2 uint8
	switch len {
	case 3:
		op2 = e.m.ReadCPU(e.cpu.PC + 2)
		fallthrough
	case 2:
		op1 = e.m.ReadCPU(e.cpu.PC + 1)
	}
	return Trace{
		CPU:            e.cpu,
		Opcode:         op,
		Operand1:       op1,
		Operand2:       op2,
		Mnemonic:       inst.Mnemonic,
		AddressingMode: inst.AddressingMode,
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
