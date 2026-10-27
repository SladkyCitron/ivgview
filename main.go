package main

import (
	"fmt"
	"image"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ncruces/zenity"
	"github.com/vibrantgio/ivg/decode"
	"github.com/vibrantgio/ivg/raster/img"
	"github.com/vibrantgio/ivg/render"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	ivgImage *ebiten.Image
	err      error
	m        sync.Mutex
}

func (g *Game) Update() error {
	if err := func() error {
		g.m.Lock()
		defer g.m.Unlock()
		return g.err
	}(); err != nil {
		return err
	}

	if files := ebiten.DroppedFiles(); files != nil {
		go func() {
			if err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				fi, err := d.Info()
				if err != nil {
					return err
				}

				if fi.IsDir() {
					return nil
				}

				ivgData, err := fs.ReadFile(files, path)
				if err != nil {
					return err
				}

				if err := g.loadImage(path, ivgData); err != nil {
					return err
				}

				return nil
			}); err != nil {
				g.m.Lock()
				if g.err == nil {
					g.err = err
				}
				g.m.Unlock()
			}
		}()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) ||
		inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		go func() {
			path, _ := zenity.SelectFile(
				zenity.Modal(),
				zenity.FileFilters{
					{"IconVG file", []string{"*.ivg", "*.iconvg"}, true},
				},
			)
			if path == "" {
				return
			}

			ivgData, err := os.ReadFile(path)
			if err != nil {
				g.m.Lock()
				if g.err == nil {
					g.err = err
				}
				g.m.Unlock()
			}

			if err := g.loadImage(path, ivgData); err != nil {
				g.m.Lock()
				if g.err == nil {
					g.err = err
				}
				g.m.Unlock()
			}
		}()
	}

	return nil
}

func (g *Game) loadImage(path string, b []byte) error {
	ivgImage := image.NewRGBA(image.Rect(0, 0, screenWidth, screenHeight))
	var z render.Renderer
	z.SetRasterizer(img.NewRasterizer(ivgImage), ivgImage.Bounds())
	if err := decode.Decode(&z, b); err != nil {
		return fmt.Errorf("error decoding IconVG data: %w", err)
	}

	eimg := ebiten.NewImageFromImage(ivgImage)

	g.m.Lock()
	g.ivgImage = eimg
	g.m.Unlock()

	ebiten.SetWindowTitle(fmt.Sprintf("%s - ivgview", filepath.Base(path)))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.ivgImage == nil {
		ebitenutil.DebugPrint(screen, "Drop an IconVG file here to view it;\nClick or press Enter to open a file dialog.")
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64((screenWidth-g.ivgImage.Bounds().Dx())/2), float64((screenHeight-g.ivgImage.Bounds().Dy())/2))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(g.ivgImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	/*
		if len(os.Args) < 2 {
			fmt.Printf("Usage: %s file.ivg\n", os.Args[0])
			os.Exit(1)
		}
	*/

	g := &Game{}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ivgview")
	if err := ebiten.RunGame(g); err != nil {
		fmt.Printf("Error running game loop: %v\n", err)
		os.Exit(1)
	}
}
