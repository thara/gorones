package gorones

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thara/gorones/cpu"
	"github.com/thara/gorones/input"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/ppu"
)

type nopFrameRenderer struct{}

func (nopFrameRenderer) UpdateFrame(*[ppu.WIDTH * ppu.HEIGHT]uint8) {}

func Test_nestest(t *testing.T) {
	f, err := os.Open("testdata/nestest.nes")
	require.NoError(t, err)
	defer f.Close()

	rom, err := mapper.ParseROM(f)
	require.NoError(t, err)
	m, err := rom.Mapper()
	require.NoError(t, err)

	var ctrl1, ctrl2 input.StandardController

	intr := cpu.NoInterrupt

	ppu := ppu.New(m, new(nopFrameRenderer))
	ticker := cpuTicker{ppu: ppu, interrupt: &intr}
	bus := cpuBus{mapper: m, ctrl1: &ctrl1, ctrl2: &ctrl2, t: &ticker}
	nes := &NES{
		cpu:       cpu.New(&ticker, &bus),
		interrupt: &intr,
	}
	nes.PowerOn()

	nes.cpu.PC = 0xC000
	// https://wiki.nesdev.com/w/index.php/CPU_power_up_state#cite_ref-1
	nes.cpu.P.Set(0x24)
	nes.cpu.Cycles = 7

	f, err = os.Open("testdata/nestest.log")
	require.NoError(t, err)
	defer f.Close()

	i := 1

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()

		tr := nes.cpu.Trace()
		nes.step()

		li, err := parseCPUTrace(line)
		require.NoError(t, err)

		require.Equal(t, *li, tr, "lineno:%d %s", i, line)

		i++
	}
	require.NoError(t, sc.Err())

	assert.EqualValues(t, 26560, nes.cpu.Trace().Cycles)

	assert.EqualValues(t, 0, bus.ReadCPU(0x0002))
	assert.EqualValues(t, 0, bus.ReadCPU(0x0003))
}

// nestest.log line example:
// C000  4C F5 C5  JMP $C5F5                       A:00 X:00 Y:00 P:24 SP:FD PPU:  0, 21 CYC:7

var nestestLogLineRe = regexp.MustCompile(
	`(?P<pc>.{4})  (?P<op>.{2}) (?P<op1>.{2}) (?P<op2>.{2}) (\s|\*)(?P<code>.{32})A:(?P<a>.{2}) X:(?P<x>.{2}) Y:(?P<y>.{2}) P:(?P<p>.{2}) SP:(?P<s>.{2}) PPU:(?P<dot>.{3}),(?P<line>.{3}) CYC:(?P<cyc>.*)`)

func parseCPUTrace(line string) (*cpu.Trace, error) {
	re := nestestLogLineRe

	var t cpu.Trace

	matches := re.FindStringSubmatch(line)
	toInt := func(name string, base int, out *int64) error {
		i := re.SubexpIndex(name)
		m := matches[i]
		if len(strings.Trim(m, " ")) == 0 {
			*out = 0
			return nil
		}
		n, err := strconv.ParseInt(matches[i], base, 64)
		if err != nil {
			return errors.WithStack(err)
		}
		*out = n
		return nil
	}

	var n int64
	err := toInt("pc", 16, &n)
	if err != nil {
		return nil, err
	}
	t.PC = uint16(n)

	err = toInt("op", 16, &n)
	if err != nil {
		return nil, err
	}
	t.Opcode = uint8(n)

	err = toInt("op1", 16, &n)
	if err != nil {
		return nil, err
	}
	t.Operand1 = uint8(n)
	err = toInt("op2", 16, &n)
	if err != nil {
		return nil, err
	}
	t.Operand2 = uint8(n)

	err = toInt("a", 16, &n)
	if err != nil {
		return nil, err
	}
	t.A = uint8(n)
	err = toInt("x", 16, &n)
	if err != nil {
		return nil, err
	}
	t.X = uint8(n)
	err = toInt("y", 16, &n)
	if err != nil {
		return nil, err
	}
	t.Y = uint8(n)
	err = toInt("p", 16, &n)
	if err != nil {
		return nil, err
	}
	t.P = cpu.NewStatus(uint8(n))
	err = toInt("s", 16, &n)
	if err != nil {
		return nil, err
	}
	t.S = uint8(n)

	err = toInt("cyc", 10, &n)
	if err != nil {
		return nil, err
	}
	t.Cycles = uint64(n)

	inst := cpu.Decode(t.Opcode)
	t.Mnemonic = inst.Mnemonic
	t.AddressingMode = inst.AddressingMode

	return &t, nil
}
