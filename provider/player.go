package provider

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/playerdb"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/google/uuid"
)

type config struct {
	player.Config
	World *world.World
}

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
	AirSupply int `json:",omitempty"`
	// MaxAirSupply is the maximum air supply the player can have.
	MaxAirSupply int `json:",omitempty"`
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
func (p *Provider) convertSavablePlayerData(dat playerData, wrld func(world.Dimension) *world.World) config {
	data := config{
		Config: player.Config{
			UUID: dat.UUID,
			Name: dat.Username,
		},
	}

	data.Position = dat.Position
	data.Velocity = dat.Velocity
	data.Rotation = cube.Rotation{dat.Yaw, dat.Pitch}
	data.Health = dat.Health
	data.MaxHealth = dat.MaxHealth
	data.FoodTick = dat.FoodTick
	data.Exhaustion = dat.ExhaustionLevel
	data.Saturation = dat.SaturationLevel
	data.EnchantmentSeed = dat.EnchantmentSeed
	data.Experience = dat.Experience
	data.AirSupply = dat.AirSupply
	data.MaxAirSupply = dat.MaxAirSupply

	gm, ok := world.GameModeByID(dat.GameMode)
	if !ok {
		p.log.Error("bedrock-gophers/provider: unknown gamemode: %d\n", dat.GameMode)
		gm = world.GameModeSurvival
	}
	data.GameMode = gm

	inv, err := ConvertSavableInventory(dat.Inventory)
	if err != nil {
		p.log.Error("bedrock-gophers/provider: error decoding inventory: %v\n", err)
	}
	data.Inventory = inventory.New(36, nil)
	for s, i := range inv.Items {
		_ = data.Inventory.SetItem(s, i)
	}

	data.EnderChestInventory = inventory.New(27, nil)
	for i, stack := range dat.EnderChestInventory {
		s, _ := stack.ToStack()
		_ = data.EnderChestInventory.SetItem(i, s)
	}
	data.FallDistance = dat.FallDistance
	data.FireTicks = dat.FireTicks

	dim, ok := world.DimensionByID(dat.Dimension)
	if ok {
		if p.settings.UseServerWorld {
			data.World = wrld(dim)
		} else {
			newWorld := p.settings.World(dim)
			if newWorld != nil {
				data.World = newWorld
			}
		}
	}

	return data
}

// convertNonSavablePlayerData converts the player data passed to a playerData struct that can be saved to disk.
func (p *Provider) convertNonSavablePlayerData(dat config) playerData {
	settings := p.settings
	playerDat := playerData{
		UUID:     dat.UUID,
		Username: dat.Name,
	}

	if settings.SaveGameMode {
		playerDat.GameMode, _ = world.GameModeID(dat.GameMode)
	}
	if settings.SaveInventory {
		invData := playerdb.InventoryData{
			Helmet:     dat.Armour.Helmet(),
			Chestplate: dat.Armour.Chestplate(),
			Leggings:   dat.Armour.Leggings(),
			Boots:      dat.Armour.Boots(),
			Items:      make([]item.Stack, len(dat.Inventory.Slots())),
		}
		for s, it := range dat.Inventory.Slots() {
			invData.Items[s] = it
		}
		playerDat.Inventory = ConvertNonSavableInventory(invData)
	}
	if settings.SaveEnderChest {
		playerDat.EnderChestInventory = make([]StackData, len(dat.EnderChestInventory.Slots()))
		for i, stack := range dat.EnderChestInventory.Slots() {
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
		playerDat.Yaw = dat.Rotation.Yaw()
		playerDat.Pitch = dat.Rotation.Pitch()
	}
	if settings.SaveHealth {
		playerDat.Health = dat.Health
		playerDat.MaxHealth = dat.MaxHealth
	}
	if settings.SaveHunger {
		playerDat.FoodTick = dat.FoodTick
		playerDat.ExhaustionLevel = dat.Exhaustion
		playerDat.SaturationLevel = dat.Saturation
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
