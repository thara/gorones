package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

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

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("gorones")
	if err := ebiten.RunGame(emu); err != nil {
		log.Fatal(err)
	}
}
