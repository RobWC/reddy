package main

import (
	"image"
	"image/color"
)

type squareMask struct {
	p image.Point
}

func NewSquareMask(x int, y int) *squareMask {
	return &squareMask{p: image.Point{X: x, Y: y}}
}

func (sm *squareMask) ColorModel() color.Model {
	return color.AlphaModel
}

func (sm *squareMask) Bounds() image.Rectangle {
	return image.Rect(sm.p.X, sm.p.Y, 128, 128)
}

func (sm *squareMask) At(x, y int) color.Color {
	xx, yy := float64(x-sm.p.X), float64(y-sm.p.Y)
	if xx*xx+yy*yy < float64(x*x+y*y) {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}
