package main

import (
	"flag"
	"github.com/hochbaum/vampire-survivors-tools/texturepacker"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"os"
	"path/filepath"
)

type cropper interface {
	SubImage(r image.Rectangle) image.Image
}

func cropFrames(filePath string, size int, texture texturepacker.PackedTexture) (map[string]image.Image, error) {
	directory := filepath.Dir(filePath)
	imgFile, err := os.Open(filepath.Join(directory, texture.Image))
	if err != nil {
		panic(err)
	}
	defer imgFile.Close()

	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}

	images := make(map[string]image.Image)
	for _, frame := range texture.Frames {
		cropped := img.(cropper).SubImage(frame.Frame.Rect())
		cropped = resize.Resize(
			uint(resizeInt(frame.SourceSize.X, size)),
			uint(resizeInt(frame.SourceSize.Y, size)),
			cropped, resize.NearestNeighbor)
		images[frame.FileName] = cropped
	}
	return images, nil
}

func resizeInt(i, resizePercent int) int {
	return int(float64(i) / 100.0 * float64(resizePercent))
}

func writeImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	out := flag.String("o", wd, "Specifies the output path.")
	size := flag.Int("size", 100, "Specifies size of the images.")
	flag.Parse()

	path := flag.Arg(0)
	sheet, err := texturepacker.Open(path)
	if err != nil {
		panic(err)
	}

	// The game's packed textures only contain a single entry right now.
	texture := sheet.Textures[0]
	images, err := cropFrames(path, *size, texture)
	if err != nil {
		panic(err)
	}

	_, err = os.Stat(*out)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(*out, os.ModeDir); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	for name, img := range images {
		if err := writeImage(filepath.Join(*out, name), img); err != nil {
			panic(err)
		}
	}
}
