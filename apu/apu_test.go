package apu

import (
	"context"
	"testing"
	"time"
)

func TestPulse(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	p := runPulse(ctx, sweepOneComplement)
	p.clock()
	assertRecvValue(t, p.output(), 0)
}
