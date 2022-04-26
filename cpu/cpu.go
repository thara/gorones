package cpu

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

// https://www.nesdev.org/wiki/Status_flags

// Status represents CPU Status by flags
type Status [6]bool

func NewStatus(b uint8) Status {
	var s Status
	s.Set(b)
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

func (s *Status) Set(b uint8) {
	s[status_C] = b&1 == 1
	s[status_Z] = (b>>1)&1 == 1
	s[status_I] = (b>>2)&1 == 1
	s[status_D] = (b>>3)&1 == 1
	s[status_V] = (b>>6)&1 == 1
	s[status_N] = (b>>7)&1 == 1
}

func (s *Status) insert(b uint8) {
	s.Set(s.u8() | b)
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
