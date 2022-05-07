package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/thara/gorones/ppu"
)

var nestest bool

func init() {
	flag.BoolVar(&nestest, "nestest", false, "init for nestest")
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] ROM\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	path := os.Args[1]

	emu, err := newEmulator(path)
	if err != nil {
		log.Fatalf("fail to initialize emulator for %s: %v", path, err)
	}

	scale := 4

	ebiten.SetWindowSize(ppu.WIDTH*scale, ppu.HEIGHT*scale)
	ebiten.SetWindowTitle("gorones")
	if err := ebiten.RunGame(emu); err != nil {
		log.Fatal(err)
	}
}
