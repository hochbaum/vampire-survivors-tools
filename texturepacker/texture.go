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
	Size   Dimension `json:"size"`
	Scale  int       `json:"scale"`
	Frames []Frame   `json:"frames"`
}

// Frame defines a frame entry of a PackedTexture.
type Frame struct {
	FileName         string
	Rotated          bool
	Trimmed          bool
	SourceSize       Dimension
	SpriteSourceSize image.Rectangle
	Frame            image.Rectangle
}

type Dimension struct {
	Width  int `json:"w"`
	Height int `json:"h"`
}

func (f *Frame) UnmarshalJSON(data []byte) error {
	values := make(map[string]interface{})
	if err := json.Unmarshal(data, &values); err != nil {
		return err
	}
	f.FileName = values["filename"].(string)
	f.Rotated = values["rotated"].(bool)
	f.Trimmed = values["trimmed"].(bool)
	f.SourceSize = dimension(values["sourceSize"].(map[string]interface{}))
	f.SpriteSourceSize = imageRect(values["spriteSourceSize"].(map[string]interface{}))
	f.Frame = imageRect(values["frame"].(map[string]interface{}))
	return nil
}

// dimension creates a Dimension instance from the provided JSON map.
func dimension(data map[string]interface{}) Dimension {
	return Dimension{Width: int(data["w"].(float64)), Height: int(data["h"].(float64))}
}

// imageRect creates an instance of image.Reactangle from the provided JSON map.
func imageRect(data map[string]interface{}) image.Rectangle {
	x := data["x"].(float64)
	y := data["y"].(float64)
	w := data["w"].(float64)
	h := data["h"].(float64)
	return image.Rect(int(x), int(y), int(x+w), int(y+h))
}
