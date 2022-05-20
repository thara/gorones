package apu

import (
	"context"
	"testing"
	"time"
)

func Test_sweep_oneComplement(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	t.Run("not enabled", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.input(0b01010101)

		s.clock()
		assertRecvValue(t, s.output(), 0)
	})

	t.Run("enabled but muted", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.update(0b10000000) // enabled

		s.input(0b00000000)

		s.clock()
		assertRecvValue(t, s.output(), 0)
	})

	t.Run("update P", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.update(0b10110000) // enabled, P = 3

		s.input(0b00000101)

		// P + 1
		s.clock()
		s.clock()
		s.clock()
		s.clock()
		assertRecvValue(t, s.output(), 1)
	})

	t.Run("negate without shift", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.update(0b10001000) // enabled, negate, no shift

		s.input(0b00000101)

		s.clock()
		assertRecvValue(t, s.output(), 0)
	})

	t.Run("negate with shift", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.update(0b10001010) // enabled, negate, 2 shift

		s.input(0b00000101)

		s.clock()
		assertRecvValue(t, s.output(), 1)
	})
}

func Test_sweep_twoComplement(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	t.Run("negate without shift", func(t *testing.T) {
		s := runSweep(ctx, sweepTwoComplement)
		s.update(0b10001000) // enabled, negate, no shift

		s.input(0b00000101)

		s.clock()
		assertRecvValue(t, s.output(), 1)
	})

	t.Run("negate with shift", func(t *testing.T) {
		s := runSweep(ctx, sweepOneComplement)
		s.update(0b10001010) // enabled, negate, 2 shift

		s.input(0b00000101)

		s.clock()
		assertRecvValue(t, s.output(), 1)
	})
}
