package vampires

import (
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"reflect"
	"testing"
)

type mockSaveStorage struct {
	data map[string][]byte
	t    *testing.T
}

func (m *mockSaveStorage) Get(key []byte, _ *opt.ReadOptions) ([]byte, error) {
	stringKey := string(key)
	assert.Contains(m.t, m.data, stringKey)
	return m.data[stringKey], nil
}

func (m *mockSaveStorage) Put([]byte, []byte, *opt.WriteOptions) error {
	m.t.Errorf("Put(...) should not have been called by a read operation")
	return nil
}

func Test_UnmarshalSave(t *testing.T) {
	fields := []struct {
		Key    string
		Input  string
		Output interface{}
	}{
		{"testKey1", `"Test-Value-1"`, "Test-Value-1"},
		{"testKey2", "32", int32(32)},
		{"testKey3", "true", true},
		{"testKey4", "914.48", float32(914.48)},
		{"testKey5", `["Rainer", "Wingl"]`, []string{"Rainer", "Wingl"}},
		{"testKey6", `{"Iam": 91, "JustA": 44, "Test": 8}`, map[string]int32{
			"Iam":   91,
			"JustA": 44,
			"Test":  8,
		}},
	}

	data := make(map[string][]byte)
	for _, test := range fields {
		data[string(createKey(test.Key))] = createValue([]byte(test.Input))
	}

	db := &mockSaveStorage{data: data, t: t}
	s := struct {
		TestKey1 string           `vs_save:"testKey1"`
		TestKey2 int32            `vs_save:"testKey2"`
		TestKey3 bool             `vs_save:"testKey3"`
		TestKey4 float32          `vs_save:"testKey4"`
		TestKey5 []string         `vs_save:"testKey5"`
		TestKey6 map[string]int32 `vs_save:"testKey6"`
	}{}
	assert.NoError(t, UnmarshalSave(db, &s))

	elem := reflect.ValueOf(&s).Elem()
	for i := 0; i < elem.NumField(); i++ {
		assert.Equal(t, elem.Field(i).Interface(), fields[i].Output)
	}
}
