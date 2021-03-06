package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/thara/gorones/input"
)

// kbStdCtrl is standard controller emulated by keyboard
type kbStdCtrl struct {
	ctrl *input.StandardController

	keys []ebiten.Key
}

func newKbStdCtrl() *kbStdCtrl {
	return &kbStdCtrl{
		ctrl: new(input.StandardController),
		keys: make([]ebiten.Key, 0, 8),
	}
}

func (c *kbStdCtrl) update() {
	c.keys = inpututil.AppendPressedKeys(c.keys[:0])

	var state uint8
	for _, k := range c.keys {
		switch k {
		case ebiten.KeyW:
			state |= input.StandardUp
		case ebiten.KeyA:
			state |= input.StandardLeft
		case ebiten.KeyS:
			state |= input.StandardDown
		case ebiten.KeyD:
			state |= input.StandardRight

		case ebiten.KeyJ:
			state |= input.StandardB
		case ebiten.KeyK:
			state |= input.StandardA

		case ebiten.KeySpace:
			state |= input.StandardStart
		case ebiten.KeyEnter:
			state |= input.StandardSelect
		}
	}
	c.ctrl.Update(state)
}
