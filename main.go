package main

import (
	"fmt"
	"image"
	"image/color"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"

	_ "embed"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ncruces/zenity"
	"github.com/vibrantgio/ivg/decode"
	"github.com/vibrantgio/ivg/raster/img"
	"github.com/vibrantgio/ivg/render"
)

//go:embed checkerboard.kage
var checkerboardKage []byte

const (
	screenWidth  = 640
	screenHeight = 480
)

type Game struct {
	checkerboardShader  *ebiten.Shader
	checkerboardEnabled bool
	ivgImage            *ebiten.Image
	disasmLines         []string
	disasmEnabled       bool
	disasmScrollOffset  int
	err                 error
	loading             bool
	m                   sync.Mutex
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
			g.m.Lock()
			g.loading = true
			g.m.Unlock()
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

			g.m.Lock()
			g.loading = true
			g.m.Unlock()

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

	if inpututil.IsKeyJustPressed(ebiten.KeyB) && !g.disasmEnabled {
		g.m.Lock()
		g.checkerboardEnabled = !g.checkerboardEnabled
		g.m.Unlock()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.m.Lock()
		g.disasmEnabled = !g.disasmEnabled
		g.m.Unlock()
	}

	if g.disasmEnabled {
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.disasmScrollOffset++
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.disasmScrollOffset--
		}

		_, yoff := ebiten.Wheel()
		if yoff > 0 {
			g.disasmScrollOffset -= int(math.Ceil(yoff))
		} else if yoff < 0 {
			g.disasmScrollOffset += int(math.Ceil(-yoff))
		}

		g.disasmScrollOffset = max(0, g.disasmScrollOffset)
	}

	return nil
}

func (g *Game) loadImage(path string, b []byte) error {
	vb, err := decode.DecodeViewBox(b)
	if err != nil {
		return fmt.Errorf("error decoding IconVG data: %w", err)
	}

	minX, minY, maxX, maxY := vb.AspectMeet(
		float32(screenWidth),
		float32(screenHeight),
		0.5, 0.5,
	)

	width := int(maxX - minX)
	height := int(maxY - minY)

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	var z render.Renderer
	z.SetRasterizer(img.NewRasterizer(rgba), rgba.Bounds())
	if err := decode.Decode(&z, b); err != nil {
		return fmt.Errorf("error decoding IconVG data: %w", err)
	}

	eimg := ebiten.NewImageFromImage(rgba)

	disasm, err := decode.Disassemble(b)
	if err != nil {
		return fmt.Errorf("error disassembling IconVG data: %w", err)
	}

	ebiten.SetWindowTitle(fmt.Sprintf("%s - ivgview", filepath.Base(path)))

	g.m.Lock()
	g.ivgImage = eimg
	g.loading = false
	g.disasmLines = strings.Split(string(disasm), "\n")
	g.m.Unlock()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.m.Lock()
	defer g.m.Unlock()

	if g.loading {
		ebitenutil.DebugPrint(screen, "Loading...")
		return
	}

	if g.ivgImage == nil {
		ebitenutil.DebugPrint(screen, "Drop an IconVG file here to view it;\nClick or press Enter to open a file dialog.")
		return
	}

	if g.disasmEnabled {
		const lineHeight = 16
		visibleLines := screen.Bounds().Dy() / lineHeight

		maxScroll := max(0, len(g.disasmLines)-visibleLines)
		g.disasmScrollOffset = min(g.disasmScrollOffset, maxScroll)

		end := min(g.disasmScrollOffset+visibleLines, len(g.disasmLines))
		visible := strings.Join(g.disasmLines[g.disasmScrollOffset:end], "\n")

		ebitenutil.DebugPrint(screen, visible)
		return
	}

	if g.checkerboardEnabled {
		op := &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"CellSize":   float32(8),
				"LightColor": []float32{0.85, 0.85, 0.85, 1},
				"DarkColor":  []float32{0.65, 0.65, 0.65, 1},
			},
		}
		screen.DrawRectShader(screenWidth, screenHeight, g.checkerboardShader, op)
	} else {
		screen.Fill(color.Black)
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
	if len(os.Args) > 2 {
		fmt.Printf("Usage: %s file.ivg\n", os.Args[0])
		os.Exit(1)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("ivgview")

	checkerboardShader, err := ebiten.NewShader(checkerboardKage)
	if err != nil {
		fmt.Printf("Error compiling checkerboard shader: %v\n", err)
		os.Exit(1)
	}

	g := &Game{checkerboardShader: checkerboardShader}

	if len(os.Args) == 2 {
		g.loading = true
		ivgData, err := os.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			os.Exit(1)
		}

		if err := g.loadImage(os.Args[1], ivgData); err != nil {
			fmt.Printf("Error loading image: %v\n", err)
			os.Exit(1)
		}
	}

	if err := ebiten.RunGame(g); err != nil {
		fmt.Printf("Error running game loop: %v\n", err)
		os.Exit(1)
	}
}
