package main_test

import (
	"image"
	"image/color"
	"testing"

	blur "github.com/pedrosland/go-experiments/scala/parallel-box-blur"
)

var src3x3 *image.NRGBA

func init() {
	img := image.NewNRGBA(image.Rect(0, 0, 3, 3))

	i := uint8(0)

	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			img.Set(x, y, color.NRGBA{i, i, i, i})
			i++
		}
	}

	src3x3 = img
}

func check(t *testing.T, img *image.NRGBA, x int, y int, value uint8) {
	c := img.At(x, y).(color.NRGBA)

	if c.R != value || c.G != value || c.B != value || c.A != value {
		t.Errorf("pixel (%d, %d) has value %v expected %d", x, y, c, value)
	}
}

func check3x3Radius1(t *testing.T, img *image.NRGBA) {
	check(t, img, 0, 0, 2)
	check(t, img, 1, 0, 2)
	check(t, img, 2, 0, 3)
	check(t, img, 0, 1, 3)
	check(t, img, 1, 1, 4)
	check(t, img, 2, 1, 4)
	check(t, img, 0, 2, 5)
	check(t, img, 1, 2, 5)
	check(t, img, 2, 2, 6)
}

func Test_Blur3x3Radius1(t *testing.T) {
	dst := blur.Blur(src3x3, 1)

	check3x3Radius1(t, dst)
}

func Test_ParallelBlur3x3Radius1Tasks1(t *testing.T) {
	dst := blur.ParallelBlur(src3x3, 1, 1)

	check3x3Radius1(t, dst)
}

func Test_ParallelBlur3x3Radius1Tasks2(t *testing.T) {
	dst := blur.ParallelBlur(src3x3, 1, 1)

	check3x3Radius1(t, dst)
}
