package cpu

// CPU emulates CPU behaviors
type CPU struct {
	// https://wiki.nesdev.org/w/index.php?title=CPU_registers

	// Accumulator, Index X/Y register
	A, X, Y uint8
	// Stack pointer
	S uint8
	// Status register
	P Status
	// Program counter
	PC uint16

	// clock cycle
	Cycles uint64

	t Ticker
	m Bus
}

// New returns new CPU emulator
func New(t Ticker, m Bus) *CPU {
	return &CPU{t: t, m: m}
}

// PowerOn initializes CPU state on power
func (c *CPU) PowerOn() {
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state

	// IRQ disabled
	c.P.Set(0x34)
	c.A = 0x00
	c.X = 0x00
	c.Y = 0x00
	c.S = 0xFD
	// frame irq disabled
	c.write(0x4017, 0x00)
	// all channels disabled
	c.write(0x4015, 0x00)
}

// Step emulates 1 CPU step.
func (c *CPU) Step(intr *Interrupt) {
	c.handleInterrupt(intr)

	op := c.fetch()
	inst := Decode(op)
	c.execute(inst)
}

func (c *CPU) fetch() uint8 {
	op := c.read(c.PC)
	c.PC += 1
	return op
}

func (c *CPU) Reset() {
	c.PC = c.readWord(0xFFFC)
	c.P.Set(uint8(status_I))
	c.S -= 3
}
