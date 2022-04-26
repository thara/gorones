package cpu

// CPU emulates CPU behaviors
type CPU struct {
	cpu

	t Ticker
	m Bus

	interrupt *interrupt
}

// Ticker is the event handler on CPU ticks.
type Ticker interface {
	// Tick handles events on CPU ticks.
	Tick()
}

// Bus is the abstraction of memory access seen from CPU emulator.
//
// This interface provides a strategy of memory access to target.
type Bus interface {
	// ReadCPU read a byte by CPU.
	//
	// A location pointed of address may be one of NES component or other RAM/ROM.
	ReadCPU(addr uint16) uint8

	// WriteCPU write a byte from any component by CPU.
	//
	// A location pointed of address may be one of NES component or other RAM/ROM.
	WriteCPU(addr uint16, value uint8)
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
func (c *CPU) Step() {
	c.handleInterrupt()

	op := c.fetch()
	inst := Decode(op)
	c.execute(inst)
}

func (c *CPU) fetch() uint8 {
	op := c.read(c.PC)
	c.PC += 1
	return op
}

func (c *CPU) tick() {
	c.t.Tick()
	c.Cycles += 1
}

func (c *CPU) tick_n(n uint) {
	for i := uint(0); i < n; i++ {
		c.t.Tick()
	}
	c.Cycles += uint64(n)
}

func (c *CPU) read(addr uint16) uint8 {
	v := c.m.ReadCPU(addr)
	c.tick()
	return v
}

func (c *CPU) readWord(addr uint16) uint16 {
	return uint16(c.read(addr)) | uint16(c.read(addr+1))<<8
}

func (c *CPU) write(addr uint16, value uint8) {
	c.m.WriteCPU(addr, value)
	c.tick()
}

func (c *CPU) readOnIndirect(addr uint16) uint16 {
	low := uint16(c.read(addr))
	// Reproduce 6502 bug - http://nesdev.com/6502bugs.txt
	high := uint16(c.read((addr & 0xFF00) | ((addr + 1) & 0x00FF)))
	return low | (high << 8)
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
