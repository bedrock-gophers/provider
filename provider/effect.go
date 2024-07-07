package provider

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"time"
)

// EffectData is a struct that contains the data of an effect that can be used to create an effect.
type EffectData struct {
	ID            int
	Amplifier     int
	Duration      time.Duration
	Ambient       bool
	ShowParticles bool
}

// ToEffect converts the EffectData to an effect.Effect.
func (e EffectData) ToEffect() effect.Effect {
	typ, ok := effect.ByID(e.ID)
	if !ok {
		return effect.Effect{}
	}
	if e.Duration == 0 {
		return effect.NewInstant(typ, e.Amplifier)
	}

	lastingType, ok := typ.(effect.LastingType)
	if !ok {
		return effect.Effect{}
	}

	ef := effect.New(lastingType, e.Amplifier, e.Duration)
	if e.Ambient {
		ef = effect.NewAmbient(lastingType, e.Amplifier, e.Duration)
	}

	if !e.ShowParticles {
		return ef.WithoutParticles()
	}
	return ef
}
