package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Jarimus/pokedexcli/internal/pokecache"
)

// Initiate global variables
var cliCommands map[string]cliCommand
var mapCache = pokecache.NewCache(5 * time.Second)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	cliCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exits Pokedex",
			callback:    commandExit,
			config:      &config{},
		},
		"help": {
			name:        "help",
			description: "Display commands",
			callback:    commandHelp,
			config:      &config{},
		},
		"map": {
			name:        "map",
			description: "List the next 20 locations of the Pokemon World",
			callback:    commandMap,
			config:      &config{},
		},
		"mapb": {
			name:        "mapb",
			description: "List the previous 20 locations of the Pokemon World",
			callback:    commandMapBack,
			config:      &config{},
		},
		"explore": {
			name:        "explore",
			description: "Explore a location for Pokemon.",
			callback:    commandExplore,
			config:      &config{},
		},
	}

	fmt.Println("####################\nWelcome to the Pokedex!")

	commandHelp("")

	for {

		fmt.Print("Pokedex > ")

		// Wait for user input
		if scanner.Scan() {

			// Clean the input: lowercase, split into words
			words := cleanInputString(scanner.Text())

			// Parse the command and they argument, if they exist.
			if len(words) > 0 {
				com := words[0]

				command, ok := cliCommands[com]

				var arg string
				if len(words) >= 2 {
					arg = words[1]
				}

				if ok {
					command.callback(arg)
				} else {
					fmt.Printf("####################\nInvalid command\n####################\n")
				}

			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error using scanner.")
		}
	}
}
