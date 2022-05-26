package apu

import "context"

type sequencer struct {
	init chan interface{}
	clk  chan interface{}
	out  chan uint8
}

func runSequencer(ctx context.Context) *sequencer {
	init := make(chan interface{})
	clock := make(chan interface{})
	out := make(chan uint8, 1)
	go func() {
		defer close(init)
		defer close(clock)
		defer close(out)

		var seq uint8
		for {
			select {
			case <-ctx.Done():
				return
			case <-init:
				seq = 0
			case <-clock:
				seq += 1
				if seq == 8 {
					seq = 0
				}
				out <- seq
			}
		}
	}()
	return &sequencer{init, clock, out}
}

func (s *sequencer) clock()               { s.clk <- struct{}{} }
func (s *sequencer) restart()             { s.init <- struct{}{} }
func (s *sequencer) output() <-chan uint8 { return s.out }
