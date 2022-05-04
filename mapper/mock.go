package mapper

// MapperMock for test
type MapperMock struct{}

func (*MapperMock) Read(addr uint16) uint8         { return 0 }
func (*MapperMock) Write(addr uint16, value uint8) {}
func (*MapperMock) Mirroring() Mirroring           { return Mirroring_Horizontal }
func (*MapperMock) PRG() []byte                    { return nil }
func (*MapperMock) CHR() []byte                    { return nil }
