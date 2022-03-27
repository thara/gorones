package cpu

// https://www.nesdev.org/wiki/Status_flags

// status represents CPU status by flags
type status [6]bool

func (s *status) u8() uint8 {
	return bit(s[status_C]) |
		bit(s[status_Z])<<1 |
		bit(s[status_I])<<2 |
		bit(s[status_D])<<3 |
		bit(s[status_V])<<6 |
		bit(s[status_N])<<7
}

func (s *status) set(b uint8) {
	s[status_C] = b&1 == 1
	s[status_Z] = (b>>1)&1 == 1
	s[status_I] = (b>>2)&1 == 1
	s[status_D] = (b>>3)&1 == 1
	s[status_V] = (b>>6)&1 == 1
	s[status_N] = (b>>7)&1 == 1
}

func (s *status) insert(b uint8) {
	s.set(s.u8() | b)
}

func (s *status) setZN(v uint8) {
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
