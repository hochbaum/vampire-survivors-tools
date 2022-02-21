package main

import (
	"flag"
	"image"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/andybons/gogif"
	"github.com/hochbaum/vampire-survivors-tools/texturepacker"
	"github.com/nfnt/resize"
)

var gifNameExp1 = regexp.MustCompile(`^(.*)_(\d*)\.png`)
var gifNameExp2 = regexp.MustCompile(`(.*)(\d)\.png`)

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
		cropped := img.(cropper).SubImage(frame.Frame)
		if cropped.Bounds().Dx() <= 6 && cropped.Bounds().Dy() <= 6 {
			continue
		}
		cropped = resize.Resize(
			uint(resizeInt(frame.SourceSize.Width, size)),
			uint(resizeInt(frame.SourceSize.Height, size)),
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

func writeImages(path string, images map[string]image.Image) error {
	for name, img := range images {
		if err := writeImage(filepath.Join(path, name), img); err != nil {
			return err
		}
	}
	return nil
}

// Returns imges in order according to their position in the gif
func orderImages(images map[string]image.Image) (map[string][]image.Image, error) {
	gifs := make(map[string][]image.Image)
	for name, img := range images {
		parts := gifNameExp1.FindStringSubmatch(name)
		// Checking if image is part of a gif
		if len(parts) != 3 {
			parts = gifNameExp2.FindStringSubmatch(name)
			if len(parts) != 3 {
				continue
			}
		}
		gifName, gifOrderRaw := parts[1], parts[2]
		gifOrder, _ := strconv.Atoi(gifOrderRaw)
		// Allocating enough memory for 128 frames
		if gifs[gifName] == nil {
			gifs[gifName] = make([]image.Image, 128)
		}
		gifs[gifName][gifOrder] = img
	}
	return gifs, nil
}

// Normalizes images according to the largest image, creates palette
func normalizeImages(images []image.Image) ([]*image.Paletted, error) {
	w, h := 0, 0
	quantizer := gogif.MedianCutQuantizer{NumColor: 64}
	normalized := make([]*image.Paletted, 0)
	// Finding biggest image bounds
	for _, img := range images {
		if img == nil {
			continue
		}
		dx, dy := img.Bounds().Dx(), img.Bounds().Dy()
		if dx > w {
			w = dx
		}
		if dy > h {
			h = dy
		}
	}
	bounds := image.Rect(0, 0, w, h)
	for _, img := range images {
		if img == nil {
			continue
		}
		srcBounds := img.Bounds()
		fixedSize := image.NewRGBA(bounds)
		// Fitting image into max bounds
		r := image.Rectangle{
			image.Point{X: bounds.Dx() - srcBounds.Dx(), Y: bounds.Dy() - srcBounds.Dy()},
			image.Point{X: bounds.Dx(), Y: bounds.Dy()}}
		draw.Draw(fixedSize, r, img, image.Point{}, draw.Src)
		palettedImage := image.NewPaletted(bounds, nil)
		quantizer.Quantize(palettedImage, bounds, fixedSize, image.Point{})
		normalized = append(normalized, palettedImage)
	}
	return normalized, nil
}

func writeGif(path string, imgs []*image.Paletted) error {
	outGif := &gif.GIF{}
	for _, img := range imgs {
		// Setting gif settings for each frame
		outGif.Image = append(outGif.Image, img)
		outGif.Delay = append(outGif.Delay, 20)
		outGif.Disposal = append(outGif.Disposal, 0x02)
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, outGif)
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	out := flag.String("o", wd, "Specifies the output path.")
	size := flag.Int("size", 100, "Specifies size of the images.")
	gifFlag := flag.Bool("gif", false, "Creates gifs from connected frames.")
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

	if *gifFlag {
		gifs, err := orderImages(images)
		if err != nil {
			panic(err)
		}
		// Normalizing and writing each gif
		for name, imageSeries := range gifs {
			normalized, err := normalizeImages(imageSeries)
			if err != nil {
				panic(err)
			}
			if err := writeGif(filepath.Join(*out, name+".gif"), normalized); err != nil {
				panic(err)
			}
		}
	} else {
		err := writeImages(*out, images)
		if err != nil {
			panic(err)
		}
	}
}
