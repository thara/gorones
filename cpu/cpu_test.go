package cpu

type busMock []uint8

func newBusMock() busMock {
	return make(busMock, 0x10000)
}

func (m busMock) ReadCPU(addr uint16) uint8 {
	return m[addr]
}

func (m busMock) WriteCPU(addr uint16, value uint8) {
	m[addr] = value
}

type tickFn func()

func (f tickFn) Tick() {
	f()
}

var tickMock = tickFn(func() {})
