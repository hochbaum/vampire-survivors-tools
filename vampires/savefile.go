package vampires

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// SaveFile wraps the contents of a Vampire Survivors save file.
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
