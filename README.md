# Pokedex CLI

A terminal-based Pokedex REPL built in Go, powered by the [PokeAPI](https://pokeapi.co/).

## Run

```bash
go run .
```

## Commands

| Command | Usage | Description |
|---------|-------|-------------|
| `help` | `help` | List all available commands |
| `exit` | `exit` | Quit the program |
| `map` | `map` | Show the next 20 location areas (numbered) |
| `mapb` | `mapb` | Show the previous 20 location areas |
| `go` | `go <n>` | Explore area #n from the last map view |
| `explore` | `explore <area-name>` | List all Pokemon in a named area |
| `catch` | `catch <pokemon> [ball]` | Attempt to catch a Pokemon |
| `inspect` | `inspect <pokemon>` | View a caught Pokemon's stats and types |
| `pokedex` | `pokedex` | List all caught Pokemon |
| `party` | `party` | Show your party with levels and XP |
| `battle` | `battle <pokemon>` | Fight a wild Pokemon using your lead party member |
| `save` | `save` | Save your progress to `pokedex_save.json` |

---

## Features

### Navigation

```
Pokedex > map
   1. canalave-city-area
   2. eterna-city-area
   ...
(use 'go <n>' to explore a numbered area)

Pokedex > go 2
Exploring eterna-city-area...
Found Pokemon:
 - budew
 - gastly
 ...
```

`map` / `mapb` paginate through 20 locations at a time. `go <n>` explores the nth area from the last `map` result.

---

### Random Encounters

Using `go <n>` has a **25% chance** of triggering a random wild Pokemon encounter instead of showing the area list.

```
Pokedex > go 3
A wild tentacool appeared! Use 'catch tentacool' to catch it!
```

---

### Catching Pokemon

```
Pokedex > catch squirtle
Throwing a Pokeball at squirtle...
squirtle was caught!
You may now inspect it with the inspect command.
squirtle was added to your party!
```

#### Ball Types

Different balls increase your catch chance:

| Ball | Command | Catch Multiplier |
|------|---------|-----------------|
| Pokeball | `catch <pokemon>` or `catch <pokemon> pokeball` | 1× |
| Great Ball | `catch <pokemon> greatball` | 1.5× |
| Ultra Ball | `catch <pokemon> ultraball` | 2× |

```
Pokedex > catch mewtwo ultraball
Throwing a Ultra Ball at mewtwo...
mewtwo was caught!
```

---

### Inspecting Pokemon

```
Pokedex > inspect squirtle
Name: squirtle
Height: 5
Weight: 90
Stats:
  -hp: 44
  -attack: 48
  -defense: 65
  -special-attack: 50
  -special-defense: 64
  -speed: 43
Types:
  - water
```

---

### Party, Battles & Leveling Up

Your party holds up to **6 Pokemon**. Use `battle` to fight wild Pokemon — your lead party member gains XP on a win. Every **100 XP** earns a level.

```
Pokedex > party
Your Party:
  1. squirtle  Lv.1  XP:0

Pokedex > battle pidgey
squirtle (Lv.1, HP:90) vs pidgey (Lv.1, HP:82)
  --- Round 1 ---
  squirtle hits pidgey for 31! (pidgey HP: 51)
  pidgey hits squirtle for 18! (squirtle HP: 72)
  ...
squirtle won and earned 25 XP!

Pokedex > party
Your Party:
  1. squirtle  Lv.1  XP:25
```

---

### Evolution

Caught Pokemon in your party will **automatically evolve** after 5 minutes. The evolution is checked before each command and announced when it happens.

```
*** caterpie evolved into metapod! ***
```

---

### Saving & Loading

Progress is saved to `pokedex_save.json` in the current directory and **automatically loaded** on startup.

```
Pokedex > save
Game saved to pokedex_save.json
```

---

### Command History (Interactive Mode)

When running interactively (not piped), use the **↑ / ↓ arrow keys** to cycle through your previous commands.

---

## Testing

```bash
# Run all unit tests
go test ./...

# Verbose output
go test ./... -v
```

### Manual test examples

```bash
# Ball types
printf "catch mewtwo pokeball\ncatch mewtwo greatball\ncatch mewtwo ultraball\nexit\n" | go run .

# Party, battle, and XP
printf "catch bulbasaur\ncatch bulbasaur\ncatch bulbasaur\nbattle pidgey\nparty\nexit\n" | go run .

# Directional exploration (run multiple times to see random encounters)
printf "map\ngo 3\nexit\n" | go run .

# Save and load persistence
printf "catch squirtle\ncatch squirtle\ncatch squirtle\nsave\nexit\n" | go run .
printf "pokedex\nparty\nexit\n" | go run .
```
