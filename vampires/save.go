package vampires

import (
	"encoding/json"
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"reflect"
)

// marshalFunc defines a function which serializes the content of the provided reflect.Value into a byte array.
type marshalFunc func(v reflect.Value) ([]byte, error)

// marshalers maps primitive types to their respective marshalFunc.
var marshalers = map[reflect.Kind]marshalFunc{
	reflect.Bool:    marshalBool,
	reflect.String:  marshalString,
	reflect.Float32: marshalFloat,
	reflect.Float64: marshalFloat,
	reflect.Int32:   marshalInt,
	reflect.Int64:   marshalInt,
	reflect.Slice:   marshalStringSlice,
	reflect.Map:     marshalStringToIntMap,
}

// saveSaveFile serializes the provided SaveFile and writes it to the LevelDB instance.
func saveSaveFile(save *SaveFile, db *leveldb.DB) error {
	v := reflect.ValueOf(save).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		key, ok := field.Tag.Lookup("vs_save")
		if !ok {
			continue
		}

		kind := field.Type.Kind()
		if marshaler, found := marshalers[kind]; found {
			bytes, err := marshaler(v.Field(i))
			if err != nil {
				return err
			}
			// Add a dummy byte in front of the value.
			bytes = append([]byte{'\x01'}, bytes...)
			if err := db.Put(createKey(key), bytes, nil); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unknown kind %q", kind)
		}
	}
	return nil
}

// marshalBool serializes the bool in `v`.
func marshalBool(v reflect.Value) ([]byte, error) {
	return json.Marshal(v.Bool())
}

// marshalString serializes the string in `v`.
func marshalString(v reflect.Value) ([]byte, error) {
	return json.Marshal(v.String())
}

// marshalFloat serializes the float in `v`.
func marshalFloat(v reflect.Value) ([]byte, error) {
	return json.Marshal(v.Float())
}

// marshalInt serializes the int in `v`.
func marshalInt(v reflect.Value) ([]byte, error) {
	return json.Marshal(v.Int())
}

// marshalStringSlice serializes the string slice in `v`.
func marshalStringSlice(v reflect.Value) ([]byte, error) {
	slice, ok := v.Interface().([]string)
	if !ok {
		return nil, fmt.Errorf("could not marshal string slice")
	}
	return json.Marshal(slice)
}

// marshalStringToIntMap serializes the map in `v`.
func marshalStringToIntMap(v reflect.Value) ([]byte, error) {
	m, ok := v.Interface().(map[string]int32)
	if !ok {
		return nil, fmt.Errorf("could not marshal string-to-int-map")
	}
	return json.Marshal(m)
}
