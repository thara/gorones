package apu

import (
	"context"

	"github.com/thara/gorones/util"
)

// https://www.nesdev.org/wiki/APU_Envelope

type envelope struct {
	s   chan interface{}
	in  chan interface{}
	upd chan byte
	out chan uint8
}

func runEnvelope(ctx context.Context) *envelope {
	s := make(chan interface{})
	in := make(chan interface{})
	upd := make(chan byte)
	out := make(chan uint8)
	go func() {
		defer close(s)
		defer close(in)
		defer close(upd)
		defer close(out)

		var (
			start      bool
			loopFlag   bool
			constVol   bool
			decayLevel uint8
		)
		var v uint
		d := runDivider(ctx, v)
		for {
			select {
			case <-ctx.Done():
				return

			case <-s:
				start = true

			case b := <-upd:
				loopFlag = util.IsSet(b, 5)
				constVol = util.IsSet(b, 4)
				v = uint(b & 0b1111)
				d.reload(v)

			case <-in:
				if start {
					start = false
					decayLevel = 15
					d.reload(v)
				} else {
					d.clock()
				}

			case <-d.output():
				d.reload(v)

				if 0 < decayLevel {
					decayLevel--
				} else if loopFlag {
					decayLevel = 15
				}

				//TODO I'm not sure about envelope output timing...
				if constVol {
					out <- uint8(v)
				} else {
					out <- decayLevel
				}
			}
		}
	}()
	return &envelope{s, in, upd, out}
}

func (e *envelope) start() {
	e.s <- struct{}{}
}

func (e *envelope) clock() {
	e.in <- struct{}{}
}

func (e *envelope) update(v byte) {
	e.upd <- v
}

func (e *envelope) output() <-chan uint8 { return e.out }
