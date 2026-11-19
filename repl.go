package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
	"pokedex/internal/pokeapi"
	"pokedex/internal/pokecache"
)

type config struct {
	Next      *string
	Previous  *string
	Cache     *pokecache.Cache
	Pokedex   map[string]pokeapi.Pokemon
	Party     []PartyMember
	LastAreas []string
	// Le registre des commandes vit dans l'état partagé, pour qu'une
	// commande puisse y accéder (commandHelp s'en sert pour lister les
	// commandes disponibles).
	Commands map[string]cliCommand
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

// startRepl reçoit la config déjà construite par main, au lieu de la créer
// elle-même : l'état partagé est ainsi initialisé en un seul endroit.
func startRepl(cfg *config) {
	if err := loadGame(cfg); err != nil {
		fmt.Printf("Warning: could not load save: %v\n", err)
	}

	interactive := term.IsTerminal(int(os.Stdin.Fd()))
	var scanner *bufio.Scanner
	if !interactive {
		scanner = bufio.NewScanner(os.Stdin)
	}

	var history []string

	for {
		if interactive {
			checkEvolutions(cfg)
		}

		var line string
		if interactive {
			var err error
			line, err = readLineRaw("Pokedex > ", history)
			if err != nil {
				return
			}
		} else {
			fmt.Print("Pokedex > ")
			if !scanner.Scan() {
				return
			}
			line = scanner.Text()
		}

		words := cleanInput(line)
		if len(words) == 0 {
			continue
		}

		history = append(history, line)

		cmd, ok := cfg.Commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := cmd.callback(cfg, words[1:]); err != nil {
			fmt.Println(err)
		}
	}
}

// readLineRaw reads a line in raw terminal mode, supporting up/down arrow history.
func readLineRaw(prompt string, history []string) (string, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, oldState)

	fmt.Print(prompt)

	var buf []byte
	histIdx := -1

	readByte := func() (byte, error) {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		return b[0], err
	}

	for {
		b, err := readByte()
		if err != nil {
			fmt.Print("\r\n")
			return string(buf), err
		}

		switch b {
		case 3: // Ctrl+C
			fmt.Print("\r\n")
			os.Exit(0)
		case 13, 10: // Enter
			fmt.Print("\r\n")
			return string(buf), nil
		case 127, 8: // Backspace
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
		case 27: // ESC — check for arrow keys
			b2, err := readByte()
			if err != nil || b2 != 91 {
				continue
			}
			b3, err := readByte()
			if err != nil {
				continue
			}
			switch b3 {
			case 65: // Up arrow
				if histIdx < len(history)-1 {
					histIdx++
					entry := history[len(history)-1-histIdx]
					fmt.Printf("\r\033[K%s%s", prompt, entry)
					buf = []byte(entry)
				}
			case 66: // Down arrow
				if histIdx > 0 {
					histIdx--
					entry := history[len(history)-1-histIdx]
					fmt.Printf("\r\033[K%s%s", prompt, entry)
					buf = []byte(entry)
				} else if histIdx == 0 {
					histIdx = -1
					fmt.Printf("\r\033[K%s", prompt)
					buf = []byte{}
				}
			}
		default:
			if b >= 32 {
				buf = append(buf, b)
				fmt.Print(string(b))
			}
		}
	}
}

// checkEvolutions checks all party members and evolves any that are ready.
func checkEvolutions(cfg *config) {
	for i := range cfg.Party {
		m := &cfg.Party[i]
		if m.FullyEvolved || time.Since(m.CaughtAt) < evolutionDelay {
			continue
		}
		species, err := pokeapi.GetPokemonSpecies(m.Name, cfg.Cache)
		if err != nil {
			continue
		}
		chain, err := pokeapi.GetEvolutionChain(species.EvolutionChain.URL, cfg.Cache)
		if err != nil {
			continue
		}
		nextName, ok := pokeapi.FindNextEvolution(chain.Chain, m.Name)
		if !ok {
			m.FullyEvolved = true
			continue
		}
		evolved, err := pokeapi.GetPokemon(nextName, cfg.Cache)
		if err != nil {
			continue
		}
		oldName := m.Name
		delete(cfg.Pokedex, oldName)
		cfg.Pokedex[nextName] = evolved
		m.Pokemon = evolved
		m.CaughtAt = time.Now()
		fmt.Printf("\n*** %s evolved into %s! ***\n", oldName, nextName)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
