package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/thara/gorones"
	"github.com/thara/gorones/input"
	"github.com/thara/gorones/mapper"
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

	nes, err := newNES(path)
	if err != nil {
		log.Fatalf("fail to initialize emulator for %s: %v", path, err)
	}
	nes.PowerOn()

	if nestest {
		fmt.Println("init for nestest")
		nes.InitNEStest()
	} else {
		fmt.Println("reset")
		nes.Reset()
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		nes.RunFrame()
	}
}

func newNES(path string) (*gorones.NES, error) {
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

	ctrl := new(input.StandardController)

	nes := gorones.NewNES(m, ctrl, ctrl, new(renderer))
	return nes, nil
}

type renderer struct{}

func (r *renderer) UpdateFrame(buf *[ppu.WIDTH * ppu.HEIGHT]uint8) {
	for i, v := range buf {
		if i%ppu.WIDTH == 0 {
			fmt.Printf("\n%03d", i/ppu.WIDTH)
		}
		c := "."
		if palette[v] != 0x000000FF {
			c = "*"
		}
		fmt.Print(c)
	}
	fmt.Print("\n=================================================\n")
}

var palette [255]uint32 = [255]uint32{
	0x7C7C7CFF, 0x0000FCFF, 0x0000BCFF, 0x4428BCFF, 0x940084FF, 0xA80020FF, 0xA81000FF, 0x881400FF,
	0x503000FF, 0x007800FF, 0x006800FF, 0x005800FF, 0x004058FF, 0x000000FF, 0x000000FF, 0x000000FF,
	0xBCBCBCFF, 0x0078F8FF, 0x0058F8FF, 0x6844FCFF, 0xD800CCFF, 0xE40058FF, 0xF83800FF, 0xE45C10FF,
	0xAC7C00FF, 0x00B800FF, 0x00A800FF, 0x00A844FF, 0x008888FF, 0x000000FF, 0x000000FF, 0x000000FF,
	0xF8F8F8FF, 0x3CBCFCFF, 0x6888FCFF, 0x9878F8FF, 0xF878F8FF, 0xF85898FF, 0xF87858FF, 0xFCA044FF,
	0xF8B800FF, 0xB8F818FF, 0x58D854FF, 0x58F898FF, 0x00E8D8FF, 0x787878FF, 0x000000FF, 0x000000FF,
	0xFCFCFCFF, 0xA4E4FCFF, 0xB8B8F8FF, 0xD8B8F8FF, 0xF8B8F8FF, 0xF8A4C0FF, 0xF0D0B0FF, 0xFCE0A8FF,
	0xF8D878FF, 0xD8F878FF, 0xB8F8B8FF, 0xB8F8D8FF, 0x00FCFCFF, 0xF8D8F8FF, 0x000000FF, 0x000000FF}
