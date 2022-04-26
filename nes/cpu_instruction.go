package nes

import (
	"fmt"
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

func (n *NES) getOperand(m addressingMode) uint16 {
	switch m {
	case implicit:
		return 0
	case accumulator:
		return uint16(n.cpu.A)
	case immediate:
		pc := n.cpu.PC
		n.cpu.PC++
		return pc
	case zeroPage:
		v := n.read(n.cpu.PC)
		n.cpu.PC++
		return uint16(v)
	case zeroPageX:
		n.tick()
		v := (uint16(n.read(n.cpu.PC)) + uint16(n.cpu.X)) & 0xFF
		n.cpu.PC++
		return uint16(v)
	case zeroPageY:
		n.tick()
		v := (uint16(n.read(n.cpu.PC)) + uint16(n.cpu.Y)) & 0xFF
		n.cpu.PC++
		return uint16(v)
	case absolute:
		v := n.readWord(n.cpu.PC)
		n.cpu.PC += 2
		return v
	case absoluteX:
		v := n.readWord(n.cpu.PC)
		n.cpu.PC += 2
		n.tick()
		return v + uint16(n.cpu.X)
	case absoluteXWithPenalty:
		v := n.readWord(n.cpu.PC)
		n.cpu.PC += 2
		if pageCrossed(uint16(n.cpu.X), v) {
			n.tick()
		}
		return v + uint16(n.cpu.X)
	case absoluteY:
		v := n.readWord(n.cpu.PC)
		n.cpu.PC += 2
		n.tick()
		return v + uint16(n.cpu.Y)
	case absoluteYWithPenalty:
		v := n.readWord(n.cpu.PC)
		n.cpu.PC += 2
		if pageCrossed(uint16(n.cpu.Y), v) {
			n.tick()
		}
		return v + uint16(n.cpu.Y)
	case relative:
		v := n.read(n.cpu.PC)
		n.cpu.PC++
		return uint16(v)
	case indirect:
		m := n.readWord(n.cpu.PC)
		v := n.readOnIndirect(m)
		n.cpu.PC += 2
		return v
	case indexedIndirect:
		m := n.read(n.cpu.PC)
		v := n.readOnIndirect(uint16(m + n.cpu.X))
		n.cpu.PC += 1
		n.tick()
		return v
	case indirectIndexed:
		m := n.read(n.cpu.PC)
		v := n.readOnIndirect(uint16(m))
		n.cpu.PC += 1
		n.tick()
		return v + uint16(n.cpu.Y)
	case indirectIndexedWithPenalty:
		m := n.read(n.cpu.PC)
		v := n.readOnIndirect(uint16(m))
		n.cpu.PC += 1
		if pageCrossed(uint16(n.cpu.Y), v) {
			n.tick()
		}
		return v + uint16(n.cpu.Y)
	}

	panic("unrecognized addressing mode")
}

func (n *NES) execute(inst instruction) {
	v := n.getOperand(inst.AddressingMode)

	switch inst.Mnemonic {
	case LDA:
		n.cpu.A = n.read(v)
		n.cpu.P.setZN(n.cpu.A)
	case LDX:
		n.cpu.X = n.read(v)
		n.cpu.P.setZN(n.cpu.X)
	case LDY:
		n.cpu.Y = n.read(v)
		n.cpu.P.setZN(n.cpu.Y)
	case STA:
		n.write(v, n.cpu.A)
	case STX:
		n.write(v, n.cpu.X)
	case STY:
		n.write(v, n.cpu.Y)
	case TAX:
		n.cpu.X = n.cpu.A
		n.cpu.P.setZN(n.cpu.X)
		n.tick()
	case TAY:
		n.cpu.Y = n.cpu.A
		n.cpu.P.setZN(n.cpu.Y)
		n.tick()
	case TXA:
		n.cpu.A = n.cpu.X
		n.cpu.P.setZN(n.cpu.A)
		n.tick()
	case TYA:
		n.cpu.A = n.cpu.Y
		n.cpu.P.setZN(n.cpu.A)
		n.tick()
	case TSX:
		n.cpu.X = n.cpu.S
		n.cpu.P.setZN(n.cpu.X)
		n.tick()
	case TXS:
		n.cpu.S = n.cpu.X
		n.tick()

	case PHA:
		n.pushStack(n.cpu.A)
		n.tick()
	case PHP:
		p := n.cpu.P.u8() | instructionB
		n.pushStack(p)
		n.tick()
	case PLA:
		n.cpu.A = n.pullStack()
		n.cpu.P.setZN(n.cpu.A)
		n.tick_n(2)
	case PLP:
		v := n.pullStack() & ^instructionB
		v |= 0b100000 // for nestest
		n.cpu.P.set(v)
		n.tick_n(2)

	case AND:
		n.and(v)
	case EOR:
		n.eor(v)
	case ORA:
		n.ora(v)
	case BIT:
		m := n.read(v)
		b := n.cpu.A & m
		n.cpu.P[status_Z] = b == 0
		n.cpu.P[status_V] = m&0x40 == 0x40
		n.cpu.P[status_N] = m&0x80 == 0x80

	case ADC:
		m := n.read(v)
		r := n.cpu.A + m
		if n.cpu.P[status_C] {
			r += 1
		}
		n.carry(m, r)
		n.cpu.A = r
		n.cpu.P.setZN(n.cpu.A)
	case SBC:
		n.sbc(v)
	case CMP:
		n.cmp(n.cpu.A, v)
	case CPX:
		n.cmp(n.cpu.X, v)
	case CPY:
		n.cmp(n.cpu.Y, v)

	case INC:
		m := n.read(v)
		r := m + 1
		n.write(v, r)
		n.cpu.P.setZN(r)
		n.tick()
	case INX:
		n.cpu.X += 1
		n.cpu.P.setZN(n.cpu.X)
		n.tick()
	case INY:
		n.cpu.Y += 1
		n.cpu.P.setZN(n.cpu.Y)
		n.tick()
	case DEC:
		m := n.read(v)
		r := m - 1
		n.write(v, r)
		n.cpu.P.setZN(r)
		n.tick()
	case DEX:
		n.cpu.X -= 1
		n.cpu.P.setZN(n.cpu.X)
		n.tick()
	case DEY:
		n.cpu.Y -= 1
		n.cpu.P.setZN(n.cpu.Y)
		n.tick()

	case ASL:
		asl := func(m *uint8) {
			n.cpu.P[status_C] = *m&0x80 == 0x80
			*m <<= 1
			n.cpu.P.setZN(*m)
			n.tick()
		}
		if inst.AddressingMode == accumulator {
			asl(&n.cpu.A)
			return
		}
		m := n.read(v)
		asl(&m)
		n.write(v, m)

	case LSR:
		lsr := func(m *uint8) {
			n.cpu.P[status_C] = *m&1 == 1
			*m >>= 1
			n.cpu.P.setZN(*m)
			n.tick()
		}
		if inst.AddressingMode == accumulator {
			lsr(&n.cpu.A)
			return
		}
		m := n.read(v)
		lsr(&m)
		n.write(v, m)
	case ROL:
		rol := func(m *uint8) {
			carry := *m & 0x80
			*m <<= 1
			if n.cpu.P[status_C] {
				*m |= 1
			}
			n.cpu.P[status_C] = carry == 0x80
			n.cpu.P.setZN(*m)
			n.tick()
		}
		if inst.AddressingMode == accumulator {
			rol(&n.cpu.A)
			return
		}
		m := n.read(v)
		rol(&m)
		n.write(v, m)
	case ROR:
		ror := func(m *uint8) {
			carry := *m & 1
			*m >>= 1
			if n.cpu.P[status_C] {
				*m |= 0x80
			}
			n.cpu.P[status_C] = carry == 1
			n.cpu.P.setZN(*m)
			n.tick()
		}
		if inst.AddressingMode == accumulator {
			ror(&n.cpu.A)
			return
		}
		m := n.read(v)
		ror(&m)
		n.write(v, m)

	case JMP:
		n.cpu.PC = v
	case JSR:
		rtn := n.cpu.PC - 1
		n.pushStackWord(rtn)
		n.cpu.PC = v
		n.tick()
	case RTS:
		n.cpu.PC = n.pullStackWord()
		n.cpu.PC += 1
		n.tick_n(3)

	case BCC:
		n.branch(v, !n.cpu.P[status_C])
	case BCS:
		n.branch(v, n.cpu.P[status_C])
	case BEQ:
		n.branch(v, n.cpu.P[status_Z])
	case BMI:
		n.branch(v, n.cpu.P[status_N])
	case BNE:
		n.branch(v, !n.cpu.P[status_Z])
	case BPL:
		n.branch(v, !n.cpu.P[status_N])
	case BVC:
		n.branch(v, !n.cpu.P[status_V])
	case BVS:
		n.branch(v, n.cpu.P[status_V])

	case CLC:
		n.cpu.P[status_C] = false
		n.tick()
	case CLD:
		n.cpu.P[status_D] = false
		n.tick()
	case CLI:
		n.cpu.P[status_I] = false
		n.tick()
	case CLV:
		n.cpu.P[status_V] = false
		n.tick()
	case SEC:
		n.cpu.P[status_C] = true
		n.tick()
	case SED:
		n.cpu.P[status_D] = true
		n.tick()
	case SEI:
		n.cpu.P[status_I] = true
		n.tick()

	case BRK:
		n.pushStackWord(n.cpu.PC)
		n.cpu.P.insert(instructionB)
		n.pushStack(n.cpu.P.u8())
		n.cpu.PC = n.readWord(0xFFFE)
		n.tick()
	case NOP:
		n.tick()
	case RTI:
		v := n.pullStack()
		n.cpu.P.set(v)
		n.cpu.PC = n.pullStackWord()
		n.tick_n(2)

	case LAX:
		m := n.read(v)
		n.cpu.A = m
		n.cpu.P.setZN(m)
		n.cpu.X = m
	case SAX:
		n.write(v, n.cpu.A&n.cpu.X)
	case DCP:
		// decrementMemory excluding tick
		m := n.read(v) - 1
		n.cpu.P.setZN(m)
		n.write(v, m)
		n.cmp(n.cpu.A, v)
	case ISB:
		// incrementMemory excluding tick
		m := n.read(v) + 1
		n.cpu.P.setZN(m)
		n.write(v, m)
		n.sbc(v)
	case SLO:
		// arithmeticShiftLeft excluding tick
		m := n.read(v)
		n.cpu.P[status_C] = m&0x80 == 0x80
		m <<= 1
		n.write(v, m)
		n.ora(v)
	case RLA:
		// rotateLeft excluding tick
		m := n.read(v)
		carry := m & 0x80
		m <<= 1
		if n.cpu.P[status_C] {
			m |= 0x01
		}
		n.cpu.P[status_C] = carry == 0x80
		n.cpu.P.setZN(m)
		n.write(v, m)
		n.and(v)
	case SRE:
		// logicalShiftRight excluding tick
		m := n.read(v)
		n.cpu.P[status_C] = m&1 == 1
		m >>= 1
		n.cpu.P.setZN(m)
		n.write(v, m)
		n.eor(v)
	case RRA:
		// rotateRight excluding tick
		m := n.read(v)
		carry := m & 1
		m >>= 1
		if n.cpu.P[status_C] {
			m |= 0x80
		}
		n.cpu.P[status_C] = carry == 1
		n.cpu.P.setZN(m)
		n.write(v, m)
		n.adc(v)
	default:
		panic(fmt.Sprintf("unrecognized mnemonic: %d", inst.Mnemonic))
	}
}

func (n *NES) and(v uint16) {
	n.cpu.A &= n.read(v)
	n.cpu.P.setZN(n.cpu.A)
}

func (n *NES) eor(v uint16) {
	n.cpu.A ^= n.read(v)
	n.cpu.P.setZN(n.cpu.A)
}

func (n *NES) ora(v uint16) {
	n.cpu.A |= n.read(v)
	n.cpu.P.setZN(n.cpu.A)
}

func (n *NES) carry(m, r uint8) {
	a7 := n.cpu.A >> 7 & 1
	m7 := m >> 7 & 1
	c6 := a7 ^ m7 ^ (r >> 7 & 1)
	c7 := (a7 & m7) | (a7 & c6) | (m7 & c6)
	n.cpu.P[status_C] = c7 == 1
	n.cpu.P[status_V] = c6^c7 == 1
}

func (n *NES) adc(v uint16) {
	m := n.read(v)
	r := n.cpu.A + m
	if n.cpu.P[status_C] {
		r += 1
	}
	n.carry(m, r)
	n.cpu.A = r
	n.cpu.P.setZN(n.cpu.A)
}

func (n *NES) sbc(v uint16) {
	m := ^n.read(v)
	r := n.cpu.A + m
	if n.cpu.P[status_C] {
		r += 1
	}
	n.carry(m, r)
	n.cpu.A = r
	n.cpu.P.setZN(n.cpu.A)
}

func (n *NES) cmp(x uint8, v uint16) {
	r := int16(x) - int16(n.read(v))
	n.cpu.P.setZN(uint8(r))
	n.cpu.P[status_C] = 0 <= r
}

func (n *NES) branch(v uint16, cond bool) {
	if !cond {
		return
	}
	n.tick()
	base := int16(n.cpu.PC)
	offset := int8(v) // to negative number
	if pageCrossed(int16(offset), base) {
		n.tick()
	}
	n.cpu.PC = uint16(base + int16(offset))
}
