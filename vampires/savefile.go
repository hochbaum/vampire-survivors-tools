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

// ParseSave reads the Vampire Survivors save file located at the specified path or an error on failure.
func ParseSave(path string) (*SaveFile, error) {
	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	var save SaveFile
	return &save, Unmarshal(db, &save)
}
