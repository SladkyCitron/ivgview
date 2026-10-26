package main

import (
	"fmt"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/vibrantgio/ivg/decode"
	"github.com/vibrantgio/ivg/raster/img"
	"github.com/vibrantgio/ivg/render"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	ivgData []byte
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	var z render.Renderer
	z.SetRasterizer(img.NewRasterizer(screen), screen.Bounds())
	if err := decode.Decode(&z, g.ivgData); err != nil {
		fmt.Printf("Error decoding IVG data: %v\n", err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file.ivg\n", os.Args[0])
		os.Exit(1)
	}

	ivgData, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	g := &Game{ivgData: ivgData}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("ivgview")
	if err := ebiten.RunGame(g); err != nil {
		fmt.Printf("Error running game loop: %v\n", err)
		os.Exit(1)
	}
}
