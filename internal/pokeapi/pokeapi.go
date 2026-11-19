package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"pokedex/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2"

// --- Pokemon ---

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func GetPokemon(name string, cache *pokecache.Cache) (Pokemon, error) {
	return fetchJSON[Pokemon](baseURL+"/pokemon/"+name, cache)
}

// --- Location areas ---

type LocationAreaResponse struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationAreaDetail struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreas(pageURL *string, cache *pokecache.Cache) (LocationAreaResponse, error) {
	url := baseURL + "/location-area/"
	if pageURL != nil {
		url = *pageURL
	}
	return fetchJSON[LocationAreaResponse](url, cache)
}

func GetLocationAreaDetail(name string, cache *pokecache.Cache) (LocationAreaDetail, error) {
	return fetchJSON[LocationAreaDetail](baseURL+"/location-area/"+name, cache)
}

// --- Evolution ---

type PokemonSpecies struct {
	EvolutionChain struct {
		URL string `json:"url"`
	} `json:"evolution_chain"`
}

type EvolutionChainLink struct {
	Species struct {
		Name string `json:"name"`
	} `json:"species"`
	EvolvesTo []EvolutionChainLink `json:"evolves_to"`
}

type EvolutionChain struct {
	Chain EvolutionChainLink `json:"chain"`
}

func GetPokemonSpecies(name string, cache *pokecache.Cache) (PokemonSpecies, error) {
	return fetchJSON[PokemonSpecies](baseURL+"/pokemon-species/"+name, cache)
}

func GetEvolutionChain(url string, cache *pokecache.Cache) (EvolutionChain, error) {
	return fetchJSON[EvolutionChain](url, cache)
}

// FindNextEvolution returns the name of the next evolution form, if any.
func FindNextEvolution(chain EvolutionChainLink, name string) (string, bool) {
	if chain.Species.Name == name {
		if len(chain.EvolvesTo) > 0 {
			return chain.EvolvesTo[0].Species.Name, true
		}
		return "", false
	}
	for _, next := range chain.EvolvesTo {
		if evName, ok := FindNextEvolution(next, name); ok {
			return evName, true
		}
	}
	return "", false
}

// --- shared HTTP helper ---

func fetchJSON[T any](url string, cache *pokecache.Cache) (T, error) {
	var result T
	if data, ok := cache.Get(url); ok {
		return result, json.Unmarshal(data, &result)
	}
	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// Une réponse 404 ou 500 a un corps, mais ce n'est pas la ressource
	// demandée : on renvoie une erreur plutôt que de la mettre en cache.
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status %s for %s", resp.Status, url)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	// Le décodage passe avant la mise en cache : du JSON invalide ne doit
	// pas être mémorisé, sinon chaque appel suivant échouerait aussi, sans
	// même retenter la requête.
	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	cache.Add(url, data)
	return result, nil
}
