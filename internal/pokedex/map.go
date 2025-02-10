package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jarimus/pokedexcli/internal/pokecache"
)

type LocationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func MapRequest(url string, cache pokecache.Cache) (LocationArea, error) {

	var body []byte
	var err error

	// First look for the url in the cache
	body, ok := cache.Get(url)

	// If the data is not cached, initiate GET.
	if !ok {
		// Get locations from the pokedex api
		resp, err := http.Get(url)
		if err != nil {
			return LocationArea{}, fmt.Errorf("error getting map: %v", err)
		}

		// read and store in memory the response body.
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return LocationArea{}, fmt.Errorf("error: %v", err)
		}
		defer resp.Body.Close()

		// store the response in the cache
		cache.Add(url, body)
	}

	var locations LocationArea

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error: %v", err)
	}

	return locations, nil
}
