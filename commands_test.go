package main

import (
	"os"
	"testing"
	"time"

	"pokedex/internal/pokeapi"
)

func TestCalcDamage(t *testing.T) {
	for i := 0; i < 200; i++ {
		dmg := calcDamage(50, 20)
		if dmg < 1 {
			t.Errorf("calcDamage returned %d, want >= 1", dmg)
		}
	}
}

func TestCalcDamageHighDefense(t *testing.T) {
	for i := 0; i < 200; i++ {
		dmg := calcDamage(1, 1000)
		if dmg < 1 {
			t.Errorf("calcDamage with high defense returned %d, want >= 1", dmg)
		}
	}
}

func TestBallTypes(t *testing.T) {
	cases := []struct {
		name       string
		multiplier float64
	}{
		{"pokeball", 1.0},
		{"greatball", 1.5},
		{"ultraball", 2.0},
	}
	for _, c := range cases {
		ball, ok := balls[c.name]
		if !ok {
			t.Errorf("ball %q not found in registry", c.name)
			continue
		}
		if ball.Multiplier != c.multiplier {
			t.Errorf("ball %q: expected multiplier %v, got %v", c.name, c.multiplier, ball.Multiplier)
		}
	}
}

func TestPartyMemberAddXP(t *testing.T) {
	p := newPartyMember(pokeapi.Pokemon{Name: "bulbasaur"})
	if p.Level != 1 {
		t.Errorf("new party member should start at level 1, got %d", p.Level)
	}

	leveled := p.addXP(99)
	if leveled {
		t.Errorf("should not level up at 99 XP")
	}
	if p.Level != 1 {
		t.Errorf("level should still be 1, got %d", p.Level)
	}

	leveled = p.addXP(1) // now 100 XP
	if !leveled {
		t.Errorf("should level up at 100 XP")
	}
	if p.Level != 2 {
		t.Errorf("expected level 2, got %d", p.Level)
	}
}

func TestPartyMemberStat(t *testing.T) {
	p := newPartyMember(pokeapi.Pokemon{
		Name: "squirtle",
		Stats: []struct {
			BaseStat int `json:"base_stat"`
			Stat     struct {
				Name string `json:"name"`
			} `json:"stat"`
		}{
			{BaseStat: 44, Stat: struct {
				Name string `json:"name"`
			}{Name: "hp"}},
		},
	})
	// stat = BaseStat + Level (1)
	if got := p.stat("hp"); got != 45 {
		t.Errorf("expected hp stat 45, got %d", got)
	}
	if got := p.stat("nonexistent"); got != 0 {
		t.Errorf("expected 0 for missing stat, got %d", got)
	}
}

func TestSaveLoad(t *testing.T) {
	t.Cleanup(func() { _ = os.Remove(saveFile) })

	cfg := &config{
		Pokedex: map[string]pokeapi.Pokemon{
			"pikachu": {Name: "pikachu", BaseExperience: 112},
		},
		Party: []PartyMember{
			{
				Pokemon:  pokeapi.Pokemon{Name: "pikachu"},
				Level:    5,
				XP:       450,
				CaughtAt: time.Now(),
			},
		},
	}

	if err := saveGame(cfg); err != nil {
		t.Fatalf("saveGame failed: %v", err)
	}

	cfg2 := &config{
		Pokedex: make(map[string]pokeapi.Pokemon),
	}
	if err := loadGame(cfg2); err != nil {
		t.Fatalf("loadGame failed: %v", err)
	}

	if _, ok := cfg2.Pokedex["pikachu"]; !ok {
		t.Errorf("expected pikachu in loaded Pokedex")
	}
	if len(cfg2.Party) != 1 {
		t.Errorf("expected 1 party member, got %d", len(cfg2.Party))
	}
	if cfg2.Party[0].Level != 5 {
		t.Errorf("expected party[0] level 5, got %d", cfg2.Party[0].Level)
	}
}

func TestFindNextEvolution(t *testing.T) {
	chain := pokeapi.EvolutionChainLink{}
	chain.Species.Name = "bulbasaur"
	ivysaur := pokeapi.EvolutionChainLink{}
	ivysaur.Species.Name = "ivysaur"
	venusaur := pokeapi.EvolutionChainLink{}
	venusaur.Species.Name = "venusaur"
	ivysaur.EvolvesTo = []pokeapi.EvolutionChainLink{venusaur}
	chain.EvolvesTo = []pokeapi.EvolutionChainLink{ivysaur}

	cases := []struct {
		input    string
		expected string
		found    bool
	}{
		{"bulbasaur", "ivysaur", true},
		{"ivysaur", "venusaur", true},
		{"venusaur", "", false},
		{"mewtwo", "", false},
	}
	for _, c := range cases {
		got, ok := pokeapi.FindNextEvolution(chain, c.input)
		if ok != c.found {
			t.Errorf("FindNextEvolution(%q): expected found=%v, got %v", c.input, c.found, ok)
		}
		if got != c.expected {
			t.Errorf("FindNextEvolution(%q): expected %q, got %q", c.input, c.expected, got)
		}
	}
}
