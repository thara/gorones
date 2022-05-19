package apu

import (
	"context"
)

// https://www.nesdev.org/wiki/APU#Glossary

type divider struct {
	p   chan uint
	in  chan interface{}
	out chan interface{}
}

func runDivider(ctx context.Context, p uint) *divider {
	reload := make(chan uint)

	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		defer close(reload)
		defer close(in)
		defer close(out)

		n := p
		for {
			select {
			case <-ctx.Done():
				return
			case p := <-reload:
				n = p
			case _, ok := <-in:
				if !ok {
					return
				}

				if n == 0 {
					n = p
					out <- struct{}{}
				} else {
					n--
				}
			}
		}
	}()
	return &divider{reload, in, out}
}

func (d *divider) reload(n uint) {
	d.p <- n
}

func (d *divider) clock() {
	d.in <- struct{}{}
}

func (d *divider) output() <-chan interface{} {
	return d.out
}
