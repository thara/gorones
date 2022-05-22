package apu

import (
	"context"
)

type lengthCounter struct {
	clk       chan interface{}
	enabledCh chan bool
	haltCh    chan bool
	loadCh    chan uint8
	counterCh chan chan uint8
}

func runLengthCounter(ctx context.Context) *lengthCounter {
	var (
		clk       = make(chan interface{})
		enabledCh = make(chan bool)
		haltCh    = make(chan bool)
		loadCh    = make(chan uint8)
		counterCh = make(chan chan uint8)
	)
	go func() {
		defer close(clk)
		defer close(enabledCh)
		defer close(haltCh)
		defer close(loadCh)
		defer close(counterCh)

		var (
			enabled bool
			halt    bool
			counter uint8
		)
		for {
			select {
			case <-ctx.Done():
				return
			case <-clk:
				if counter != 0 && !halt {
					counter--
				}
			case enabled = <-enabledCh:
				if !enabled {
					counter = 0
				}
			case halt = <-haltCh:
			case v := <-loadCh:
				if enabled {
					n := (v & 0b11111000) >> 3
					counter = lengthTable[n]
				}

			case ch := <-counterCh:
				ch <- counter
				close(ch)
			}
		}
	}()
	return &lengthCounter{clk, enabledCh, haltCh, loadCh, counterCh}
}

func (c *lengthCounter) load(v uint8)      { c.loadCh <- v }
func (c *lengthCounter) setEnabled(v bool) { c.enabledCh <- v }
func (c *lengthCounter) setHalt(v bool)    { c.haltCh <- v }
func (c *lengthCounter) clock()            { c.clk <- struct{}{} }
func (c *lengthCounter) counter() uint8 {
	ch := make(chan uint8)
	c.counterCh <- ch
	return <-ch
}

var lengthTable = []uint8{
	10, 254, 20, 2, 40, 4, 80, 6, 160, 8, 60, 10, 14, 12, 26, 14,
	12, 16, 24, 18, 48, 20, 96, 22, 192, 24, 72, 26, 16, 28, 32, 30,
}
