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

func writeGif(path string, img *gif.GIF) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return gif.EncodeAll(file, img)
}

func main() {
	gifNameExp1, _ := regexp.Compile(`^(.*)_i*(\d\d)`)
	gifNameExp2, _ := regexp.Compile(`(.*)(\d)\.png`)
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	out := flag.String("o", wd, "Specifies the output path.")
	size := flag.Int("size", 100, "Specifies size of the images.")
	gifFlag := flag.Bool("gif", false, "Creates gifs from connected frames.")
	gifs := make(map[string][]image.Image)
	maxGif := make(map[string]image.Point)
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
		if *gifFlag {
			parts := gifNameExp1.FindStringSubmatch(name)
			// Checking if image is part of a gif
			if len(parts) != 3 {
				parts = gifNameExp2.FindStringSubmatch(name)
			}
			if len(parts) == 3 {
				gifName, gifOrderRaw := parts[1], parts[2]
				gifOrder, _ := strconv.Atoi(gifOrderRaw)
				// Ordering the images in the right sequence
				if gifs[gifName] == nil {
					gifs[gifName] = make([]image.Image, 20)
					maxGif[gifName] = image.Point{}
				}
				//Getting the max image size for the gif
				if img.Bounds().Dx() > maxGif[gifName].X {
					maxGif[gifName] = image.Point{X: img.Bounds().Max.X, Y: maxGif[gifName].Y}
				}
				if img.Bounds().Dy() > maxGif[gifName].Y {
					maxGif[gifName] = image.Point{X: maxGif[gifName].X, Y: img.Bounds().Max.Y}
				}
				gifs[gifName][gifOrder] = img
				continue
			}
		}
		if err := writeImage(filepath.Join(*out, name), img); err != nil {
			panic(err)
		}
	}
	for name, curGif := range gifs {
		outGif := &gif.GIF{}
		bounds := image.Rect(0, 0, maxGif[name].X, maxGif[name].Y)
		quantizer := gogif.MedianCutQuantizer{NumColor: 64}

		for _, simage := range curGif {
			if simage == nil {
				continue
			}
			srcBounds := simage.Bounds()
			fixedSize := image.NewRGBA(bounds)
			r := image.Rectangle{
				image.Point{X: bounds.Dx() - srcBounds.Dx(), Y: bounds.Dy() - srcBounds.Dy()},
				image.Point{X: bounds.Dx(), Y: bounds.Dy()}}
			draw.Draw(fixedSize, r, simage, image.Point{}, draw.Src)
			palettedImage := image.NewPaletted(bounds, nil)
			quantizer.Quantize(palettedImage, bounds, fixedSize, image.Point{})

			// Add new frame to animated GIF
			outGif.Image = append(outGif.Image, palettedImage)
			outGif.Delay = append(outGif.Delay, 20)
			outGif.Disposal = append(outGif.Disposal, 0x02)
		}
		if err := writeGif(filepath.Join(*out, name+".gif"), outGif); err != nil {
			panic(err)
		}
	}
}
