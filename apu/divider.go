package apu

import (
	"context"

	"golang.org/x/exp/constraints"
)

// https://www.nesdev.org/wiki/APU#Glossary

type divider[T constraints.Unsigned] struct {
	p   chan T
	in  chan interface{}
	out chan interface{}
}

func runDivider[T constraints.Unsigned](ctx context.Context, p T) *divider[T] {
	reload := make(chan T)

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
	return &divider[T]{reload, in, out}
}

func (d *divider[T]) reload(n T)                 { d.p <- n }
func (d *divider[T]) clock()                     { d.in <- struct{}{} }
func (d *divider[T]) output() <-chan interface{} { return d.out }
