package provider

import "time"

// Settings represents settings for the player provider.
type Settings struct {
	// Path is the path to the directory where the player data is saved.
	Path string
	// FlushRate is the rate at which the player data is flushed to from memory.
	FlushRate time.Duration
	// AutoSave is true if the player data should be saved automatically.
	AutoSave         bool
	SavePosition     bool
	SaveVelocity     bool
	SaveRotation     bool
	SaveHealth       bool
	SaveHunger       bool
	SaveFoodTick     bool
	SaveExhaustion   bool
	SaveSaturation   bool
	SaveAbsorption   bool
	SaveEnchantment  bool
	SaveExperience   bool
	SaveGameMode     bool
	SaveInventory    bool
	SaveEffects      bool
	SaveEnderChest   bool
	SaveAirSupply    bool
	SaveFallDistance bool
	SaveFireTicks    bool
	SaveDimension    bool
}

func DefaultSettings() Settings {
	return Settings{
		Path:             "assets/players/",
		FlushRate:        time.Minute,
		AutoSave:         true,
		SavePosition:     true,
		SaveVelocity:     true,
		SaveRotation:     true,
		SaveHealth:       true,
		SaveHunger:       true,
		SaveFoodTick:     true,
		SaveExhaustion:   true,
		SaveSaturation:   true,
		SaveAbsorption:   true,
		SaveEnchantment:  true,
		SaveExperience:   true,
		SaveGameMode:     true,
		SaveInventory:    true,
		SaveEffects:      true,
		SaveEnderChest:   true,
		SaveAirSupply:    true,
		SaveFallDistance: true,
		SaveFireTicks:    true,
		SaveDimension:    true,
	}
}
