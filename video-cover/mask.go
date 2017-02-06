package cover

import (
	"image"
	"image/color"
	"math"
)

type gradmask struct {
	h int
	w int
}

func (m *gradmask) ColorModel() color.Model {
	return color.AlphaModel
}

func (m *gradmask) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.w, m.h)
}

func (m *gradmask) At(x, y int) color.Color {
	if x > m.w/2 {
		return color.Alpha{255}
	}

	if x < m.w/4 {
		return color.Alpha{0}
	}

	return color.Alpha{uint8(255.0 * math.Pow(float64(x-m.w/4)/float64(m.w/4), 2))}
}
