package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Jarimus/pokedexcli/internal/pokecache"
	"github.com/Jarimus/pokedexcli/internal/pokedex"
)

// Initiate global variables
var cliCommands map[string]cliCommand
var Cache = pokecache.NewCache(5 * time.Second)
var OurPokemon []pokedex.Pokemon

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	cliCommands = getCommands()

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
