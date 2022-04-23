package cpu

func Decode(opcode uint8) instruction {
	return instructions[opcode]
}

var (
	instructions [256]instruction = func() [256]instruction {
		var s []instruction
		for i := 0; i < 256; i++ {
			s = append(s, newInstruction(uint8(i)))
		}
		var a [256]instruction
		copy(a[:], s[:256])
		return a
	}()
)

func newInstruction(op uint8) instruction {
	switch op {
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
