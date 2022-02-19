package vampires

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_MarshalSave(t *testing.T) {
	fields := []struct {
		Key             string
		SerializedValue string
		RawValue        interface{}
	}{
		{"testKey1", `"Test-Value-1"`, "Test-Value-1"},
		{"testKey2", "32", int32(32)},
		{"testKey3", "true", true},
		{"testKey4", "914", float32(914)},
		{"testKey5", `["Rainer","Wingl"]`, []string{"Rainer", "Wingl"}},
		{"testKey6", `{"Iam":91,"JustA":44,"Test":8}`, map[string]int32{
			"Iam":   91,
			"JustA": 44,
			"Test":  8,
		}},
	}

	dummySaveFile := &struct {
		TestKey1 string           `vs_save:"testKey1"`
		TestKey2 int32            `vs_save:"testKey2"`
		TestKey3 bool             `vs_save:"testKey3"`
		TestKey4 float32          `vs_save:"testKey4"`
		TestKey5 []string         `vs_save:"testKey5"`
		TestKey6 map[string]int32 `vs_save:"testKey6"`
	}{}

	expected := &SerializedSaveFile{}
	for i, test := range fields {
		setStructField(dummySaveFile, i, test.RawValue)
		expected.Entries = append(expected.Entries, createSerializedSaveFileEntry(test.Key, test.SerializedValue))
	}

	actual, err := MarshalSave(dummySaveFile)
	assert.NoError(t, err)

	for _, entry := range expected.Entries {
		assert.Contains(t, actual.Entries, entry)
	}
}

func setStructField(target interface{}, index int, value interface{}) {
	reflect.ValueOf(target).Elem().Field(index).Set(reflect.ValueOf(value))
}

func createSerializedSaveFileEntry(key, value string) SerializedSaveFileEntry {
	return SerializedSaveFileEntry{Key: createKey(key), Value: createValue([]byte(value))}
}
