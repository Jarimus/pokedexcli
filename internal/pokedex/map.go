package pokedex

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func MapRequest(url string) (LocationArea, error) {

	resp, err := http.Get(url)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error getting map: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error: %v", err)
	}

	defer resp.Body.Close()

	var locations LocationArea

	err = json.Unmarshal(body, &locations)
	if err != nil {
		return LocationArea{}, fmt.Errorf("error: %v", err)
	}

	return locations, nil
}
