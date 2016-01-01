// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 58.
//!+

// Surface computes an SVG rendering of a 3-D surface function.
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
)

const (
	cells   = 100         // number of grid cells
	xyrange = 30.0        // axis ranges (-xyrange..+xyrange)
	angle   = math.Pi / 6 // angle of x, y axes (=30°)
)

type imageInfo struct {
	width   float64
	height  float64
	xyscale float64 // pixels per x or y unit
	zscale  float64 // pixels per z unit
}

var sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)

func main() {
	//createSvg(os.Stdout, newImageInfo(600, 320))

	http.HandleFunc("/image.svg", handleImageRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleImageRequest(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	width, err := strconv.ParseFloat(query.Get("width"), 64)
	if err != nil {
		width = 600
	}

	height, err := strconv.ParseFloat(query.Get("height"), 64)
	if err != nil {
		height = 320
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	createSvg(w, newImageInfo(width, height))
}

func newImageInfo(width, height float64) imageInfo {
	return imageInfo{
		width:   width,
		height:  height,
		xyscale: width / 2 / xyrange, // pixels per x or y unit
		zscale:  height * 0.4,        // pixels per z unit
	}
}

func createSvg(w io.Writer, info imageInfo) {
	fmt.Fprintf(w, "<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; stroke-width: 0.7' "+
		"width='%d' height='%d'>", int(info.width), int(info.height))

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay, z1, ok1 := corner(i+1, j, info)
			bx, by, z2, ok2 := corner(i, j, info)
			cx, cy, z3, ok3 := corner(i, j+1, info)
			dx, dy, z4, ok4 := corner(i+1, j+1, info)

			if !ok1 || !ok2 || !ok3 || !ok4 {
				continue
			}

			color := calculateColor(average(z1, z2, z3, z4), info)

			fmt.Fprintf(w, "<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='%s'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, color)
		}
	}

	fmt.Fprintln(w, "</svg>")
}

func corner(i, j int, info imageInfo) (float64, float64, float64, bool) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/cells - 0.5)
	y := xyrange * (float64(j)/cells - 0.5)

	// Compute surface height z.
	z := f(x, y)
	if math.IsNaN(z) {
		return 0, 0, 0, false
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := info.width/2 + (x-y)*cos30*info.xyscale
	sy := info.height/2 + (x+y)*sin30*info.xyscale - z*info.zscale
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

func calculateColor(z float64, info imageInfo) string {
	// z has range of -0.22 - 0.99
	const minColor = 0xff0000
	const maxColor = 0x0000ff
	const colorRange = 256
	const zRange = 1.22

	bitColor := int((z + 0.22) / zRange * colorRange)
	colorR := bitColor
	colorB := 0xff - bitColor
	hexColor := fmt.Sprintf("#%02x00%02x", colorR, colorB)

	return hexColor
}

//!-
