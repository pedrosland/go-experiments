package main

import (
	"image"
	"image/color"
)

func clamp(val int, min int, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}

type NRGBASum struct {
	R, G, B, A uint32
}

func (s *NRGBASum) Add(c *color.NRGBA) {
	s.R += uint32(c.R)
	s.G += uint32(c.G)
	s.B += uint32(c.B)
	s.A += uint32(c.A)
}

func (s NRGBASum) ToNRGBA() *color.NRGBA {
	return &color.NRGBA{uint8(s.R), uint8(s.G), uint8(s.B), uint8(s.A)}
}

func BoxBlurKernel(src *image.NRGBA, x int, y int, radius int) *color.NRGBA {
	b := src.Bounds()

	startX := clamp(x-radius, 0, b.Dx()-1)
	startY := clamp(y-radius, 0, b.Dy()-1)

	endX := clamp(x+radius, 0, b.Dx()-1)
	endY := clamp(y+radius, 0, b.Dy()-1)

	curX := startX
	curY := startY
	total := NRGBASum{0, 0, 0, 0}
	numPixels := uint32(0)

	for curX <= endX {
		for curY <= endY {
			c := src.At(curX, curY).(color.NRGBA)
			total.Add(&c)
			numPixels++

			curY++
		}
		curX++
		curY = startY
	}

	return NRGBASum{total.R / numPixels, total.G / numPixels, total.B / numPixels, total.A / numPixels}.ToNRGBA()
}
