package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"slices"
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
func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "\nExit Pokedex.",
			callback:    commandExit,
			config:      &config{},
		},
		"help": {
			name:        "help",
			description: "\nDisplay commands.",
			callback:    commandHelp,
			config:      &config{},
		},
		"map": {
			name:        "map",
			description: "\nList the next 20 locations of the Pokemon World.",
			callback:    commandMap,
			config:      &config{},
		},
		"mapb": {
			name:        "mapb",
			description: "\nList the previous 20 locations of the Pokemon World.",
			callback:    commandMapBack,
			config:      &config{},
		},
		"explore": {
			name:        "explore",
			description: "\nExplore <location_name>\nExplore an area to find Pokemon.",
			callback:    commandExplore,
			config:      &config{},
		},
		"catch": {
			name:        "catch",
			description: "\nCatch <pokemon_name>\nAttempt to catch a Pokemon.\nSuccess chance increases for every Pokemon you have caught.",
			callback:    commandCatch,
			config:      &config{},
		},
		"inspect": {
			name:        "inspect",
			description: "\nInspect the statistics of a Pokemon you have caught.",
			callback:    commandInspect,
			config:      &config{},
		},
		"pokedex": {
			name:        "pokedex",
			description: "\nList the Pokemmon you have caught.",
			callback:    commandPokedex,
			config:      &config{},
		},
		"save": {
			name:        "save",
			description: "\nSave into a file the Pokemon you have caught.",
			callback:    commandSave,
			config:      &config{},
		},
		"load": {
			name:        "load",
			description: "\nLoad from a file the Pokemon you have caught.",
			callback:    commandLoad,
		},
	}
}

func commandExit(string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(string) error {
	fmt.Println("####################\nCommands:")
	fmt.Println()
	for _, command := range cliCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
		fmt.Println()
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("####################")
	return nil
}

func commandMap(string) error {

	var locations pokedex.Locations
	var err error

	// Get the locations from the pokedex api. If the locations have been previously received, show the next 20 locations (config.Next)
	if cliCommands["map"].config.Next == "" {
		locations, err = pokedex.LocationsRequest("https://pokeapi.co/api/v2/location-area/", Cache)
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	} else {
		locations, err = pokedex.LocationsRequest(cliCommands["map"].config.Next, Cache)
		if err != nil {
			return fmt.Errorf("mapreqest error: %v", err)
		}
	}

	cliCommands["map"].config.Prev = locations.Previous
	cliCommands["map"].config.Next = locations.Next

	fmt.Println("####################")

	for _, location := range locations.Results {
		time.Sleep(50 * time.Millisecond)
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
		locations, err = pokedex.LocationsRequest(previousLocations, Cache)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		cliCommands["map"].config.Prev = locations.Previous
		cliCommands["map"].config.Next = locations.Next

		fmt.Println("####################")

		for _, location := range locations.Results {
			time.Sleep(50 * time.Millisecond)
			fmt.Println(location.Name)
		}

		fmt.Println("####################")

	}

	return nil
}

func commandExplore(areaName string) error {

	if areaName == "" {
		println("Please give an area to explore as an argument.")
		return fmt.Errorf("no argument")
	}

	area, err := pokedex.AreaRequest(areaName, Cache)
	if err != nil {
		println("Area not found.")
		return fmt.Errorf("error: %w", err)
	}

	fmt.Printf("Exploring %s...\n", area.Name)
	time.Sleep(time.Second)

	if len(area.Encounters) > 0 {
		fmt.Printf("Found the following Pokemon:\n")
		for _, encounter := range area.Encounters {
			time.Sleep(150 * time.Millisecond)
			fmt.Printf("- %s\n", encounter.Pokemon.Name)
		}

	} else {
		fmt.Println("No pokemon here")
	}

	return nil
}

func commandCatch(target_pokemon string) error {
	// base exp 0 - 255, normalize to 0 - 100. Random number from 1 - 100 must be equal to or higher than base exp.

	// Set up a list of pokemon names user has caught to compare the target pokemon to
	var PokemonCaughtNames []string
	for _, pokemon := range PokemonCaught {
		PokemonCaughtNames = append(PokemonCaughtNames, pokemon.Name)
	}

	// Compare target and user's pokemon. Return if already caught.
	if slices.Contains(PokemonCaughtNames, target_pokemon) {
		fmt.Printf("You have already caught a %s\n", target_pokemon)
		return nil
	}

	// Request pokemon data from the server
	pokemon, err := pokedex.PokemonRequest(target_pokemon, Cache)
	if err != nil {
		return fmt.Errorf("error requesting pokemon from server: %w", err)
	}

	// Try to catch the pokemon
	targetValue := pokemon.BaseExperience * 100 / 255
	ourBonus := len(PokemonCaught)
	ourRoll := (rand.Intn(100) + 1) + ourBonus // Random number from 1 - 100, +1/pokemon caught
	successChance := float64(100 - (targetValue - ourBonus))

	fmt.Printf("Throwing a Pokeball at %s...\nSuccess chance: %.1f %%\n", target_pokemon, successChance)

	time.Sleep(time.Second)

	if ourRoll >= targetValue {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		PokemonCaught = append(PokemonCaught, pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}

	return nil
}

func commandInspect(target_pokemon string) error {

	for _, pokemon := range PokemonCaught {
		if pokemon.Name == target_pokemon {
			fmt.Printf("Name: %s\n", pokemon.Name)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Height: %d\n", pokemon.Height)
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Weight: %d\n", pokemon.Weight)
			time.Sleep(100 * time.Millisecond)

			fmt.Printf("Stats:\n")
			for _, stat := range pokemon.Stats {
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("\t- %s: %d\n", stat.Stat.Name, stat.BaseStat)
			}

			fmt.Printf("Types:\n")
			for _, t := range pokemon.Types {
				time.Sleep(100 * time.Millisecond)
				fmt.Printf("\t- %s\n", t.Type.Name)
			}

			return nil
		}
	}

	fmt.Println("You have not caught that pokemon.")

	return nil
}

func commandPokedex(string) error {

	if len(PokemonCaught) == 0 {
		println("You have not caught any Pokemon.")
		return nil
	}

	println("Your Pokedex:")

	for _, pokemon := range PokemonCaught {
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("\t- %s\n", pokemon.Name)
	}

	return nil
}

func commandSave(string) error {

	if len(PokemonCaught) == 0 {
		println("No Pokemon caught. Save aborted.")
		return nil
	}

	fileData, err := json.Marshal(PokemonCaught)
	if err != nil {
		return fmt.Errorf("error building json from the list of Pokemon structs: %w", err)
	}

	println("Caught Pokemon saved successfully to 'PokemonCaught.json'.")

	return os.WriteFile("PokemonCaught.json", fileData, 0644)

}

func commandLoad(string) error {
	data, err := os.ReadFile("PokemonCaught.json")
	if err != nil {
		println("Error loading data: ", err)
		return err
	}
	err = json.Unmarshal(data, &PokemonCaught)
	if err != nil {
		println("Error decoding 'PokemonCaught.json': ", err)
		return err
	}

	println("Save file 'PokemonCaught.json' loaded successfully!")
	println("Number of Pokemon caught: ", len(PokemonCaught))

	return nil
}
