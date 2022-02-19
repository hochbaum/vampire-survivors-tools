package vampires

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_createKey(t *testing.T) {
	data := []struct {
		Input  string
		Output []byte
	}{
		{"CapacitorStorage.CheatCodeUsed", []byte("_file://\x00\x01CapacitorStorage.CheatCodeUsed")},
		{"CapacitorStorage.Achievements", []byte("_file://\x00\x01CapacitorStorage.Achievements")},
		{"CapacitorStorage.BLuck", []byte("_file://\x00\x01CapacitorStorage.BLuck")},
	}
	for _, entry := range data {
		key := createKey(entry.Input)
		assert.ElementsMatch(t, key, entry.Output,
			"format of actual key should be %q but is %q", string(entry.Output), string(key))
	}
}

func Test_scanStructKeys(t *testing.T) {
	s, expected := createTestStructForScanning()
	actual, err := scanStructTags(s, "test")

	assert.NoError(t, err)

	for k, v := range expected {
		assert.Contains(t, actual, k, "actual map does not contain expected key %q", k)
		assert.Equal(t, v, actual[k], "actual map value %q is not equal to expected value: %v", actual[k])
	}

	assert.Equal(t, expected, actual)
}

func createTestStructForScanning() (interface{}, map[string]reflect.Value) {
	s := &struct {
		Field1 string   `test:"field1"`
		Field2 int32    `test:"field2"`
		Field3 []string `test:"field3"`
	}{}
	vals := map[string]reflect.Value{
		"field1": reflect.ValueOf(s).Elem().Field(0),
		"field2": reflect.ValueOf(s).Elem().Field(1),
		"field3": reflect.ValueOf(s).Elem().Field(2),
	}
	return s, vals
}
