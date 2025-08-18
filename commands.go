package main

import (
	"github.com/Specialized101/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, *pokecache.Cache) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Help menu showing all available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the next 20 location areas",
			callback:    commandMapf,
		},
		"mapb": {
			name:        "mapb",
			description: "Display the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Display the list of all the Pokemons found at a location",
			callback:    commandExplore,
		},
	}
}
