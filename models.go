package main

import (
	"pokedex/internal/pokeapi"
	"time"
)

type BallType struct {
	Name       string
	Multiplier float64
}

var balls = map[string]BallType{
	"pokeball":  {Name: "Pokeball", Multiplier: 1.0},
	"greatball": {Name: "Great Ball", Multiplier: 1.5},
	"ultraball": {Name: "Ultra Ball", Multiplier: 2.0},
}

const (
	defaultBall    = "pokeball"
	maxPartySize   = 6
	evolutionDelay = 5 * time.Minute
)

type PartyMember struct {
	pokeapi.Pokemon
	Level        int       `json:"level"`
	XP           int       `json:"xp"`
	CaughtAt     time.Time `json:"caught_at"`
	FullyEvolved bool      `json:"fully_evolved"`
}

func newPartyMember(p pokeapi.Pokemon) PartyMember {
	return PartyMember{
		Pokemon:  p,
		Level:    1,
		CaughtAt: time.Now(),
	}
}

func (pm *PartyMember) addXP(amount int) bool {
	pm.XP += amount
	newLevel := pm.XP/100 + 1
	if newLevel > pm.Level {
		pm.Level = newLevel
		return true
	}
	return false
}

func (pm PartyMember) stat(name string) int {
	for _, s := range pm.Stats {
		if s.Stat.Name == name {
			return s.BaseStat + pm.Level
		}
	}
	return 0
}
