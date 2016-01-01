// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"math"
	"os"
)

const (
	width, height = 600, 320            // canvas size in pixels
	cells         = 100                 // number of grid cells
	xyrange       = 30.0                // axis ranges (-xyrange..+xyrange)
	xyscale       = width / 2 / xyrange // pixels per x or y unit
	zscale        = height * 0.4        // pixels per z unit
	angle         = math.Pi / 6         // angle of x, y axes (=30°)
)

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	fmt.Printf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", width, height)
	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, z1, ok1 := corner(i+1, j)
			bx, by, z2, ok2 := corner(i, j)
			cx, cy, z3, ok3 := corner(i, j+1)
			dx, dy, z4, ok4 := corner(i+1, j+1)

			if !ok1 || !ok2 || !ok3 || !ok4 {
				continue
			}

			color := calculateColor(average(z1, z2, z3, z4))

			fmt.Printf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}
	fmt.Println("</svg>")
}

func corner(i, j int) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	if math.IsNaN(z) {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale
	return sx, sy, z, true
}

func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}

func average(nums ...float64) float64 {
	var total, num float64
	var i int

	for i, num = range nums {
		total += num
	}

	return total / float64(i+1)
}

func calculateColor(z float64) string {
	// z has range of -0.22 - 0.99
	const minColor = 0xff0000
	const maxColor = 0x0000ff
	const colorRange = 256
	const zRange = 1.22

	bitColor := int((z + 0.22) / zRange * colorRange)
	colorR := bitColor
	colorB := 0xff - bitColor
	hexColor := fmt.Sprintf("#%02x00%02x", colorR, colorB)

	fmt.Fprintf(os.Stderr, "%v #%06x %s\n", bitColor, (colorR<<16)|colorB, hexColor)

	return hexColor
}

//!-
