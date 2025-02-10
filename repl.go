package main

import (
	"fmt"
	"os"
	"strings"

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
	callback    func() error
	config      *config
}

// Command functions
func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	for _, command := range cliCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}
	return nil
}

func commandMap() error {

	var locations pokedex.LocationArea
	var err error

	// Get the locations from the pokedex api. If the locations have been previously received, show the next 20 locations (config.Next)
	if cliCommands["map"].config.Next == "" {
		locations, err = pokedex.MapRequest("https://pokeapi.co/api/v2/location-area/")
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	} else {
		locations, err = pokedex.MapRequest(cliCommands["map"].config.Next)
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	}

	cliCommands["map"].config.Prev = locations.Previous
	cliCommands["map"].config.Next = locations.Next

	for _, location := range locations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapBack() error {
	previousLocations := cliCommands["map"].config.Prev
	var locations pokedex.LocationArea
	var err error

	if previousLocations == "" {
		fmt.Println("You are at the first page.")
		return nil
	} else {
		locations, err = pokedex.MapRequest(previousLocations)
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
