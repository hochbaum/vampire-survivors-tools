package texturepacker

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_unmarshalPackedTextures(t *testing.T) {
	expected, data := createTestSheet()
	r := strings.NewReader(string(data))
	actual, err := unmarshalPackedTextures(r)
	assert.NoError(t, err)
	assert.Equal(t, expected, *actual)
}

func createTestSheet() (Sheet, []byte) {
	sheet := Sheet{
		Textures: []PackedTexture{
			{
				Image:  "image01.png",
				Format: "RGBA8888",
				Size:   jsonDimension{13, 37},
				Scale:  1,
				Frames: []Frame{
					{
						FileName:   "frame01.png",
						Rotated:    true,
						Trimmed:    false,
						SourceSize: jsonDimension{900, 1},
						SpriteSourceSize: jsonRectangle{
							jsonDimension: jsonDimension{5, 6},
							X:             42,
							Y:             360,
						},
						Frame: jsonRectangle{
							jsonDimension: jsonDimension{90, 12},
							X:             413,
							Y:             600,
						},
					},
				},
			},
		},
		Metadata: struct {
			App         string `json:"app"`
			Version     string `json:"version"`
			SmartUpdate string `json:"smartupdate"`
		}{
			App:         "dummyApp",
			Version:     "dummyVersion",
			SmartUpdate: "dummyUpdate",
		},
	}
	data, _ := json.Marshal(&sheet)
	return sheet, data
}
