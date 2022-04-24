package cpu

// Emu emulates CPU behaviors
type Emu struct {
	cpu CPU

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
func NewEmu(t Ticker, m Bus) *Emu {
	return &Emu{t: t, m: m}
}

// PowerOn initializes CPU state on power
func (e *Emu) PowerOn() {
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state

	// IRQ disabled
	e.cpu.P.set(0x34)
	e.cpu.A = 0x00
	e.cpu.X = 0x00
	e.cpu.Y = 0x00
	e.cpu.S = 0xFD
	// frame irq disabled
	e.write(0x4017, 0x00)
	// all channels disabled
	e.write(0x4015, 0x00)
}

// Step emulates 1 CPU step.
func (e *Emu) Step() {
	e.handleInterrupt()

	op := e.fetch()
	inst := Decode(op)
	e.execute(inst)
}

func (e *Emu) fetch() uint8 {
	op := e.read(e.cpu.PC)
	e.cpu.PC += 1
	return op
}

func (e *Emu) tick() {
	e.t.Tick()
	e.cpu.Cycles += 1
}

func (e *Emu) tick_n(n uint) {
	for i := uint(0); i < n; i++ {
		e.t.Tick()
	}
	e.cpu.Cycles += uint64(n)
}

func (e *Emu) read(addr uint16) uint8 {
	v := e.m.ReadCPU(addr)
	e.tick()
	return v
}

func (e *Emu) readWord(addr uint16) uint16 {
	return uint16(e.read(addr)) | uint16(e.read(addr+1))<<8
}

func (e *Emu) write(addr uint16, value uint8) {
	e.m.WriteCPU(addr, value)
	e.tick()
}

func (e *Emu) readOnIndirect(addr uint16) uint16 {
	low := uint16(e.read(addr))
	// Reproduce 6502 bug - http://nesdev.com/6502bugs.txt
	high := uint16(e.read((addr & 0xFF00) | ((addr + 1) & 0x00FF)))
	return low | (high << 8)
}

func (e *Emu) pushStack(v uint8) {
	e.write(uint16(e.cpu.S)+0x0100, v)
	e.cpu.S -= 1
}

func (e *Emu) pushStackWord(v uint16) {
	e.pushStack(uint8(v >> 8))
	e.pushStack(uint8(v & 0xFF))
}

func (e *Emu) pullStack() uint8 {
	e.cpu.S += 1
	return e.read(uint16(e.cpu.S) + 0x0100)
}

func (e *Emu) pullStackWord() uint16 {
	return uint16(e.pullStack()) | uint16(e.pullStack())<<8
}

func (e *Emu) handleInterrupt() {
	if e.interrupt == nil {
		return
	}

	intr := *e.interrupt
	switch intr {
	case NMI: // interrupted
	case IRQ:
		if !e.cpu.P[status_I] {
			return
		}
		// interrupted
	default:
		return
	}

	e.tick_n(2)
	e.pushStackWord(e.cpu.PC)
	// https://wiki.nesdev.com/w/index.php/Status_flags#The_B_flag
	// http://visual6502.org/wiki/index.php?title=6502_BRK_and_B_bit
	e.pushStack(e.cpu.P.u8() | interruptB)
	e.cpu.P[status_I] = true
	e.cpu.PC = e.readWord(intr.vector())
	e.interrupt = nil
}

func (e *Emu) InitNESTest() {
	e.cpu.PC = 0xC000
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_ref-1
	e.cpu.P.set(0x24)
	e.cpu.Cycles = 7
}
