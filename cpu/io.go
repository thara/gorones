package cpu

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
