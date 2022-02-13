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

type unmarshalFunc func(data []byte, v reflect.Value) error

var unmarshalers = map[reflect.Kind]unmarshalFunc{
	reflect.Slice:  unmarshalSlice,
	reflect.Bool:   unmarshalBool,
	reflect.String: unmarshalString,
}

// Unmarshal reads the input struct, reads its tags and unmarshalls the levelDB contents to the struct which `out` points
// to respectively.
func Unmarshal(db *leveldb.DB, out interface{}) error {
	v := reflect.ValueOf(out).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot unmarshal, output type must be a pointer")
	}

	for i := 0; i < v.NumField(); i++ {
		typ := v.Type().Field(i)
		key, ok := typ.Tag.Lookup("vs_save")
		if !ok {
			continue
		}

		read, err := db.Get(createKey(key), nil)
		if err != nil {
			panic(err)
		}

		// Strip the first byte because its rubbish.
		read = read[1:]
		kind := typ.Type.Kind()

		if unmarshaler, found := unmarshalers[kind]; found {
			if err := unmarshaler(read, v.Field(i)); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown kind %q", kind)
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

func unmarshalString(data []byte, v reflect.Value) error {
	v.SetString(string(data))
	return nil
}
