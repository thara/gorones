package cpu

import (
	"fmt"
)

// https://www.nesdev.org/wiki/6502_instructions

type instruction struct {
	Mnemonic       mnemonic
	AddressingMode addressingMode
}

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

func (c *CPU) execute(inst instruction) {
	v := c.getOperand(inst.AddressingMode)

	switch inst.Mnemonic {
	case LDA:
		c.a = c.read(v)
		c.p.setZN(c.a)
	case LDX:
		c.x = c.read(v)
		c.p.setZN(c.x)
	case LDY:
		c.y = c.read(v)
		c.p.setZN(c.y)
	case STA:
		c.write(v, c.a)
	case STX:
		c.write(v, c.x)
	case STY:
		c.write(v, c.y)
	case TAX:
		c.x = c.a
		c.p.setZN(c.x)
		c.tick()
	case TAY:
		c.y = c.a
		c.p.setZN(c.y)
		c.tick()
	case TXA:
		c.a = c.x
		c.p.setZN(c.a)
		c.tick()
	case TYA:
		c.a = c.y
		c.p.setZN(c.a)
		c.tick()
	case TSX:
		c.x = c.s
		c.p.setZN(c.x)
		c.tick()
	case TXS:
		c.s = c.x
		c.tick()

	case PHA:
		c.pushStack(c.a)
		c.tick()
	case PHP:
		p := c.p.u8() | instructionB
		c.pushStack(p)
		c.tick()
	case PLA:
		c.a = c.pullStack()
		c.p.setZN(c.a)
		c.tick_n(2)
	case PLP:
		v := c.pullStack() & ^instructionB
		v |= 0b100000 // for nestest
		c.p.set(v)
		c.tick_n(2)

	case AND:
		c.and(v)
	case EOR:
		c.eor(v)
	case ORA:
		c.ora(v)
	case BIT:
		m := c.read(v)
		b := c.a & m
		c.p[status_Z] = b == 0
		c.p[status_V] = m&0x40 == 0x40
		c.p[status_N] = m&0x80 == 0x80

	case ADC:
		m := c.read(v)
		r := c.a + m
		if c.p[status_C] {
			r += 1
		}
		c.carry(m, r)
		c.a = r
		c.p.setZN(c.a)
	case SBC:
		c.sbc(v)
	case CMP:
		c.cmp(c.a, v)
	case CPX:
		c.cmp(c.x, v)
	case CPY:
		c.cmp(c.y, v)

	case INC:
		m := c.read(v)
		r := m + 1
		c.write(v, r)
		c.p.setZN(r)
		c.tick()
	case INX:
		c.x += 1
		c.p.setZN(c.x)
		c.tick()
	case INY:
		c.y += 1
		c.p.setZN(c.y)
		c.tick()
	case DEC:
		m := c.read(v)
		r := m - 1
		c.write(v, r)
		c.p.setZN(r)
		c.tick()
	case DEX:
		c.x -= 1
		c.p.setZN(c.x)
		c.tick()
	case DEY:
		c.y -= 1
		c.p.setZN(c.y)
		c.tick()

	case ASL:
		asl := func(m *uint8) {
			c.p[status_C] = *m&0x80 == 0x80
			*m <<= 1
			c.p.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			asl(&c.a)
			return
		}
		m := c.read(v)
		asl(&m)
		c.write(v, m)

	case LSR:
		lsr := func(m *uint8) {
			c.p[status_C] = *m&1 == 1
			*m >>= 1
			c.p.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			lsr(&c.a)
			return
		}
		m := c.read(v)
		lsr(&m)
		c.write(v, m)
	case ROL:
		rol := func(m *uint8) {
			carry := *m & 0x80
			*m <<= 1
			if c.p[status_C] {
				*m |= 1
			}
			c.p[status_C] = carry == 0x80
			c.p.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			rol(&c.a)
			return
		}
		m := c.read(v)
		rol(&m)
		c.write(v, m)
	case ROR:
		ror := func(m *uint8) {
			carry := *m & 1
			*m >>= 1
			if c.p[status_C] {
				*m |= 0x80
			}
			c.p[status_C] = carry == 1
			c.p.setZN(*m)
			c.tick()
		}
		if inst.AddressingMode == accumulator {
			ror(&c.a)
			return
		}
		m := c.read(v)
		ror(&m)
		c.write(v, m)

	case JMP:
		c.pc = v
	case JSR:
		rtn := c.pc - 1
		c.pushStackWord(rtn)
		c.pc = v
		c.tick()
	case RTS:
		c.pc = c.pullStackWord()
		c.pc += 1
		c.tick_n(3)

	case BCC:
		c.branch(v, !c.p[status_C])
	case BCS:
		c.branch(v, c.p[status_C])
	case BEQ:
		c.branch(v, c.p[status_Z])
	case BMI:
		c.branch(v, c.p[status_N])
	case BNE:
		c.branch(v, !c.p[status_Z])
	case BPL:
		c.branch(v, !c.p[status_N])
	case BVC:
		c.branch(v, !c.p[status_V])
	case BVS:
		c.branch(v, c.p[status_V])

	case CLC:
		c.p[status_C] = false
		c.tick()
	case CLD:
		c.p[status_D] = false
		c.tick()
	case CLI:
		c.p[status_I] = false
		c.tick()
	case CLV:
		c.p[status_V] = false
		c.tick()
	case SEC:
		c.p[status_C] = true
		c.tick()
	case SED:
		c.p[status_D] = true
		c.tick()
	case SEI:
		c.p[status_I] = true
		c.tick()

	case BRK:
		c.pushStackWord(c.pc)
		c.p.insert(instructionB)
		c.pushStack(c.p.u8())
		c.pc = c.readWord(0xFFFE)
		c.tick()
	case NOP:
		c.tick()
	case RTI:
		v := c.pullStack() & ^c.p.u8()
		v |= 0b10000 // for nestest
		c.p.set(v)
		c.pc = c.pullStackWord()
		c.tick_n(2)

	case LAX:
		m := c.read(v)
		c.a = m
		c.p.setZN(m)
		c.x = m
	case SAX:
		c.write(v, c.a&c.x)
	case DCP:
		// decrementMemory excluding tick
		m := c.read(v) - 1
		c.p.setZN(m)
		c.write(v, m)
		c.cmp(c.a, v)
	case ISB:
		// incrementMemory excluding tick
		m := c.read(v) + 1
		c.p.setZN(m)
		c.write(v, m)
		c.sbc(v)
	case SLO:
		// arithmeticShiftLeft excluding tick
		m := c.read(v)
		c.p[status_C] = m&0x80 == 0x80
		m <<= 1
		c.write(v, m)
		c.ora(v)
	case RLA:
		// rotateLeft excluding tick
		m := c.read(v)
		carry := m & 0x80
		m <<= 1
		if c.p[status_C] {
			m |= 0x01
		}
		c.p[status_C] = carry == 0x80
		c.p.setZN(m)
		c.write(v, m)
		c.and(v)
	case SRE:
		// logicalShiftRight excluding tick
		m := c.read(v)
		c.p[status_C] = m&1 == 1
		m >>= 1
		c.p.setZN(m)
		c.write(v, m)
		c.eor(v)
	case RRA:
		// rotateRight excluding tick
		m := c.read(v)
		carry := m & 1
		m >>= 1
		if c.p[status_C] {
			m |= 0x80
		}
		c.p[status_C] = carry == 1
		c.p.setZN(m)
		c.write(v, m)
		c.adc(v)
	default:
		panic(fmt.Sprintf("unrecognized mnemonic: %d", inst.Mnemonic))
	}
}

func (c *CPU) and(v uint16) {
	c.a &= c.read(v)
	c.p.setZN(c.a)
}

func (c *CPU) eor(v uint16) {
	c.a ^= c.read(v)
	c.p.setZN(c.a)
}

func (c *CPU) ora(v uint16) {
	c.a |= c.read(v)
	c.p.setZN(c.a)
}

func (c *CPU) carry(m, r uint8) {
	a7 := c.a >> 7 & 1
	m7 := m >> 7 & 1
	c6 := a7 ^ m7 ^ (r >> 7 & 1)
	c7 := (a7 & m7) | (a7 & c6) | (m7 & c6)
	c.p[status_C] = c7 == 1
	c.p[status_V] = c6^c7 == 1
}

func (c *CPU) adc(v uint16) {
	m := c.read(v)
	r := c.a + m
	if c.p[status_C] {
		r += 1
	}
	c.carry(m, r)
	c.a = r
	c.p.setZN(c.a)
}

func (c *CPU) sbc(v uint16) {
	m := ^c.read(v)
	r := c.a + m
	if c.p[status_C] {
		r += 1
	}
	c.carry(m, r)
	c.a = r
	c.p.setZN(c.a)
}

func (c *CPU) cmp(x uint8, v uint16) {
	r := int16(x) - int16(c.read(v))
	c.p.setZN(uint8(r))
	c.p[status_C] = 0 <= r
}

func (c *CPU) branch(v uint16, cond bool) {
	if !cond {
		return
	}
	c.tick()
	base := int16(c.pc)
	offset := int8(v) // to negative number
	if pageCrossed(int16(offset), base) {
		c.tick()
	}
	c.pc = uint16(base + int16(offset))
}
