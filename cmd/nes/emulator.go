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

	ctrl1 *kbStdCtrl
	ctrl2 *kbStdCtrl

	renderer *renderer
}

func newEmulator(path string, audio *Audio) (*Emulator, error) {
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
	fmt.Println(m)

	ctrl1 := newKbStdCtrl()
	ctrl2 := newKbStdCtrl()

	renderer := newRenderer()

	var emu Emulator
	emu.renderer = renderer
	emu.ctrl1 = ctrl1
	emu.ctrl2 = ctrl2

	emu.nes = gorones.NewNES(m, ctrl1.ctrl, ctrl2.ctrl, renderer, audio)
	emu.nes.PowerOn()

	if nestest {
		fmt.Println("init for nestest")
		emu.nes.InitNEStest()
	} else {
		fmt.Println("reset")
		emu.nes.Reset()
	}

	return &emu, nil
}

func (e *Emulator) Update() error {
	e.ctrl1.update()
	e.ctrl2.update()
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
