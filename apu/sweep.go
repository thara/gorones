package apu

import (
	"context"

	"github.com/thara/gorones/util"
)

// https://www.nesdev.org/wiki/APU_Sweep

type sweep struct {
	p   chan uint16
	clk chan interface{}
	upd chan byte
	out chan uint8
}

func runSweep(ctx context.Context, cmp sweepComplement) *sweep {
	in := make(chan uint16)
	clk := make(chan interface{})
	upd := make(chan byte)
	out := make(chan uint8)

	go func() {
		defer close(in)
		defer close(clk)
		defer close(upd)
		defer close(out)

		enabled := false
		var p uint8
		negate := false
		var shift uint8

		muted := false

		var targetPeriod uint16

		d := runDivider(ctx, p)

		for {
			select {
			case <-ctx.Done():
				return
			case b := <-upd:
				enabled = util.IsSet(b, 7)
				p = b & 0b01110000 >> 4
				negate = util.IsSet(b, 3)
				shift = b & 0b111
				d.reload(p)

			case <-clk:
				d.clock()

			case currentPeriod := <-in:
				c := currentPeriod >> shift
				if negate {
					c = ^c + uint16(cmp)
				}
				targetPeriod = currentPeriod + uint16(c)
				muted = currentPeriod == 0 || 0x7FF < targetPeriod

			case <-d.output():
				if enabled && !muted {
					out <- 1
				} else {
					out <- 0
				}
			}
		}
	}()
	return &sweep{in, clk, upd, out}
}

func (s *sweep) input(n uint16)       { s.p <- n }
func (s *sweep) clock()               { s.clk <- struct{}{} }
func (s *sweep) update(v byte)        { s.upd <- v }
func (s *sweep) output() <-chan uint8 { return s.out }

type sweepComplement uint16

const (
	sweepOneComplement sweepComplement = 0
	sweepTwoComplement                 = 1
)
