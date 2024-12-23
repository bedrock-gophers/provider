package provider

import (
	"bytes"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-jose/go-jose/v3/json"
	"github.com/google/uuid"
)

// Provider is a struct that is used to store player data in memory and save it to disk.
type Provider struct {
	log *slog.Logger

	path   string
	dataMu sync.Mutex
	data   map[uuid.UUID]config

	settings Settings
	closed   bool
}

// NewProvider returns a new player provider with the settings passed.
func NewProvider(conf *server.Config, settings Settings) *Provider {
	prov := &Provider{
		log:      conf.Log,
		path:     strings.TrimSuffix(settings.Path, "/"),
		data:     make(map[uuid.UUID]config),
		settings: settings,
	}
	conf.PlayerProvider = prov

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

			p.data = make(map[uuid.UUID]config)
			p.dataMu.Unlock()
		}
	}
}

// SavePlayer saves the player data passed to the player provider.
func (p *Provider) SavePlayer(pl *player.Player) error {
	return p.Save(pl.UUID(), pl.Data(), pl.Tx().World())
}

// Save saves the player data passed to the player provider.
func (p *Provider) Save(uuid uuid.UUID, cfg player.Config, w *world.World) error {
	data := config{Config: cfg, World: w}

	p.dataMu.Lock()
	p.data[uuid] = data
	p.dataMu.Unlock()

	if !p.settings.AutoSave {
		return nil
	}
	return p.savePlayerData(uuid, data)
}

// filePath returns the file path for the UUID passed.
func (p *Provider) filePath(uuid uuid.UUID) string {
	return p.settings.Path + "/" + uuid.String() + ".json"
}

// savePlayerData saves the player data passed to the player provider to disk.
func (p *Provider) savePlayerData(uuid uuid.UUID, data config) error {
	playerDat := p.convertNonSavablePlayerData(data)

	buf, err := json.MarshalIndent(playerDat, "", "\t")
	if err != nil {
		return err
	}

	_ = os.MkdirAll(p.settings.Path, os.ModePerm)
	return os.WriteFile(p.filePath(uuid), buf, 0644)
}

// Load loads the player data for the UUID passed from the player provider.
func (p *Provider) Load(uuid uuid.UUID, wrld func(world.Dimension) *world.World) (player.Config, *world.World, error) {
	p.dataMu.Lock()
	data, ok := p.data[uuid]
	p.dataMu.Unlock()
	if ok {
		return data.Config, data.World, nil
	}
	var playerDat playerData

	_ = os.MkdirAll(p.settings.Path, os.ModePerm)
	buf, err := os.ReadFile(p.filePath(uuid))
	if err != nil {
		if os.IsNotExist(err) {
			msg := p.settings.FirstJoinMessage
			if len(msg) > 0 {
				p.log.Info(p.settings.FirstJoinMessage, uuid)
			}
			return player.Config{}, nil, errors.New("bedrock-gophers/provider: player data not found")
		}
		p.log.Error("bedrock-gophers/provider: error reading file: %s", err)
		return player.Config{}, nil, err
	}

	dec := json.NewDecoder(bytes.NewReader(buf))
	dec.SetNumberType(json.UnmarshalIntOrFloat)
	err = dec.Decode(&playerDat)

	if err != nil {
		p.log.Error("bedrock-gophers/provider: error unmarshalling: %s", err)
		return player.Config{}, nil, err
	}
	dat := p.convertSavablePlayerData(playerDat, wrld)
	p.dataMu.Lock()
	p.data[uuid] = dat
	p.dataMu.Unlock()
	return dat.Config, dat.World, nil
}

// Close closes the player provider and flushes the player data to disk.
func (p *Provider) Close() error {
	p.closed = true
	p.data = nil
	return nil
}
