package vampires

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
)

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

// readSaveFile reads the raw save file from the LevelDB and returns its contents wrapped inside an SaveFile instance.
func readSaveFile(db *leveldb.DB) (*SaveFile, error) {
	save := new(SaveFile)
	v := reflect.ValueOf(save).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		key, ok := field.Tag.Lookup("vs_save")
		if !ok {
			continue
		}

		read, err := db.Get(createKey(key), nil)
		if err != nil {
			if err == leveldb.ErrNotFound {
				continue
			}
			panic(err)
		}

		// Strip the \x01 byte.
		read = read[1:]
		kind := field.Type.Kind()

		if unmarshaler, found := unmarshalers[kind]; found {
			if err := unmarshaler(read, v.Field(i)); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("unknown kind %q", kind)
		}
	}
	return save, nil
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
	var b bool
	if err := json.Unmarshal(data, &b); err != nil {
		return err
	}
	v.SetBool(b)
	return nil
}

// unmarshalString reads a string from the save file and assigns it to a struct field.
func unmarshalString(data []byte, v reflect.Value) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	v.SetString(str)
	return nil
}

// unmarshalFloat reads a float from the save file and assigns it to a struct field.
func unmarshalFloat(data []byte, v reflect.Value) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	v.SetFloat(f)
	return nil
}

// unmarshalInt reads an int from the save file and assigns it to a struct field.
func unmarshalInt(data []byte, v reflect.Value) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
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
