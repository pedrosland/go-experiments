package main

import (
	"image"
	"math"
	"sync"
)

func Blur(src *image.NRGBA, radius int) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(b)

	width := b.Dx()

	verticalBlur(src, dst, radius, 0, width)

	return dst
}

func ParallelBlur(src *image.NRGBA, radius, numTasks int) *image.NRGBA {
	b := src.Bounds()
	dst := image.NewNRGBA(b)

	wg := sync.WaitGroup{}

	width := b.Dx()
	i := 0

	step := int(math.Ceil(float64(width) / math.Min(float64(numTasks), float64(width))))
	start := 0
	end := start + step

	for i < width {
		wg.Add(1)
		go func(start, end int){
			verticalBlur(src, dst, radius, start, end)
			wg.Done()
		}(start, end)

		i += step
		start = end

		if end + step < width {
			end += step
		} else {
			end = width
		}
	}

	wg.Wait()

	return dst
}


func verticalBlur(src, dst *image.NRGBA, radius, start, end int) *image.NRGBA {
	b := src.Bounds()

	height := b.Dy()

	for x := start; x < end; x++ {
		for y := 0; y < height; y++ {
			c := *BoxBlurKernel(src, x, y, radius)
			dst.SetNRGBA(x, y, c)
		}
	}

	return dst
}

