package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jarimus/pokedexcli/internal/pokecache"
)

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	Moves []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func PokemonRequest(pokemon string, cache pokecache.Cache) (Pokemon, error) {
	var body []byte
	var err error

	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon

	// First look for the url in the cache
	body, ok := cache.Get(url)

	// If the data is not cached, initiate GET.
	if !ok {
		// Get locations from the pokedex api
		resp, err := http.Get(url)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error retrieving pokemon from the server: %w", err)
		}

		// read and store in memory the response body.
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return Pokemon{}, fmt.Errorf("error: %v", err)
		}
		defer resp.Body.Close()

		// store the response in the cache
		cache.Add(url, body)
	}

	var found_pokemon Pokemon

	err = json.Unmarshal(body, &found_pokemon)
	if err != nil {
		return Pokemon{}, fmt.Errorf("error: %v", err)
	}

	return found_pokemon, nil
}
