package mapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseROM(t *testing.T) {
	f, err := os.Open("../testdata/nestest.nes")
	require.NoError(t, err)
	defer f.Close()

	rom, err := ParseROM(f)
	require.NoError(t, err)

	assert.EqualValues(t, 0, rom.header.mapperNO)
	assert.EqualValues(t, 1, rom.header.prgROMSize)
	assert.EqualValues(t, 1, rom.header.chrROMSize)
	assert.EqualValues(t, Mirroring_Horizontal, rom.header.mirroring)
}
