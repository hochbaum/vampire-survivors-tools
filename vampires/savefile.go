package vampires

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"reflect"
)

// SaveFile wraps the contents of a Vampire Survivors save file.
//
// Vampire Survivors uses LevelDB for storing save files, it's located at `%APPDATA%/Vampire Survivors/Local Storage`.
// The LevelDB keys are prefixed with `_file://` followed by a 0-byte and a 1-byte.
// The LevelDB values are prefixed with a 1-byte.
type SaveFile struct {
	Achievements         []string `vs_save:"CapacitorStorage.Achievements"`
	BoughtCharacters     []string `vs_save:"CapacitorStorage.BoughtCharacters"`
	BoughtPowerups       []string `vs_save:"CapacitorStorage.BoughtPowerups"`
	CollectedItems       []string `vs_save:"CapacitorStorage.CollectedItems"`
	CollectedWeapons     []string `vs_save:"CapacitorStorage.CollectedWeapons"`
	UnlockedCharacters   []string `vs_save:"CapacitorStorage.UnlockedCharacters"`
	UnlockedHypers       []string `vs_save:"CapacitorStorage.UnlockedHypers"`
	UnlockedPowerUpRanks []string `vs_save:"CapacitorStorage.UnlockedPowerUpRanks"`
	UnlockedStages       []string `vs_save:"CapacitorStorage.UnlockedStages"`
	UnlockedWeapons      []string `vs_save:"CapacitorStorage.UnlockedWeapons"`

	CheatCodeUsed        bool `vs_save:"CapacitorStorage.CheatCodeUsed"`
	DamageNumbersEnabled bool `vs_save:"CapacitorStorage.DamageNumbersEnabled"`
	FlashingVfxEnabled   bool `vs_save:"CapacitorStorage.FlashingVFXEnabled"`
	JoystickVisible      bool `vs_save:"CapacitorStorage.JoystickVisible"`
	SelectedHyper        bool `vs_save:"CapacitorStorage.SelectedHyper"`
	StreamSafeEnabled    bool `vs_save:"CapacitorStorage.StreamSafeEnabled"`

	Language          string `vs_save:"CapacitorStorage.Language"`
	SelectedCharacter string `vs_save:"CapacitorStorage.SelectedCharacter"`
	SelectedStage     string `vs_save:"CapacitorStorage.SelectedStage"`

	Coins         float64 `vs_save:"CapacitorStorage.Coins"`
	LifetimeCoins float64 `vs_save:"CapacitorStorage.LifetimeCoins"`
	LifetimeHeal  float64 `vs_save:"CapacitorStorage.LifetimeHeal"`
	MusicVolume   float64 `vs_save:"CapacitorStorage.MusicVolume"`
	SoundsVolume  float64 `vs_save:"CapacitorStorage.SoundsVolume"`

	BLuck            int32 `vs_save:"CapacitorStorage.BLuck"`
	LifetimeSurvived int32 `vs_save:"CapacitorStorage.LifetimeSurvived"`

	DestroyedCount map[string]int32 `vs_save:"CapacitorStorage.DestroyedCount"`
	KillCount      map[string]int32 `vs_save:"CapacitorStorage.KillCount"`
	PickupCount    map[string]int32 `vs_save:"CapacitorStorage.PickupCount"`
}

// SerializedSaveFileEntry defines a serialized entry of a Vampire Survivors save file. Its fields follow the rules
// explained in the SaveFile doc.
type SerializedSaveFileEntry struct {
	Key, Value []byte
}

// SerializedSaveFile defines the save file itself, containing multiple SerializedSaveFileEntry s.
type SerializedSaveFile struct {
	Entries []SerializedSaveFileEntry
}

// OpenSaveFile opens a Vampire Survivors save file located at the provided path and returns it wrapped in a SaveFile
// instance, as well as the LevelDB itself, which must be closed by the user.
func OpenSaveFile(path string) (*SaveFile, *leveldb.DB, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, nil, err
	}
	save := new(SaveFile)
	return save, db, UnmarshalSave(db, save)
}

// StoreSaveFile writes the SaveFile to the provided LevelDB, which you can obtain by using OpenSaveFile.
func StoreSaveFile(save *SaveFile, db SaveStorage) error {
	serialized, err := MarshalSave(save)
	if err != nil {
		return err
	}
	return writeSaveToDB(serialized, db)
}

// scanStructTags collects the reflect.Value s tagged by the provided tag in the provided struct. It maps the key of the
// struct tag to its respective reflect.Value.
func scanStructTags(i interface{}, tag string) (map[string]reflect.Value, error) {
	elem := reflect.ValueOf(i).Elem()
	if !elem.CanAddr() {
		return nil, fmt.Errorf("input type must be a pointer type")
	}

	result := make(map[string]reflect.Value)
	for i := 0; i < elem.NumField(); i++ {
		field := elem.Type().Field(i)
		value := elem.Field(i)

		key, ok := field.Tag.Lookup(tag)
		if !ok {
			continue
		}
		result[key] = value
	}

	return result, nil
}

// SaveStorage defines a wrapper for leveldb.DB.
type SaveStorage interface {
	Get(key []byte, ro *opt.ReadOptions) ([]byte, error)
	Put(key []byte, value []byte, wo *opt.WriteOptions) error
}

// createKey formats and serializes the provided string to be a valid LevelDB key.
func createKey(key string) []byte {
	return []byte("_file://\x00\x01" + key)
}

// createValue formats the provided bytes to be a valid LevelDB value.
func createValue(value []byte) []byte {
	return append([]byte{'\x01'}, value...)
}
