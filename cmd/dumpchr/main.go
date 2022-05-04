package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/thara/gorones/mapper"
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
	if err := run(path); err != nil {
		log.Fatal(err)
	}
}

func run(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", path, err)
	}
	defer f.Close()

	rom, err := mapper.ParseROM(f)
	if err != nil {
		return fmt.Errorf("fail to open %s: %v", path, err)
	}

	m, err := rom.Mapper()
	if err != nil {
		return fmt.Errorf("fail to get mapper %s: %v", path, err)
	}

	ps := loadCHRPatterns(m)
	fmt.Printf("number of patterns: %d\n", len(ps))

	img := image.NewRGBA(image.Rectangle{Max: image.Point{X: len(ps) * 8, Y: 8}})

	for i, p := range ps {
		p.write(i, img)
	}

	f, _ = os.Create("chr.png")
	if err != nil {
		return fmt.Errorf("fail to create file: %v", err)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		return fmt.Errorf("fail to create png file: %v", err)
	}
	return nil
}
