package apu

import (
	"context"

	"github.com/thara/gorones/util"
)

// https://www.nesdev.org/wiki/APU_Frame_Counter

type frameCounter struct {
	upd chan byte

	quarter chan interface{}
	half    chan interface{}
	irq     chan bool
}

func runFrameCounter(ctx context.Context, divider *divider) *frameCounter {
	upd := make(chan byte)
	quarter := make(chan interface{})
	half := make(chan interface{})
	irq := make(chan bool)

	go func() {
		defer close(quarter)
		defer close(half)
		defer close(irq)
		defer close(upd)

		interruptInhibit := false
		step := 0

		seq4 := func() {
			switch step % 4 {
			case 0:
				quarter <- struct{}{}
			case 1:
				quarter <- struct{}{}
				half <- struct{}{}
			case 2:
				quarter <- struct{}{}
			case 3:
				quarter <- struct{}{}
				half <- struct{}{}
				if !interruptInhibit {
					irq <- true
				}
			}
			step++
		}
		seq5 := func() {
			switch step % 5 {
			case 0:
				quarter <- struct{}{}
			case 1:
				quarter <- struct{}{}
				half <- struct{}{}
			case 2:
				quarter <- struct{}{}
			case 3:
			case 4:
				quarter <- struct{}{}
				half <- struct{}{}
			}
			step++
		}

		seq := seq4
		for {
			select {
			case <-ctx.Done():
				return
			case b := <-upd:
				if util.IsSet(b, 7) {
					seq = seq5
				} else {
					seq = seq4
				}
				interruptInhibit = util.IsSet(b, 6)
				if interruptInhibit {
					irq <- false
				}
			case _, ok := <-divider.output():
				if !ok {
					return
				}
				seq()
			}
		}
	}()
	return &frameCounter{upd, quarter, half, irq}
}

func (c *frameCounter) update(b byte)                    { c.upd <- b }
func (c *frameCounter) quarterFrame() <-chan interface{} { return c.quarter }
func (c *frameCounter) halfFrame() <-chan interface{}    { return c.half }
func (c *frameCounter) frameInterrupt() <-chan bool      { return c.irq }
