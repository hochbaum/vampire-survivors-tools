package texturepacker

import (
	"encoding/json"
	"image"
	"io"
	"os"
)

// Open opens and parses the texturepacker-packed sprite sheet located at the specified path.
func Open(path string) (*Sheet, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return unmarshalPackedTextures(file)
}

// unmarshalPackedTextures parses a texturepacker-packed sprite sheet from the provided reader.
func unmarshalPackedTextures(r io.Reader) (sheet *Sheet, err error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	sheet = new(Sheet)
	return sheet, json.Unmarshal(data, sheet)
}

// Sheet defines a sprite sheet of the Vampire Survivors game, packed using texturepacker.
type Sheet struct {
	Textures []PackedTexture `json:"textures"`
	Metadata struct {
		App         string `json:"app"`
		Version     string `json:"version"`
		SmartUpdate string `json:"smartupdate"`
	} `json:"meta"`
}

// PackedTexture defines a texture packed by texturepacker.
type PackedTexture struct {
	Image  string    `json:"image"`
	Format string    `json:"format"`
	Size   jsonPoint `json:"size"`
	Scale  int       `json:"scale"`
	Frames []Frame   `json:"frames"`
}

// Frame defines a frame entry of a PackedTexture.
type Frame struct {
	FileName         string        `json:"filename"`
	Rotated          bool          `json:"rotated"`
	Trimmed          bool          `json:"trimmed"`
	SourceSize       jsonPoint     `json:"sourceSize"`
	SpriteSourceSize jsonRectangle `json:"spriteSourceSize"`
	Frame            jsonRectangle `json:"frame"`
}

// jsonPoint defines the JSON representation of an image.Point.
type jsonPoint struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// Point returns the jsonPoint as image.Point.
func (p jsonPoint) Point() image.Point {
	return image.Pt(p.X, p.Y)
}

// jsonRectangle defines the JSON representation of an image.Rectangle.
type jsonRectangle struct {
	jsonPoint
	Width  int `json:"w"`
	Height int `json:"h"`
}

// Rect returns the jsonRectangle as image.Rectangle.
func (r jsonRectangle) Rect() image.Rectangle {
	return image.Rect(r.X, r.Y, r.X+r.Width, r.Y+r.Height)
}
