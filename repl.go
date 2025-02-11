package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Jarimus/pokedexcli/internal/pokedex"
)

// Input clean up to parse the commands
func cleanInputString(text string) []string {
	cleanWords := strings.Fields(strings.ToLower(text))

	return cleanWords
}

// command structs
type config struct {
	Next string
	Prev string
}

type cliCommand struct {
	name        string
	description string
	callback    func(string) error
	config      *config
}

// Command functions
func commandExit(string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(string) error {
	fmt.Println("####################\nCommands:")
	for _, command := range cliCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	fmt.Println("####################")
	return nil
}

func commandMap(string) error {

	var locations pokedex.Locations
	var err error

	// Get the locations from the pokedex api. If the locations have been previously received, show the next 20 locations (config.Next)
	if cliCommands["map"].config.Next == "" {
		locations, err = pokedex.LocationsRequest("https://pokeapi.co/api/v2/location-area/", mapCache)
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	} else {
		locations, err = pokedex.LocationsRequest(cliCommands["map"].config.Next, mapCache)
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	}

	cliCommands["map"].config.Prev = locations.Previous
	cliCommands["map"].config.Next = locations.Next

	fmt.Println("####################")

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	fmt.Println("####################")

	return nil
}

func commandMapBack(string) error {
	previousLocations := cliCommands["map"].config.Prev
	var locations pokedex.Locations
	var err error

	if previousLocations == "" {
		fmt.Println("You are at the first page.")
		return nil
	} else {
		locations, err = pokedex.LocationsRequest(previousLocations, mapCache)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		cliCommands["map"].config.Prev = locations.Previous
		cliCommands["map"].config.Next = locations.Next

		for _, location := range locations.Results {
			fmt.Println(location.Name)
		}
	}

	return nil
}

func commandExplore(areaName string) error {

	if areaName == "" {
		println("Please give an area to explore as an argument.")
		return fmt.Errorf("no argument")
	}

	area, err := pokedex.AreaRequest(areaName, mapCache)
	if err != nil {
		println("Area not found.")
		return fmt.Errorf("error: %w", err)
	}

	fmt.Printf("Exploring %s...\n", area.Name)
	time.Sleep(time.Second)

	if len(area.PokemonList) > 0 {
		fmt.Printf("Found the following Pokemon:\n")
		for _, pokemon := range area.PokemonList {
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("- %s\n", pokemon.Pokemon.Name)
		}

	} else {
		fmt.Println("No pokemon here")
	}

	return nil
}
