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

// unmarshalFunc defines a function which takes data read from the LevelDB, converts it and unmarshalls it into the
// provided reflect.Value.
type unmarshalFunc func(data []byte, v reflect.Value) error

// unmarshalers maps primitive types to their respective unmarshalFunc.
var unmarshalers = map[reflect.Kind]unmarshalFunc{
	reflect.Bool:    unmarshalBool,
	reflect.String:  unmarshalString,
	reflect.Float32: unmarshalFloat,
	reflect.Float64: unmarshalFloat,
	reflect.Int32:   unmarshalInt,
	reflect.Int64:   unmarshalInt,

	// TODO: Make this a bit more generic.
	// String slices and string to int maps are the only types of collections needed right now but in case we get more
	// collection types in a future update, this should be changed.
	reflect.Slice: unmarshalStringSlice,
	reflect.Map:   unmarshalStringToIntMap,
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

// unmarshalStringSlice reads a string slice from the save file and assigns it to a struct field.
func unmarshalStringSlice(data []byte, v reflect.Value) error {
	slice := new([]string)
	if err := json.Unmarshal(data, slice); err != nil {
		return err
	}
	v.Set(reflect.ValueOf(*slice))
	return nil
}

// unmarshalBool reads a bool from the save file and assigns it to a struct field.
func unmarshalBool(data []byte, v reflect.Value) error {
	b, err := strconv.ParseBool(string(data))
	if err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

// unmarshalString reads a string from the save file and assigns it to a struct field.
func unmarshalString(data []byte, v reflect.Value) error {
	v.SetString(string(data))
	return nil
}

// unmarshalFloat reads a float from the save file and assigns it to a struct field.
func unmarshalFloat(data []byte, v reflect.Value) error {
	f, err := strconv.ParseFloat(string(data), 64)
	if err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

// unmarshalInt reads an int from the save file and assigns it to a struct field.
func unmarshalInt(data []byte, v reflect.Value) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	v.SetInt(i)
	return nil
}

// unmarshalStringToIntMap reads a map with string-keys and int-values from the save file and assigns it to a struct field.
func unmarshalStringToIntMap(data []byte, v reflect.Value) error {
	m := new(map[string]int32)
	if err := json.Unmarshal(data, m); err != nil {
		return err
	}
	v.Set(reflect.ValueOf(*m))
	return nil
}
