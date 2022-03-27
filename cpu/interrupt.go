package cpu

// interrupt Kinds of CPU interrupts
type interrupt uint8

// currently supports NMI and IRQ only
const (
	_ interrupt = iota
	NMI
	IRQ
)
