package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/thara/gorones"
	"github.com/thara/gorones/mapper"
	"github.com/thara/gorones/ppu"
)

type Emulator struct {
	nes *gorones.NES

	renderer *renderer
}

func newEmulator(path string) (*Emulator, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fail to open %s: %v", path, err)
	}
	defer f.Close()

	rom, err := mapper.ParseROM(f)
	if err != nil {
		return nil, fmt.Errorf("fail to open %s: %v", path, err)
	}

	m, err := rom.Mapper()
	if err != nil {
		return nil, fmt.Errorf("fail to get mapper %s: %v", path, err)
	}

	ctrl1 := newKbStdCtrl()
	ctrl2 := newKbStdCtrl()

	renderer := newRenderer()

	var emu Emulator
	emu.renderer = renderer

	emu.nes = gorones.NewNES(m, ctrl1.ctrl, ctrl2.ctrl, renderer)
	emu.nes.PowerOn()

	return &emu, nil
}

func (e *Emulator) Update() error {
	e.nes.RunFrame()
	return nil
}

func (e *Emulator) Draw(screen *ebiten.Image) {
	screen.ReplacePixels(e.renderer.pixels())
	ebitenutil.DebugPrint(screen, fmt.Sprintf("tps: %f", ebiten.CurrentTPS()))
}

func (e *Emulator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ppu.WIDTH, ppu.HEIGHT
}
