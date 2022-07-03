package cpu

import "fmt"

// Trace is a snapshot of CPU state
type Trace struct {
	A, X, Y uint8
	S       uint8
	P       Status
	PC      uint16

	Cycles uint64

	Opcode   uint8
	Operand1 uint8
	Operand2 uint8

	Mnemonic       mnemonic
	AddressingMode addressingMode
}

// Trace current CPU state and return snapshot
func (e *CPU) Trace() Trace {
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
		A:              e.A,
		X:              e.X,
		Y:              e.Y,
		S:              e.S,
		P:              e.P,
		PC:             e.PC,
		Cycles:         e.Cycles,
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

func (t Trace) String() string {
	len := t.AddressingMode.instructionLength()
	var op string
	switch len {
	case 1:
		op = fmt.Sprintf("%02X      ", t.Opcode)
	case 2:
		op = fmt.Sprintf("%02X %02X   ", t.Opcode, t.Operand1)
	case 3:
		op = fmt.Sprintf("%02X %02X %02X", t.Opcode, t.Operand1, t.Operand2)
	default:
		panic("illegal addressing mode")
	}
	return fmt.Sprintf(
		"%04X  %s  A:%02X X:%02X Y:%02X P:%02X SP:%02X",
		t.PC,
		op,
		t.A,
		t.X,
		t.Y,
		t.P.u8(),
		t.S,
	)
}
