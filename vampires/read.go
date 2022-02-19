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

// unmarshalers maps types to their respective unmarshalFunc.
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

// UnmarshalSave reads the entries from the SaveStorage and unmarshalls them into the fields tagged with `vs_save` in the
// provided interface.
//
// If a referenced LevelDB key could not be found in the database, this function does not return an error but prints a
// warning, as new save files don't contain every possible key.
func UnmarshalSave(db SaveStorage, i interface{}) error {
	taggedFields, err := scanStructTags(i, "vs_save")
	if err != nil {
		return err
	}

	for key, value := range taggedFields {
		data, err := db.Get(createKey(key), nil)
		if err == leveldb.ErrNotFound {
			fmt.Printf("warning: ignoring field tagged with %s at it is not present in the levelDB\n", key)
			continue
		} else if err != nil {
			return err
		}

		data = data[1:]
		fieldType := value.Type().Kind()

		unmarshaler, ok := unmarshalers[fieldType]
		if !ok {
			return fmt.Errorf("could not find suitable unmarshaler for type %s", fieldType)
		}

		if err := unmarshaler(data, value); err != nil {
			return err
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
