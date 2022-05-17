package apu

import (
	"context"
)

// https://www.nesdev.org/wiki/APU#Glossary

type divider struct {
	in  chan interface{}
	out chan interface{}
}

func runDivider(ctx context.Context, p int) *divider {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		defer close(in)
		defer close(out)

		n := p
		for {
			select {
			case <-ctx.Done():
				return
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
	return &divider{in, out}
}

func (d *divider) clock() {
	d.in <- struct{}{}
}

func (d *divider) output() <-chan interface{} {
	return d.out
}
