package vampires

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// SaveFile wraps the contents of a Vampire Survivors save file.
type SaveFile struct {
	Achievements   []string `vs_save:"CapacitorStorage.Achievements"`
	UnlockedHypers []string `vs_save:"CapacitorStorage.UnlockedHypers"`
}

// ParseSave reads the Vampire Survivors save file located at the specified path or an error on failure.
func ParseSave(path string) (*SaveFile, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var save SaveFile
	if err := Unmarshal(db, &save); err != nil {
		return nil, err
	}
	return &save, nil
}
