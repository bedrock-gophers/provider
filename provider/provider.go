package provider

import (
	"errors"
	"fmt"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/google/uuid"
	"os"
	"strings"
	"sync"
	"time"
)

// Provider is a struct that is used to store player data in memory and save it to disk.
type Provider struct {
	path   string
	dataMu sync.Mutex
	data   map[uuid.UUID]player.Data

	settings Settings
	closed   bool
}

// NewProvider returns a new player provider with the settings passed.
func NewProvider(settings Settings) *Provider {
	prov := &Provider{
		path:     strings.TrimSuffix(settings.Path, "/"),
		data:     make(map[uuid.UUID]player.Data),
		settings: settings,
	}

	go prov.startCacheTicker()
	return prov
}

// Settings returns the settings of the player provider.
func (p *Provider) Settings() Settings {
	return p.settings
}

// startCacheTicker starts a ticker that flushes the player data to disk at the rate specified in the settings.
func (p *Provider) startCacheTicker() {
	ticker := time.NewTicker(p.settings.FlushRate)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if p.closed {
				return
			}

			p.dataMu.Lock()
			if !p.settings.AutoSave {
				for id, data := range p.data {
					_ = p.savePlayerData(id, data)
				}
			}

			p.data = make(map[uuid.UUID]player.Data)
			p.dataMu.Unlock()
		}
	}
}

// SavePlayer saves the player data passed to the player provider.
func (p *Provider) SavePlayer(pl *player.Player) error {
	return p.Save(pl.UUID(), pl.Data())
}

// Save saves the player data passed to the player provider.
func (p *Provider) Save(uuid uuid.UUID, data player.Data) error {
	p.dataMu.Lock()
	p.data[uuid] = data
	p.dataMu.Unlock()

	if !p.settings.AutoSave {
		return nil
	}
	return p.savePlayerData(uuid, data)
}

// savePlayerData saves the player data passed to the player provider to disk.
func (p *Provider) savePlayerData(uuid uuid.UUID, data player.Data) error {
	filePath := fmt.Sprintf("%s/%s.json", p.settings.Path, uuid.String())
	playerDat := p.convertNonSavablePlayerData(data)

	buf, err := json.MarshalIndent(playerDat, "", "\t")
	if err != nil {
		return err
	}

	_ = os.MkdirAll(p.settings.Path, os.ModePerm)
	return os.WriteFile(filePath, buf, 0644)
}

// Load loads the player data for the UUID passed from the player provider.
func (p *Provider) Load(uuid uuid.UUID, wrld func(world.Dimension) *world.World) (player.Data, error) {
	p.dataMu.Lock()
	data, ok := p.data[uuid]
	p.dataMu.Unlock()
	if ok {
		return data, nil
	}
	var playerDat playerData

	_ = os.MkdirAll(p.settings.Path, os.ModePerm)
	filePath := fmt.Sprintf("%s/%s.json", p.settings.Path, uuid.String())
	buf, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[+] User with UUID %s is joining for the first time.\n", uuid)
			return player.Data{}, errors.New("player data not found")
		}
		fmt.Println("error reading file: ", err)
		return player.Data{}, err
	}
	err = json.Unmarshal(buf, &playerDat)
	if err != nil {
		fmt.Println("error unmarshalling: ", err)
		return player.Data{}, err
	}
	dat := p.convertSavablePlayerData(playerDat, wrld)
	p.dataMu.Lock()
	p.data[uuid] = dat
	p.dataMu.Unlock()
	return dat, nil
}

// Close closes the player provider and flushes the player data to disk.
func (p *Provider) Close() error {
	p.closed = true
	p.data = nil
	return nil
}
