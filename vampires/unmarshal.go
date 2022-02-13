package vampires

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
	"strconv"
)

// createKey creates a levelDB key from the specified path.
func createKey(key string) []byte {
	return []byte("_file://\x00\x01" + key)
}

// Unmarshal reads the input struct, reads its tags and unmarshals the levelDB contents to the struct which out points
// to respectively.
func Unmarshal(db *leveldb.DB, out interface{}) error {
	v := reflect.ValueOf(out).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot unmashal, output type must be a pointer")
	}

	findPath := func(tag reflect.StructTag) (string, error) {
		if path, ok := tag.Lookup("vs_save"); ok {
			return path, nil
		}
		return "", fmt.Errorf("no vs_save tag could be found")
	}

	for i := 0; i < v.NumField(); i++ {
		typ := v.Type().Field(i)
		key, err := findPath(typ.Tag)
		if err != nil {
			continue
		}

		read, err := db.Get(createKey(key), nil)
		if err != nil {
			panic(err)
		}

		// Strip the first byte because its rubbish.
		read = read[1:]

		switch typ.Type.Kind() {
		case reflect.Slice:
			if err := unmarshalSlice(read, v.Field(i)); err != nil {
				return err
			}
		case reflect.Bool:
			if err := unmarshalBool(read, v.Field(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// unmarshalSlice reads a string slice from the save file and assigns it to a struct field.
func unmarshalSlice(data []byte, v reflect.Value) error {
	slice := new([]string)
	if err := json.Unmarshal(data, slice); err != nil {
		return err
	}
	v.Set(reflect.ValueOf(*slice))
	return nil
}

func unmarshalBool(data []byte, v reflect.Value) error {
	b, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}
