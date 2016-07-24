package main

import (
	"os"
	"log"
	"image"
	"image/png"
	_ "image/jpeg"
	"flag"
	"image/draw"
)

func main() {
	inputFile := flag.String("input", "", "Input file to blur")
	radius := flag.Int("radius", 1, "Blur radius")
	numTasks := flag.Int("tasks", 0, "Number of concurrent tasks - 0 to use Blur()")

	flag.Parse()

	if *inputFile == "" {
		log.Fatal("Please specify an input file with -input")
	}

	reader, err := os.Open(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	castNRGBA, ok := img.(*image.NRGBA)
	if !ok {
		castNRGBA = image.NewNRGBA(img.Bounds())
		draw.Draw(castNRGBA, img.Bounds(), img, image.Point{0, 0}, draw.Src)
	}

	var dest *image.NRGBA
	if *numTasks < 1 {
		dest = Blur(castNRGBA, 1)
	} else {
		dest = ParallelBlur(castNRGBA, *radius, *numTasks)
	}

	writer, err := os.Create("./result.png")
	if err != nil {
		log.Fatal(err)
	}

	png.Encode(writer, dest)
}
