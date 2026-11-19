package main

import (
	"time"

	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

func main() {
	// Tout l'état partagé par le REPL et ses commandes est construit ici,
	// puis passé à startRepl.
	cfg := &config{
		Cache:    pokecache.NewCache(5 * time.Minute),
		Pokedex:  make(map[string]pokeapi.Pokemon),
		Commands: getCommands(),
	}
	startRepl(cfg)
}
