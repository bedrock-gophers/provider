package provider

import (
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

// Settings represents settings for the player provider.
type Settings struct {
	// FirstJoinMessage is the message that is written in console when a player joins for the first time.
	FirstJoinMessage string
	// Path is the path to the directory where the player data is saved.
	Path string
	// FlushRate is the rate at which the player data is flushed to from memory.
	FlushRate time.Duration
	// UseServerWorld is true if the server should use the world of the server.
	UseServerWorld bool
	// World is a function that returns the world of the specified dimension.
	World func(world.Dimension) *world.World
	// AutoSave is true if the player data should be saved automatically.
	AutoSave bool

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
		FirstJoinMessage: "[+] User with UUID %s is joining for the first time.",
		Path:             "assets/players/",
		FlushRate:        time.Minute,
		UseServerWorld:   true,
		World:            func(dimension world.Dimension) *world.World { return nil },
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
