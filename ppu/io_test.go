package ppu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thara/gorones/mapper"
)

func Test_toNTAddr(t *testing.T) {
	tests := []struct {
		mirroring mapper.Mirroring
		addr      uint16
		expected  uint16
	}{
		{mapper.Mirroring_Vertical, 0x2000, 0},
		{mapper.Mirroring_Vertical, 0x23FF, 0x03FF},
		{mapper.Mirroring_Vertical, 0x2400, 0x0400},
		{mapper.Mirroring_Vertical, 0x27FF, 0x07FF},
		{mapper.Mirroring_Vertical, 0x2800, 0},
		{mapper.Mirroring_Vertical, 0x2BFF, 0x03FF},
		{mapper.Mirroring_Vertical, 0x2C00, 0x0400},
		{mapper.Mirroring_Vertical, 0x2FFF, 0x07FF},
		{mapper.Mirroring_Horizontal, 0x2000, 0},
		{mapper.Mirroring_Horizontal, 0x23FF, 0x03FF},
		{mapper.Mirroring_Horizontal, 0x2400, 0},
		{mapper.Mirroring_Horizontal, 0x27FF, 0x03FF},
		{mapper.Mirroring_Horizontal, 0x2800, 0x0800},
		{mapper.Mirroring_Horizontal, 0x2BFF, 0x0BFF},
		{mapper.Mirroring_Horizontal, 0x2C00, 0x0800},
		{mapper.Mirroring_Horizontal, 0x2FFF, 0x0BFF},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			result := toNTAddr(tt.addr, tt.mirroring)
			assert.EqualValues(t, tt.expected, result)
		})
	}
}
