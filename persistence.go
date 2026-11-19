package main

import (
	"encoding/json"
	"os"

	"pokedex/internal/pokeapi"
)

const saveFile = "pokedex_save.json"

type saveData struct {
	Pokedex map[string]pokeapi.Pokemon `json:"pokedex"`
	Party   []PartyMember              `json:"party"`
}

func saveGame(cfg *config) error {
	data := saveData{
		Pokedex: cfg.Pokedex,
		Party:   cfg.Party,
	}
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(saveFile, b, 0644)
}

func loadGame(cfg *config) error {
	b, err := os.ReadFile(saveFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var data saveData
	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}
	if data.Pokedex != nil {
		cfg.Pokedex = data.Pokedex
	}
	if data.Party != nil {
		cfg.Party = data.Party
	}
	return nil
}
