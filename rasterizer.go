package main

/*
var _ raster.Rasterizer = (*rasterizer)(nil)

type rasterizer struct {
	path vector.Path
	posX float32
	posY float32
	size image.Point
	Dst  *ebiten.Image
}

func (r *rasterizer) Reset(w, h int) {
	r.path = vector.Path{}
	r.size = image.Point{X: w, Y: h}
	r.posX = 0
	r.posY = 0
}

func (r *rasterizer) Size() image.Point {
	return r.size
}

func (r *rasterizer) Bounds() image.Rectangle {
	return image.Rect(0, 0, r.size.X, r.size.Y)
}

func (r *rasterizer) Pen() (x, y float32) {
	return r.posX, r.posY
}

func (r *rasterizer) MoveTo(ax, ay float32) {
	r.path.MoveTo(ax, ay)
	r.posX = ax
	r.posY = ay
}

func (r *rasterizer) LineTo(bx, by float32) {
	r.path.LineTo(bx, by)
	r.posX = bx
	r.posY = by
}

func (r *rasterizer) QuadTo(bx, by, cx, cy float32) {
	r.path.QuadTo(bx, by, cx, cy)
	r.posX = cx
	r.posY = cy
}

func (r *rasterizer) CubeTo(bx, by, cx, cy, dx, dy float32) {
	r.path.CubicTo(bx, by, cx, cy, dx, dy)
	r.posX = dx
	r.posY = dy
}

func (r *rasterizer) ClosePath() {
	r.path.Close()
}

func (r *rasterizer) Draw(rect image.Rectangle, src image.Image, sp image.Point) {
	if r.Dst == nil {
		panic("ivgview: Rasterizer.Dst is nil")
	}
	if rect.Empty() {
		r.path.Reset()
		return
	}

	bounds := r.Dst.Bounds()
	source := ebiten.NewImage(rect.Dx(), rect.Dy())
	defer source.Deallocate()

	switch src := src.(type) {
	case *ebiten.Image:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(src.Bounds().Min.X-sp.X), float64(src.Bounds().Min.Y-sp.Y))
		source.DrawImage(src, op)
	case *image.Uniform:
		source.Fill(src.C)
	default:
		converted := ebiten.NewImage(rect.Dx(), rect.Dy())
		defer converted.Deallocate()
		for x := 0; x < rect.Dx(); x++ {
			for y := 0; y < rect.Dy(); y++ {
				converted.Set(x, y, src.At(x+sp.X, y+sp.Y))
			}
		}
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(src.Bounds().Min.X-sp.X), float64(src.Bounds().Min.Y-sp.Y))
		source.DrawImage(converted, op)
	}

	mask := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	defer mask.Deallocate()
	vector.FillPath(mask, &r.path, nil, &vector.DrawPathOptions{AntiAlias: true})
	r.path.Reset()

	layer := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	defer layer.Deallocate()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(rect.Min.X), float64(rect.Min.Y))
	layer.DrawImage(source, op)

	op = &ebiten.DrawImageOptions{Blend: ebiten.BlendDestinationIn}
	layer.DrawImage(mask, op)

	r.Dst.DrawImage(layer, nil)
}
*/
