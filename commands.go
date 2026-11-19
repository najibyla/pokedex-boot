package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"pokedex/internal/pokeapi"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help":    {name: "help", description: "Displays a help message", callback: commandHelp},
		"exit":    {name: "exit", description: "Exit the Pokedex", callback: commandExit},
		"map":     {name: "map", description: "Show next 20 location areas", callback: commandMap},
		"mapb":    {name: "mapb", description: "Show previous 20 location areas", callback: commandMapb},
		"explore": {name: "explore", description: "explore <area> — list Pokemon in an area", callback: commandExplore},
		"go":      {name: "go", description: "go <n> — explore area #n from last map (may trigger random encounter)", callback: commandGoArea},
		"catch":   {name: "catch", description: "catch <pokemon> [pokeball|greatball|ultraball]", callback: commandCatch},
		"inspect": {name: "inspect", description: "inspect <pokemon> — view caught Pokemon's details", callback: commandInspect},
		"pokedex": {name: "pokedex", description: "List all caught Pokemon", callback: commandPokedex},
		"party":   {name: "party", description: "Show your party members and their levels", callback: commandParty},
		"battle":  {name: "battle", description: "battle <pokemon> — fight a wild Pokemon with your lead", callback: commandBattle},
		"save":    {name: "save", description: "Save your progress to disk", callback: commandSave},
	}
}

func commandHelp(cfg *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range cfg.Commands {
		fmt.Printf("  %-8s %s\n", cmd.name, cmd.description)
	}
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args []string) error {
	resp, err := pokeapi.GetLocationAreas(cfg.Next, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	cfg.LastAreas = nil
	for i, area := range resp.Results {
		cfg.LastAreas = append(cfg.LastAreas, area.Name)
		fmt.Printf("  %2d. %s\n", i+1, area.Name)
	}
	fmt.Println("(use 'go <n>' to explore a numbered area)")
	return nil
}

func commandMapb(cfg *config, args []string) error {
	if cfg.Previous == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	resp, err := pokeapi.GetLocationAreas(cfg.Previous, cfg.Cache)
	if err != nil {
		return err
	}
	cfg.Next = resp.Next
	cfg.Previous = resp.Previous
	cfg.LastAreas = nil
	for i, area := range resp.Results {
		cfg.LastAreas = append(cfg.LastAreas, area.Name)
		fmt.Printf("  %2d. %s\n", i+1, area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: explore <location-area>")
	}
	return exploreArea(cfg, args[0], false)
}

func commandGoArea(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: go <number>")
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 || n > len(cfg.LastAreas) {
		return fmt.Errorf("pick a number between 1 and %d", len(cfg.LastAreas))
	}
	return exploreArea(cfg, cfg.LastAreas[n-1], true)
}

func exploreArea(cfg *config, name string, randomEncounters bool) error {
	if randomEncounters && rand.Intn(100) < 25 {
		detail, err := pokeapi.GetLocationAreaDetail(name, cfg.Cache)
		if err != nil {
			return err
		}
		if len(detail.PokemonEncounters) > 0 {
			enc := detail.PokemonEncounters[rand.Intn(len(detail.PokemonEncounters))]
			wild := enc.Pokemon.Name
			fmt.Printf("A wild %s appeared! Use 'catch %s' to catch it!\n", wild, wild)
			return nil
		}
	}

	fmt.Printf("Exploring %s...\n", name)
	detail, err := pokeapi.GetLocationAreaDetail(name, cfg.Cache)
	if err != nil {
		return err
	}
	fmt.Println("Found Pokemon:")
	for _, enc := range detail.PokemonEncounters {
		fmt.Printf(" - %s\n", enc.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: catch <pokemon> [pokeball|greatball|ultraball]")
	}
	name := args[0]
	ballName := defaultBall
	if len(args) >= 2 {
		ballName = args[1]
	}
	ball, ok := balls[ballName]
	if !ok {
		return fmt.Errorf("unknown ball: %q (use pokeball, greatball, or ultraball)", ballName)
	}

	fmt.Printf("Throwing a %s at %s...\n", ball.Name, name)
	pokemon, err := pokeapi.GetPokemon(name, cfg.Cache)
	if err != nil {
		return err
	}

	threshold := int(float64(40) * ball.Multiplier)
	if rand.Intn(pokemon.BaseExperience) < threshold {
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
		cfg.Pokedex[name] = pokemon
		if len(cfg.Party) < maxPartySize {
			cfg.Party = append(cfg.Party, newPartyMember(pokemon))
			fmt.Printf("%s was added to your party!\n", name)
		}
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: inspect <pokemon>")
	}
	pokemon, ok := cfg.Pokedex[args[0]]
	if !ok {
		fmt.Println("you have not caught that pokemon")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
	return nil
}

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")
	for name := range cfg.Pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandParty(cfg *config, args []string) error {
	if len(cfg.Party) == 0 {
		fmt.Println("Your party is empty! Catch some Pokemon.")
		return nil
	}
	fmt.Println("Your Party:")
	for i, m := range cfg.Party {
		fmt.Printf("  %d. %s  Lv.%d  XP:%d\n", i+1, m.Name, m.Level, m.XP)
	}
	return nil
}

func commandBattle(cfg *config, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: battle <pokemon>")
	}
	if len(cfg.Party) == 0 {
		fmt.Println("Your party is empty! Catch some Pokemon first.")
		return nil
	}

	wildPokemon, err := pokeapi.GetPokemon(args[0], cfg.Cache)
	if err != nil {
		return err
	}
	wild := newPartyMember(wildPokemon)

	won, xp := simulateBattle(cfg.Party[0], wild)
	if won {
		fmt.Printf("%s won and earned %d XP!\n", cfg.Party[0].Name, xp)
		if cfg.Party[0].addXP(xp) {
			fmt.Printf("*** %s leveled up to Lv.%d! ***\n", cfg.Party[0].Name, cfg.Party[0].Level)
		}
	} else {
		fmt.Printf("%s was defeated!\n", cfg.Party[0].Name)
	}
	return nil
}

func commandSave(cfg *config, args []string) error {
	if err := saveGame(cfg); err != nil {
		return err
	}
	fmt.Printf("Game saved to %s\n", saveFile)
	return nil
}
