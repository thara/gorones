package nes

type NES struct {
	cpu CPU

	t Ticker
	m Bus

	interrupt *interrupt
}

func (n *NES) PowerOn() {
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state

	// IRQ disabled
	n.cpu.P.set(0x34)
	n.cpu.A = 0x00
	n.cpu.X = 0x00
	n.cpu.Y = 0x00
	n.cpu.S = 0xFD
	// frame irq disabled
	n.write(0x4017, 0x00)
	// all channels disabled
	n.write(0x4015, 0x00)
}

func (n *NES) step() {
	n.handleInterrupt()

	op := n.fetch()
	inst := Decode(op)
	n.execute(inst)
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

func (n *NES) tick() {
	n.t.Tick()
	n.cpu.Cycles += 1
}

func (n *NES) tick_n(N uint) {
	for i := uint(0); i < N; i++ {
		n.t.Tick()
	}
	n.cpu.Cycles += uint64(N)
}

type cpuTicker struct {
	cycles uint64
}

func (c *cpuTicker) Tick() {
	c.cycles += 1
}
