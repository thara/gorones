package cpu

import "fmt"

/// cpu has cpu registers and clock cycle
type cpu struct {
	// https://wiki.nesdev.org/w/index.php?title=CPU_registers

	// Accumulator, Index X/Y register
	A, X, Y uint8
	// Stack pointer
	S uint8
	// Status register
	P Status
	// Program counter
	PC uint16

	Cycles uint64
}

func (c *cpu) reset() {
	*c = cpu{}
}

// https://www.nesdev.org/wiki/Status_flags

// Status represents CPU Status by flags
type Status [6]bool

func NewStatus(b uint8) Status {
	var s Status
	s.set(b)
	return s
}

func (s *Status) u8() uint8 {
	return bit(s[status_C]) |
		bit(s[status_Z])<<1 |
		bit(s[status_I])<<2 |
		bit(s[status_D])<<3 |
		bit(s[status_V])<<6 |
		bit(s[status_N])<<7
}

func (s *Status) set(b uint8) {
	s[status_C] = b&1 == 1
	s[status_Z] = (b>>1)&1 == 1
	s[status_I] = (b>>2)&1 == 1
	s[status_D] = (b>>3)&1 == 1
	s[status_V] = (b>>6)&1 == 1
	s[status_N] = (b>>7)&1 == 1
}

func (s *Status) insert(b uint8) {
	s.set(s.u8() | b)
}

func (s *Status) setZN(v uint8) {
	s[status_Z] = v == 0
	s[status_N] = v&0x80 == 0x80
}

// statusFlag is index of `status` struct
type statusFlag uint8

const (
	status_C statusFlag = iota // Carry
	status_Z                   // Zero
	status_I                   // Interrupt Disable
	status_D                   // Decimal
	status_V                   // Overflow
	status_N                   // Negative

	// B flags
	// https://wiki.nesdev.org/w/index.php?title=Status_flags#The_B_flag
	interruptB   uint8 = 0b00100000
	instructionB uint8 = 0b00110000
)

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
