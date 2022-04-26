package cpu

import (
	"fmt"
)

// https://www.nesdev.org/wiki/6502_instructions

type instruction struct {
	Mnemonic       mnemonic
	AddressingMode addressingMode
}

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

type mnemonic uint8

const (
	_ mnemonic = iota

	// Load/Store Operations
	LDA
	LDX
	LDY
	STA
	STX
	STY
	// Register Operations
	TAX
	TSX
	TAY
	TXA
	TXS
	TYA
	// Stack instructions
	PHA
	PHP
	PLA
	PLP
	// Logical instructions
	AND
	EOR
	ORA
	BIT
	// Arithmetic instructions
	ADC
	SBC
	CMP
	CPX
	CPY
	// Increment/Decrement instructions
	INC
	INX
	INY
	DEC
	DEX
	DEY
	// Shift instructions
	ASL
	LSR
	ROL
	ROR
	// Jump instructions
	JMP
	JSR
	RTS
	RTI
	// Branch instructions
	BCC
	BCS
	BEQ
	BMI
	BNE
	BPL
	BVC
	BVS
	// Flag control instructions
	CLC
	CLD
	CLI
	CLV
	SEC
	SED
	SEI
	// Misc
	BRK
	NOP
	// Unofficial
	LAX
	SAX
	DCP
	ISB
	SLO
	RLA
	SRE
	RRA
)

func Decode(opcode uint8) instruction {
	switch opcode {
	case 0x69:
		return instruction{ADC, immediate}
	case 0x65:
		return instruction{ADC, zeroPage}
	case 0x75:
		return instruction{ADC, zeroPageX}
	case 0x6D:
		return instruction{ADC, absolute}
	case 0x7D:
		return instruction{ADC, absoluteXWithPenalty}
	case 0x79:
		return instruction{ADC, absoluteYWithPenalty}
	case 0x61:
		return instruction{ADC, indexedIndirect}
	case 0x71:
		return instruction{ADC, indirectIndexedWithPenalty}

	case 0x29:
		return instruction{AND, immediate}
	case 0x25:
		return instruction{AND, zeroPage}
	case 0x35:
		return instruction{AND, zeroPageX}
	case 0x2D:
		return instruction{AND, absolute}
	case 0x3D:
		return instruction{AND, absoluteXWithPenalty}
	case 0x39:
		return instruction{AND, absoluteYWithPenalty}
	case 0x21:
		return instruction{AND, indexedIndirect}
	case 0x31:
		return instruction{AND, indirectIndexedWithPenalty}

	case 0x0A:
		return instruction{ASL, accumulator}
	case 0x06:
		return instruction{ASL, zeroPage}
	case 0x16:
		return instruction{ASL, zeroPageX}
	case 0x0E:
		return instruction{ASL, absolute}
	case 0x1E:
		return instruction{ASL, absoluteX}

	case 0x90:
		return instruction{BCC, relative}
	case 0xB0:
		return instruction{BCS, relative}
	case 0xF0:
		return instruction{BEQ, relative}

	case 0x24:
		return instruction{BIT, zeroPage}
	case 0x2C:
		return instruction{BIT, absolute}

	case 0x30:
		return instruction{BMI, relative}
	case 0xD0:
		return instruction{BNE, relative}
	case 0x10:
		return instruction{BPL, relative}

	case 0x00:
		return instruction{BRK, implicit}

	case 0x50:
		return instruction{BVC, relative}
	case 0x70:
		return instruction{BVS, relative}

	case 0x18:
		return instruction{CLC, implicit}
	case 0xD8:
		return instruction{CLD, implicit}
	case 0x58:
		return instruction{CLI, implicit}
	case 0xB8:
		return instruction{CLV, implicit}

	case 0xC9:
		return instruction{CMP, immediate}
	case 0xC5:
		return instruction{CMP, zeroPage}
	case 0xD5:
		return instruction{CMP, zeroPageX}
	case 0xCD:
		return instruction{CMP, absolute}
	case 0xDD:
		return instruction{CMP, absoluteXWithPenalty}
	case 0xD9:
		return instruction{CMP, absoluteYWithPenalty}
	case 0xC1:
		return instruction{CMP, indexedIndirect}
	case 0xD1:
		return instruction{CMP, indirectIndexedWithPenalty}

	case 0xE0:
		return instruction{CPX, immediate}
	case 0xE4:
		return instruction{CPX, zeroPage}
	case 0xEC:
		return instruction{CPX, absolute}
	case 0xC0:
		return instruction{CPY, immediate}
	case 0xC4:
		return instruction{CPY, zeroPage}
	case 0xCC:
		return instruction{CPY, absolute}

	case 0xC6:
		return instruction{DEC, zeroPage}
	case 0xD6:
		return instruction{DEC, zeroPageX}
	case 0xCE:
		return instruction{DEC, absolute}
	case 0xDE:
		return instruction{DEC, absoluteX}

	case 0xCA:
		return instruction{DEX, implicit}
	case 0x88:
		return instruction{DEY, implicit}

	case 0x49:
		return instruction{EOR, immediate}
	case 0x45:
		return instruction{EOR, zeroPage}
	case 0x55:
		return instruction{EOR, zeroPageX}
	case 0x4D:
		return instruction{EOR, absolute}
	case 0x5D:
		return instruction{EOR, absoluteXWithPenalty}
	case 0x59:
		return instruction{EOR, absoluteYWithPenalty}
	case 0x41:
		return instruction{EOR, indexedIndirect}
	case 0x51:
		return instruction{EOR, indirectIndexedWithPenalty}

	case 0xE6:
		return instruction{INC, zeroPage}
	case 0xF6:
		return instruction{INC, zeroPageX}
	case 0xEE:
		return instruction{INC, absolute}
	case 0xFE:
		return instruction{INC, absoluteX}

	case 0xE8:
		return instruction{INX, implicit}
	case 0xC8:
		return instruction{INY, implicit}

	case 0x4C:
		return instruction{JMP, absolute}
	case 0x6C:
		return instruction{JMP, indirect}

	case 0x20:
		return instruction{JSR, absolute}

	case 0xA9:
		return instruction{LDA, immediate}
	case 0xA5:
		return instruction{LDA, zeroPage}
	case 0xB5:
		return instruction{LDA, zeroPageX}
	case 0xAD:
		return instruction{LDA, absolute}
	case 0xBD:
		return instruction{LDA, absoluteXWithPenalty}
	case 0xB9:
		return instruction{LDA, absoluteYWithPenalty}
	case 0xA1:
		return instruction{LDA, indexedIndirect}
	case 0xB1:
		return instruction{LDA, indirectIndexedWithPenalty}

	case 0xA2:
		return instruction{LDX, immediate}
	case 0xA6:
		return instruction{LDX, zeroPage}
	case 0xB6:
		return instruction{LDX, zeroPageY}
	case 0xAE:
		return instruction{LDX, absolute}
	case 0xBE:
		return instruction{LDX, absoluteYWithPenalty}

	case 0xA0:
		return instruction{LDY, immediate}
	case 0xA4:
		return instruction{LDY, zeroPage}
	case 0xB4:
		return instruction{LDY, zeroPageX}
	case 0xAC:
		return instruction{LDY, absolute}
	case 0xBC:
		return instruction{LDY, absoluteXWithPenalty}

	case 0x4A:
		return instruction{LSR, accumulator}
	case 0x46:
		return instruction{LSR, zeroPage}
	case 0x56:
		return instruction{LSR, zeroPageX}
	case 0x4E:
		return instruction{LSR, absolute}
	case 0x5E:
		return instruction{LSR, absoluteX}

	case 0x09:
		return instruction{ORA, immediate}
	case 0x05:
		return instruction{ORA, zeroPage}
	case 0x15:
		return instruction{ORA, zeroPageX}
	case 0x0D:
		return instruction{ORA, absolute}
	case 0x1D:
		return instruction{ORA, absoluteXWithPenalty}
	case 0x19:
		return instruction{ORA, absoluteYWithPenalty}
	case 0x01:
		return instruction{ORA, indexedIndirect}
	case 0x11:
		return instruction{ORA, indirectIndexedWithPenalty}

	case 0x48:
		return instruction{PHA, implicit}
	case 0x08:
		return instruction{PHP, implicit}
	case 0x68:
		return instruction{PLA, implicit}
	case 0x28:
		return instruction{PLP, implicit}

	case 0x2A:
		return instruction{ROL, accumulator}
	case 0x26:
		return instruction{ROL, zeroPage}
	case 0x36:
		return instruction{ROL, zeroPageX}
	case 0x2E:
		return instruction{ROL, absolute}
	case 0x3E:
		return instruction{ROL, absoluteX}

	case 0x6A:
		return instruction{ROR, accumulator}
	case 0x66:
		return instruction{ROR, zeroPage}
	case 0x76:
		return instruction{ROR, zeroPageX}
	case 0x6E:
		return instruction{ROR, absolute}
	case 0x7E:
		return instruction{ROR, absoluteX}

	case 0x40:
		return instruction{RTI, implicit}
	case 0x60:
		return instruction{RTS, implicit}

	case 0xE9:
		return instruction{SBC, immediate}
	case 0xE5:
		return instruction{SBC, zeroPage}
	case 0xF5:
		return instruction{SBC, zeroPageX}
	case 0xED:
		return instruction{SBC, absolute}
	case 0xFD:
		return instruction{SBC, absoluteXWithPenalty}
	case 0xF9:
		return instruction{SBC, absoluteYWithPenalty}
	case 0xE1:
		return instruction{SBC, indexedIndirect}
	case 0xF1:
		return instruction{SBC, indirectIndexedWithPenalty}

	case 0x38:
		return instruction{SEC, implicit}
	case 0xF8:
		return instruction{SED, implicit}
	case 0x78:
		return instruction{SEI, implicit}

	case 0x85:
		return instruction{STA, zeroPage}
	case 0x95:
		return instruction{STA, zeroPageX}
	case 0x8D:
		return instruction{STA, absolute}
	case 0x9D:
		return instruction{STA, absoluteX}
	case 0x99:
		return instruction{STA, absoluteY}
	case 0x81:
		return instruction{STA, indexedIndirect}
	case 0x91:
		return instruction{STA, indirectIndexed}

	case 0x86:
		return instruction{STX, zeroPage}
	case 0x96:
		return instruction{STX, zeroPageY}
	case 0x8E:
		return instruction{STX, absolute}
	case 0x84:
		return instruction{STY, zeroPage}
	case 0x94:
		return instruction{STY, zeroPageX}
	case 0x8C:
		return instruction{STY, absolute}

	case 0xAA:
		return instruction{TAX, implicit}
	case 0xA8:
		return instruction{TAY, implicit}
	case 0xBA:
		return instruction{TSX, implicit}
	case 0x8A:
		return instruction{TXA, implicit}
	case 0x9A:
		return instruction{TXS, implicit}
	case 0x98:
		return instruction{TYA, implicit}

	case 0x04, 0x44, 0x64:
		return instruction{NOP, zeroPage}
	case 0x0C:
		return instruction{NOP, absolute}
	case 0x14, 0x34, 0x54, 0x74, 0xD4, 0xF4:
		return instruction{NOP, zeroPageX}
	case 0x1A, 0x3A, 0x5A, 0x7A, 0xDA, 0xEA, 0xFA:
		return instruction{NOP, implicit}
	case 0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC:
		return instruction{NOP, absoluteXWithPenalty}
	case 0x80, 0x82, 0x89, 0xc2, 0xE2:
		return instruction{NOP, immediate}

	// unofficial
	case 0xEB:
		return instruction{SBC, immediate}

	case 0xA3:
		return instruction{LAX, indexedIndirect}
	case 0xA7:
		return instruction{LAX, zeroPage}
	case 0xAB:
		return instruction{LAX, immediate}
	case 0xAF:
		return instruction{LAX, absolute}
	case 0xB3:
		return instruction{LAX, indirectIndexedWithPenalty}
	case 0xB7:
		return instruction{LAX, zeroPageY}
	case 0xBF:
		return instruction{LAX, absoluteYWithPenalty}

	case 0x83:
		return instruction{SAX, indexedIndirect}
	case 0x87:
		return instruction{SAX, zeroPage}
	case 0x8F:
		return instruction{SAX, absolute}
	case 0x97:
		return instruction{SAX, zeroPageY}

	case 0xC3:
		return instruction{DCP, indexedIndirect}
	case 0xC7:
		return instruction{DCP, zeroPage}
	case 0xCF:
		return instruction{DCP, absolute}
	case 0xD3:
		return instruction{DCP, indirectIndexed}
	case 0xD7:
		return instruction{DCP, zeroPageX}
	case 0xDB:
		return instruction{DCP, absoluteY}
	case 0xDF:
		return instruction{DCP, absoluteX}

	case 0xE3:
		return instruction{ISB, indexedIndirect}
	case 0xE7:
		return instruction{ISB, zeroPage}
	case 0xEF:
		return instruction{ISB, absolute}
	case 0xF3:
		return instruction{ISB, indirectIndexed}
	case 0xF7:
		return instruction{ISB, zeroPageX}
	case 0xFB:
		return instruction{ISB, absoluteY}
	case 0xFF:
		return instruction{ISB, absoluteX}

	case 0x03:
		return instruction{SLO, indexedIndirect}
	case 0x07:
		return instruction{SLO, zeroPage}
	case 0x0F:
		return instruction{SLO, absolute}
	case 0x13:
		return instruction{SLO, indirectIndexed}
	case 0x17:
		return instruction{SLO, zeroPageX}
	case 0x1B:
		return instruction{SLO, absoluteY}
	case 0x1F:
		return instruction{SLO, absoluteX}

	case 0x23:
		return instruction{RLA, indexedIndirect}
	case 0x27:
		return instruction{RLA, zeroPage}
	case 0x2F:
		return instruction{RLA, absolute}
	case 0x33:
		return instruction{RLA, indirectIndexed}
	case 0x37:
		return instruction{RLA, zeroPageX}
	case 0x3B:
		return instruction{RLA, absoluteY}
	case 0x3F:
		return instruction{RLA, absoluteX}

	case 0x43:
		return instruction{SRE, indexedIndirect}
	case 0x47:
		return instruction{SRE, zeroPage}
	case 0x4F:
		return instruction{SRE, absolute}
	case 0x53:
		return instruction{SRE, indirectIndexed}
	case 0x57:
		return instruction{SRE, zeroPageX}
	case 0x5B:
		return instruction{SRE, absoluteY}
	case 0x5F:
		return instruction{SRE, absoluteX}

	case 0x63:
		return instruction{RRA, indexedIndirect}
	case 0x67:
		return instruction{RRA, zeroPage}
	case 0x6F:
		return instruction{RRA, absolute}
	case 0x73:
		return instruction{RRA, indirectIndexed}
	case 0x77:
		return instruction{RRA, zeroPageX}
	case 0x7B:
		return instruction{RRA, absoluteY}
	case 0x7F:
		return instruction{RRA, absoluteX}

	default:
		return instruction{NOP, implicit}
	}
}

func (c *CPU) getOperand(m addressingMode) uint16 {
	switch m {
	case implicit:
		return 0
	case accumulator:
		return uint16(c.A)
	case immediate:
		pc := c.PC
		c.PC++
		return pc
	case zeroPage:
		v := c.read(c.PC)
		c.PC++
		return uint16(v)
	case zeroPageX:
		c.tick()
		v := (uint16(c.read(c.PC)) + uint16(c.X)) & 0xFF
		c.PC++
		return uint16(v)
	case zeroPageY:
		c.tick()
		v := (uint16(c.read(c.PC)) + uint16(c.Y)) & 0xFF
		c.PC++
		return uint16(v)
	case absolute:
		v := c.readWord(c.PC)
		c.PC += 2
		return v
	case absoluteX:
		v := c.readWord(c.PC)
		c.PC += 2
		c.tick()
		return v + uint16(c.X)
	case absoluteXWithPenalty:
		v := c.readWord(c.PC)
		c.PC += 2
		if pageCrossed(uint16(c.X), v) {
			c.tick()
		}
		return v + uint16(c.X)
	case absoluteY:
		v := c.readWord(c.PC)
		c.PC += 2
		c.tick()
		return v + uint16(c.Y)
	case absoluteYWithPenalty:
		v := c.readWord(c.PC)
		c.PC += 2
		if pageCrossed(uint16(c.Y), v) {
			c.tick()
		}
		return v + uint16(c.Y)
	case relative:
		v := c.read(c.PC)
		c.PC++
		return uint16(v)
	case indirect:
		m := c.readWord(c.PC)
		v := c.readOnIndirect(m)
		c.PC += 2
		return v
	case indexedIndirect:
		m := c.read(c.PC)
		v := c.readOnIndirect(uint16(m + c.X))
		c.PC += 1
		c.tick()
		return v
	case indirectIndexed:
		m := c.read(c.PC)
		v := c.readOnIndirect(uint16(m))
		c.PC += 1
		c.tick()
		return v + uint16(c.Y)
	case indirectIndexedWithPenalty:
		m := c.read(c.PC)
		v := c.readOnIndirect(uint16(m))
		c.PC += 1
		if pageCrossed(uint16(c.Y), v) {
			c.tick()
		}
		return v + uint16(c.Y)
	}

	panic("unrecognized addressing mode")
}

func (c *CPU) execute(inst instruction) {
	v := c.getOperand(inst.AddressingMode)

	switch inst.Mnemonic {
	case LDA:
		c.A = c.read(v)
		c.P.setZN(c.A)
	case LDX:
		c.X = c.read(v)
		c.P.setZN(c.X)
	case LDY:
		c.Y = c.read(v)
		c.P.setZN(c.Y)
	case STA:
		c.write(v, c.A)
	case STX:
		c.write(v, c.X)
	case STY:
		c.write(v, c.Y)
	case TAX:
		c.X = c.A
		c.P.setZN(c.X)
		c.tick()
	case TAY:
		c.Y = c.A
		c.P.setZN(c.Y)
		c.tick()
	case TXA:
		c.A = c.X
		c.P.setZN(c.A)
		c.tick()
	case TYA:
		c.A = c.Y
		c.P.setZN(c.A)
		c.tick()
	case TSX:
		c.X = c.S
		c.P.setZN(c.X)
		c.tick()
	case TXS:
		c.S = c.X
		c.tick()

	case PHA:
		c.pushStack(c.A)
		c.tick()
	case PHP:
		p := c.P.u8() | instructionB
		c.pushStack(p)
		c.tick()
	case PLA:
		c.A = c.pullStack()
		c.P.setZN(c.A)
		c.tick_n(2)
	case PLP:
		v := c.pullStack() & ^instructionB
		v |= 0b100000 // for nestest
		c.P.Set(v)
		c.tick_n(2)

	case AND:
		c.and(v)
	case EOR:
		c.eor(v)
	case ORA:
		c.ora(v)
	case BIT:
		m := c.read(v)
		b := c.A & m
		c.P[status_Z] = b == 0
		c.P[status_V] = m&0x40 == 0x40
		c.P[status_N] = m&0x80 == 0x80

	case ADC:
		m := c.read(v)
		r := c.A + m
		if c.P[status_C] {
			r += 1
		}
		c.carry(m, r)
		c.A = r
		c.P.setZN(c.A)
	case SBC:
		c.sbc(v)
	case CMP:
		c.cmp(c.A, v)
	case CPX:
		c.cmp(c.X, v)
	case CPY:
		c.cmp(c.Y, v)

	case INC:
		m := c.read(v)
		r := m + 1
		c.write(v, r)
		c.P.setZN(r)
		c.tick()
	case INX:
		c.X += 1
		c.P.setZN(c.X)
		c.tick()
	case INY:
		c.Y += 1
		c.P.setZN(c.Y)
		c.tick()
	case DEC:
		m := c.read(v)
		r := m - 1
		c.write(v, r)
		c.P.setZN(r)
		c.tick()
	case DEX:
		c.X -= 1
		c.P.setZN(c.X)
		c.tick()
	case DEY:
		c.Y -= 1
		c.P.setZN(c.Y)
		c.tick()

	case ASL:
		asl := func(m *uint8) {
			c.P[status_C] = *m&0x80 == 0x80
			*m <<= 1
			c.P.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			asl(&c.A)
			return
		}
		m := c.read(v)
		asl(&m)
		c.write(v, m)

	case LSR:
		lsr := func(m *uint8) {
			c.P[status_C] = *m&1 == 1
			*m >>= 1
			c.P.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			lsr(&c.A)
			return
		}
		m := c.read(v)
		lsr(&m)
		c.write(v, m)
	case ROL:
		rol := func(m *uint8) {
			carry := *m & 0x80
			*m <<= 1
			if c.P[status_C] {
				*m |= 1
			}
			c.P[status_C] = carry == 0x80
			c.P.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			rol(&c.A)
			return
		}
		m := c.read(v)
		rol(&m)
		c.write(v, m)
	case ROR:
		ror := func(m *uint8) {
			carry := *m & 1
			*m >>= 1
			if c.P[status_C] {
				*m |= 0x80
			}
			c.P[status_C] = carry == 1
			c.P.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			ror(&c.A)
			return
		}
		m := c.read(v)
		ror(&m)
		c.write(v, m)

	case JMP:
		c.PC = v
	case JSR:
		rtn := c.PC - 1
		c.pushStackWord(rtn)
		c.PC = v
		c.tick()
	case RTS:
		c.PC = c.pullStackWord()
		c.PC += 1
		c.tick_n(3)

	case BCC:
		c.branch(v, !c.P[status_C])
	case BCS:
		c.branch(v, c.P[status_C])
	case BEQ:
		c.branch(v, c.P[status_Z])
	case BMI:
		c.branch(v, c.P[status_N])
	case BNE:
		c.branch(v, !c.P[status_Z])
	case BPL:
		c.branch(v, !c.P[status_N])
	case BVC:
		c.branch(v, !c.P[status_V])
	case BVS:
		c.branch(v, c.P[status_V])

	case CLC:
		c.P[status_C] = false
		c.tick()
	case CLD:
		c.P[status_D] = false
		c.tick()
	case CLI:
		c.P[status_I] = false
		c.tick()
	case CLV:
		c.P[status_V] = false
		c.tick()
	case SEC:
		c.P[status_C] = true
		c.tick()
	case SED:
		c.P[status_D] = true
		c.tick()
	case SEI:
		c.P[status_I] = true
		c.tick()

	case BRK:
		c.pushStackWord(c.PC)
		c.P.insert(instructionB)
		c.pushStack(c.P.u8())
		c.PC = c.readWord(0xFFFE)
		c.tick()
	case NOP:
		c.tick()
	case RTI:
		v := c.pullStack()
		c.P.Set(v)
		c.PC = c.pullStackWord()
		c.tick_n(2)

	case LAX:
		m := c.read(v)
		c.A = m
		c.P.setZN(m)
		c.X = m
	case SAX:
		c.write(v, c.A&c.X)
	case DCP:
		// decrementMemory excluding tick
		m := c.read(v) - 1
		c.P.setZN(m)
		c.write(v, m)
		c.cmp(c.A, v)
	case ISB:
		// incrementMemory excluding tick
		m := c.read(v) + 1
		c.P.setZN(m)
		c.write(v, m)
		c.sbc(v)
	case SLO:
		// arithmeticShiftLeft excluding tick
		m := c.read(v)
		c.P[status_C] = m&0x80 == 0x80
		m <<= 1
		c.write(v, m)
		c.ora(v)
	case RLA:
		// rotateLeft excluding tick
		m := c.read(v)
		carry := m & 0x80
		m <<= 1
		if c.P[status_C] {
			m |= 0x01
		}
		c.P[status_C] = carry == 0x80
		c.P.setZN(m)
		c.write(v, m)
		c.and(v)
	case SRE:
		// logicalShiftRight excluding tick
		m := c.read(v)
		c.P[status_C] = m&1 == 1
		m >>= 1
		c.P.setZN(m)
		c.write(v, m)
		c.eor(v)
	case RRA:
		// rotateRight excluding tick
		m := c.read(v)
		carry := m & 1
		m >>= 1
		if c.P[status_C] {
			m |= 0x80
		}
		c.P[status_C] = carry == 1
		c.P.setZN(m)
		c.write(v, m)
		c.adc(v)
	default:
		panic(fmt.Sprintf("unrecognized mnemonic: %d", inst.Mnemonic))
	}
}

func (c *CPU) and(v uint16) {
	c.A &= c.read(v)
	c.P.setZN(c.A)
}

func (c *CPU) eor(v uint16) {
	c.A ^= c.read(v)
	c.P.setZN(c.A)
}

func (c *CPU) ora(v uint16) {
	c.A |= c.read(v)
	c.P.setZN(c.A)
}

func (c *CPU) carry(m, r uint8) {
	a7 := c.A >> 7 & 1
	m7 := m >> 7 & 1
	c6 := a7 ^ m7 ^ (r >> 7 & 1)
	c7 := (a7 & m7) | (a7 & c6) | (m7 & c6)
	c.P[status_C] = c7 == 1
	c.P[status_V] = c6^c7 == 1
}

func (c *CPU) adc(v uint16) {
	m := c.read(v)
	r := c.A + m
	if c.P[status_C] {
		r += 1
	}
	c.carry(m, r)
	c.A = r
	c.P.setZN(c.A)
}

func (c *CPU) sbc(v uint16) {
	m := ^c.read(v)
	r := c.A + m
	if c.P[status_C] {
		r += 1
	}
	c.carry(m, r)
	c.A = r
	c.P.setZN(c.A)
}

func (c *CPU) cmp(x uint8, v uint16) {
	r := int16(x) - int16(c.read(v))
	c.P.setZN(uint8(r))
	c.P[status_C] = 0 <= r
}

func (c *CPU) branch(v uint16, cond bool) {
	if !cond {
		return
	}
	c.tick()
	base := int16(c.PC)
	offset := int8(v) // to negative number
	if pageCrossed(int16(offset), base) {
		c.tick()
	}
	c.PC = uint16(base + int16(offset))
}

func (c *CPU) pushStack(v uint8) {
	c.write(uint16(c.S)+0x0100, v)
	c.S -= 1
}

func (c *CPU) pushStackWord(v uint16) {
	c.pushStack(uint8(v >> 8))
	c.pushStack(uint8(v & 0xFF))
}

func (c *CPU) pullStack() uint8 {
	c.S += 1
	return c.read(uint16(c.S) + 0x0100)
}

func (c *CPU) pullStackWord() uint16 {
	return uint16(c.pullStack()) | uint16(c.pullStack())<<8
}
