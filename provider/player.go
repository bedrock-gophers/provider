package provider

import (
	"fmt"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

// playerData is a struct that is used to store the data of a player in a format that can be saved to a file.
type playerData struct {
	// UUID is the player's unique identifier for their account.
	UUID uuid.UUID
	// Username is the last username the player joined with.
	Username string
	// Position is the last position the player was located at.
	// Velocity is the speed at which the player was moving.
	Position, Velocity mgl64.Vec3 `json:",omitempty"`
	// Yaw and Pitch represent the rotation of the player.
	Yaw, Pitch float64 `json:",omitempty"`
	// Health, MaxHealth ...
	Health, MaxHealth float64 `json:",omitempty"`
	// Hunger is the amount of hunger points the player currently has, shown on the hunger bar.
	// This should be between 0-20.
	Hunger int `json:",omitempty"`
	// FoodTick this variable is used when the hunger exceeds 17 or is equal to 0. It is used to heal
	// the player using saturation or make the player starve and receive damage if the hunger is at 0.
	// This value should be between 0-80.
	FoodTick int `json:",omitempty"`
	// ExhaustionLevel determines how fast the hunger level depletes and is controlled by the kinds
	// of food the player has eaten. SaturationLevel determines how fast the saturation level depletes.
	ExhaustionLevel, SaturationLevel float64 `json:",omitempty"`
	// AbsorptionLevel represents the amount of extra health the player has. This health does not regenerate
	// and is lost when the player takes damage.
	AbsorptionLevel float64 `json:",omitempty"`
	// EnchantmentSeed is the seed used to generate the enchantments from enchantment tables.
	EnchantmentSeed int64 `json:",omitempty"`
	// Experience is the current experience the player has.
	Experience int `json:",omitempty"`
	// AirSupply is the current tick of the player's air supply.
	AirSupply int64 `json:",omitempty"`
	// MaxAirSupply is the maximum air supply the player can have.
	MaxAirSupply int64 `json:",omitempty"`
	// GameMode is the last gamemode the playerData had, like creative or survival.
	GameMode int `json:",omitempty"`
	// Inventory contains all the items in the inventory, including armor, main inventory and offhand.
	Inventory InventoryData `json:",omitempty"`
	// EnderChestInventory contains the items in the player's ender chest.
	EnderChestInventory []StackData `json:",omitempty"`
	// Effects contains all the currently active potions effects the player has.
	Effects []EffectData `json:",omitempty"`
	// FireTicks is the amount of ticks the player will be on fire for.
	FireTicks int64 `json:",omitempty"`
	// FallDistance is the distance the player has currently been falling. This is used to calculate fall damage.
	FallDistance float64 `json:",omitempty"`
	// Dimension is the dimension the player is in.
	Dimension int `json:",omitempty"`
}

// convertSavablePlayerData converts the player data passed to a playerData struct that can be saved to disk.
func (p *Provider) convertSavablePlayerData(dat playerData, wrld func(world.Dimension) *world.World) player.Data {
	data := player.Data{
		UUID:     dat.UUID,
		Username: dat.Username,
	}
	data.Position = dat.Position
	data.Velocity = dat.Velocity
	data.Yaw = dat.Yaw
	data.Pitch = dat.Pitch
	data.Health = dat.Health
	data.MaxHealth = dat.MaxHealth
	data.Hunger = dat.Hunger
	data.FoodTick = dat.FoodTick
	data.ExhaustionLevel = dat.ExhaustionLevel
	data.SaturationLevel = dat.SaturationLevel
	data.AbsorptionLevel = dat.AbsorptionLevel
	data.EnchantmentSeed = dat.EnchantmentSeed
	data.Experience = dat.Experience
	data.AirSupply = dat.AirSupply
	data.MaxAirSupply = dat.MaxAirSupply

	gm, ok := world.GameModeByID(dat.GameMode)
	if !ok {
		fmt.Printf("unknown gamemode: %d\n", dat.GameMode)
		gm = world.GameModeSurvival
	}
	data.GameMode = gm

	inv, err := ConvertSavableInventory(dat.Inventory)
	if err != nil {
		fmt.Printf("error decoding inventory: %v\n", err)
	}
	data.Inventory = inv
	data.EnderChestInventory = make([]item.Stack, len(dat.EnderChestInventory))
	for i, stack := range dat.EnderChestInventory {
		data.EnderChestInventory[i], _ = stack.ToStack()
	}
	data.FallDistance = dat.FallDistance
	data.FireTicks = dat.FireTicks

	dim, ok := world.DimensionByID(dat.Dimension)
	if ok {
		data.World = wrld(dim)
		if data.World == nil {
			newWorld := p.settings.World(dim)
			if newWorld != nil {
				data.World = newWorld
			}
		}
	}

	return data
}

// convertNonSavablePlayerData converts the player data passed to a playerData struct that can be saved to disk.
func (p *Provider) convertNonSavablePlayerData(dat player.Data) playerData {
	settings := p.settings
	playerDat := playerData{
		UUID:     dat.UUID,
		Username: dat.Username,
	}

	if settings.SaveGameMode {
		playerDat.GameMode, _ = world.GameModeID(dat.GameMode)
	}
	if settings.SaveInventory {
		playerDat.Inventory = ConvertNonSavableInventory(dat.Inventory)
	}
	if settings.SaveEnderChest {
		playerDat.EnderChestInventory = make([]StackData, len(dat.EnderChestInventory))
		for i, stack := range dat.EnderChestInventory {
			playerDat.EnderChestInventory[i] = StackToData(stack)
		}
	}
	if settings.SavePosition {
		playerDat.Position = dat.Position
	}
	if settings.SaveVelocity {
		playerDat.Velocity = dat.Velocity
	}
	if settings.SaveRotation {
		playerDat.Yaw = dat.Yaw
		playerDat.Pitch = dat.Pitch
	}
	if settings.SaveHealth {
		playerDat.Health = dat.Health
		playerDat.MaxHealth = dat.MaxHealth
	}
	if settings.SaveHunger {
		playerDat.Hunger = dat.Hunger
		playerDat.FoodTick = dat.FoodTick
		playerDat.ExhaustionLevel = dat.ExhaustionLevel
		playerDat.SaturationLevel = dat.SaturationLevel
	}
	if settings.SaveAbsorption {
		playerDat.AbsorptionLevel = dat.AbsorptionLevel
	}
	if settings.SaveEnchantment {
		playerDat.EnchantmentSeed = dat.EnchantmentSeed
	}
	if settings.SaveExperience {
		playerDat.Experience = dat.Experience
	}
	if settings.SaveAirSupply {
		playerDat.AirSupply = dat.AirSupply
		playerDat.MaxAirSupply = dat.MaxAirSupply
	}
	if settings.SaveFallDistance {
		playerDat.FallDistance = dat.FallDistance
	}
	if settings.SaveFireTicks {
		playerDat.FireTicks = dat.FireTicks
	}
	if settings.SaveDimension {
		playerDat.Dimension, _ = world.DimensionID(dat.World.Dimension())
	}

	return playerDat
}
