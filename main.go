package main

import (
	"bufio"
	"fmt"
	"os"
)

var cliCommands map[string]cliCommand

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
	}

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Commands:")

	for _, command := range cliCommands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	for {

		fmt.Print("Pokedex > ")

		// Wait for user input
		if scanner.Scan() {

			// Clean the input: lowercase, split into words
			words := cleanInputString(scanner.Text())

			// Print first word of the input, if it exists. Stop program on "exit"
			if len(words) > 0 {
				input := words[0]

				command, ok := cliCommands[input]
				if ok {
					command.callback()
				}

			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}
