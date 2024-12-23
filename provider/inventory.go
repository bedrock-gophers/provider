package provider

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/playerdb"
	"github.com/df-mc/dragonfly/server/world"
)

// ConvertSavableInventory converts an InventoryData to a player.InventoryData.
func ConvertSavableInventory(data InventoryData) (playerdb.InventoryData, error) {
	var inv playerdb.InventoryData

	inv.Items = make([]item.Stack, len(data.Items))
	for i, d := range data.Items {
		stack, err := d.ToStack()
		if err != nil {
			return inv, err
		}
		inv.Items[i] = stack
	}
	inv.Boots, _ = data.Boots.ToStack()
	inv.Leggings, _ = data.Leggings.ToStack()
	inv.Chestplate, _ = data.Chestplate.ToStack()
	inv.Helmet, _ = data.Helmet.ToStack()
	inv.OffHand, _ = data.OffHand.ToStack()
	inv.MainHandSlot = data.MainHandSlot
	return inv, nil
}

// ConvertNonSavableInventory converts a player.InventoryData to an InventoryData.
func ConvertNonSavableInventory(inv playerdb.InventoryData) InventoryData {
	data := InventoryData{
		Items:        make([]StackData, len(inv.Items)),
		Boots:        StackToData(inv.Boots),
		Leggings:     StackToData(inv.Leggings),
		Chestplate:   StackToData(inv.Chestplate),
		Helmet:       StackToData(inv.Helmet),
		OffHand:      StackToData(inv.OffHand),
		MainHandSlot: inv.MainHandSlot,
	}
	for i, stack := range inv.Items {
		data.Items[i] = StackToData(stack)
	}
	return data
}

// StackToData converts an item.Stack to a StackData.
func StackToData(stack item.Stack) StackData {
	if stack.Empty() {
		return StackData{}
	}
	name, meta := stack.Item().EncodeItem()
	return StackData{
		Name:         name,
		Meta:         meta,
		Count:        stack.Count(),
		CustomName:   stack.CustomName(),
		Lore:         stack.Lore(),
		Damage:       stack.Durability(),
		AnvilCost:    stack.AnvilCost(),
		Data:         stack.Values(),
		Enchantments: EnchantmentsToData(stack.Enchantments()),
	}
}

// InventoryData represents the data of an inventory.
type InventoryData struct {
	Items        []StackData `json:",omitempty"`
	Boots        StackData   `json:",omitempty"`
	Leggings     StackData   `json:",omitempty"`
	Chestplate   StackData   `json:",omitempty"`
	Helmet       StackData   `json:",omitempty"`
	OffHand      StackData   `json:",omitempty"`
	MainHandSlot uint32      `json:",omitempty"`
}

// StackData represents the data of an item stack.
type StackData struct {
	Name  string `json:",omitempty"`
	Meta  int16  `json:",omitempty"`
	Count int    `json:",omitempty"`

	CustomName   string            `json:",omitempty"`
	Lore         []string          `json:",omitempty"`
	Damage       int               `json:",omitempty"`
	AnvilCost    int               `json:",omitempty"`
	Data         map[string]any    `json:",omitempty"`
	Enchantments []EnchantmentData `json:",omitempty"`
}

// ToStack converts the StackData to an item.Stack.
func (i StackData) ToStack() (item.Stack, error) {
	it, ok := world.ItemByName(i.Name, i.Meta)
	if !ok {
		return item.Stack{}, nil
	}
	stack := item.NewStack(it, i.Count)
	if len(i.CustomName) > 0 {
		stack = stack.WithCustomName(i.CustomName)
	}
	if len(i.Lore) > 0 {
		stack = stack.WithLore(i.Lore...)
	}
	stack = stack.WithDurability(i.Damage)
	stack = stack.WithAnvilCost(i.AnvilCost)
	for key, value := range i.Data {
		stack = stack.WithValue(key, value)
	}

	for _, ench := range i.Enchantments {
		en := ench.ToEnchantment()
		if en != nil {
			stack = stack.WithEnchantments(item.NewEnchantment(en, ench.Level))
		}
	}

	return stack, nil
}

// EnchantmentData represents the data of an enchantment.
type EnchantmentData struct {
	Name  string `json:",omitempty"`
	Level int    `json:",omitempty"`
}

// ToEnchantment converts the EnchantmentData to an item.EnchantmentType.
func (e EnchantmentData) ToEnchantment() item.EnchantmentType {
	for _, en := range item.Enchantments() {
		if en.Name() == e.Name {
			return en
		}
	}
	return nil
}

// EnchantmentsToData converts a slice of item.Enchantment to a slice of EnchantmentData.
func EnchantmentsToData(enchants []item.Enchantment) []EnchantmentData {
	data := make([]EnchantmentData, len(enchants))
	for i, ench := range enchants {
		data[i] = EnchantmentData{
			Name:  ench.Type().Name(),
			Level: ench.Level(),
		}
	}
	return data
}
