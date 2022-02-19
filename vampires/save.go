package vampires

import (
	"encoding/json"
	"fmt"
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

// MarshalSave serializes save file wrapper provided and returns a SerializedSaveFile handle.
func MarshalSave(i interface{}) (*SerializedSaveFile, error) {
	serialized := new(SerializedSaveFile)
	taggedFields, err := scanStructTags(i, "vs_save")
	if err != nil {
		return nil, err
	}

	for key, value := range taggedFields {
		fieldType := value.Type().Kind()
		marshaler, ok := marshalers[fieldType]
		if !ok {
			return nil, fmt.Errorf("could not find suitable marshaler for type %s", fieldType)
		}

		data, err := marshaler(value)
		if err != nil {
			return nil, err
		}
		serialized.Entries = append(serialized.Entries, SerializedSaveFileEntry{createKey(key), createValue(data)})
	}

	return serialized, nil
}

// writeSaveToDB writes a SerializedSaveFile to the provided LevelDB.
func writeSaveToDB(serialized *SerializedSaveFile, db SaveStorage) error {
	for _, entry := range serialized.Entries {
		if err := db.Put(entry.Key, entry.Value, nil); err != nil {
			return err
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
		return nil, fmt.Errorf("could not MarshalSave string slice")
	}
	return json.Marshal(slice)
}

// marshalStringToIntMap serializes the map in `v`.
func marshalStringToIntMap(v reflect.Value) ([]byte, error) {
	m, ok := v.Interface().(map[string]int32)
	if !ok {
		return nil, fmt.Errorf("could not MarshalSave string-to-int-map")
	}
	return json.Marshal(m)
}
