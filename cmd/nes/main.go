package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/gordonklaus/portaudio"
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

	portaudio.Initialize()
	defer portaudio.Terminate()
	host, err := portaudio.DefaultHostApi()
	if err != nil {
		log.Fatalln(err)
	}
	audio := &Audio{channel: make(chan float32, 44100)}
	stream, err := portaudio.OpenStream(
		portaudio.HighLatencyParameters(nil, host.DefaultOutputDevice),
		func(out []float32) {
			for i := range out {
				select {
				case sample := <-audio.channel:
					out[i] = sample
				default:
					out[i] = 0
				}
			}
		},
	)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stream.Start(); err != nil {
		log.Fatalln(err)
	}
	audio.stream = stream
	defer audio.stream.Close()

	emu, err := newEmulator(path, audio)
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

type Audio struct {
	stream  *portaudio.Stream
	channel chan float32
}

func (a *Audio) Write(v float32) {
	fmt.Println(v)
	select {
	case a.channel <- v:
	default:
	}
}
