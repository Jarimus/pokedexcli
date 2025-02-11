package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jarimus/pokedexcli/internal/pokecache"
)

type Locations struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Area struct {
	Name        string `json:"name"`
	PokemonList []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func LocationsRequest(url string, cache pokecache.Cache) (Locations, error) {

	var body []byte
	var err error

	// First look for the url in the cache
	body, ok := cache.Get(url)

	// If the data is not cached, initiate GET.
	if !ok {
		// Get locations from the pokedex api
		resp, err := http.Get(url)
		if err != nil {
			return Locations{}, fmt.Errorf("error getting map: %v", err)
		}

		// read and store in memory the response body.
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return Locations{}, fmt.Errorf("error: %v", err)
		}
		defer resp.Body.Close()

		// store the response in the cache
		cache.Add(url, body)
	}

	var locations Locations

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return Locations{}, fmt.Errorf("error: %v", err)
	}

	return locations, nil
}

func AreaRequest(url string, cache pokecache.Cache) (Area, error) {

	var body []byte
	var err error

	url = "https://pokeapi.co/api/v2/location-area/" + url

	// First look for the url in the cache
	body, ok := cache.Get(url)

	// If the data is not cached, initiate GET.
	if !ok {
		// Get locations from the pokedex api
		resp, err := http.Get(url)
		if err != nil {
			return Area{}, fmt.Errorf("error getting map: %w", err)
		}

		// read and store in memory the response body.
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return Area{}, fmt.Errorf("error: %w", err)
		}
		defer resp.Body.Close()

		// store the response in the cache
		cache.Add(url, body)
	}

	var area Area

	err = json.Unmarshal(body, &area)
	if err != nil {
		return Area{}, fmt.Errorf("error: %w", err)
	}

	return area, nil
}
