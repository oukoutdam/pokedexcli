package main

import (
	"errors"
	"fmt"
	"maps"
	"os"
	"slices"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

var cliCommands map[string]cliCommand

func init() {
	cliCommands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of 20 location areas in the Pokemon world. Each subsequent call to map will display the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of 20 previous location areas in the Pokemon world. Each subsequent call to map will display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "explore <area-name>: Lists the pokemon found in the location area",
			callback:    commandExplore,
		},
	}
}

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, name := range slices.Sorted(maps.Keys(cliCommands)) {
		fmt.Printf("%s: %s\n", name, cliCommands[name].description)
	}
	return nil
}

func commandMap(cfg *config, args ...string) error {
	areas, err := cfg.pokeClient.ListLocationAreas(cfg.next)
	if err != nil {
		return err
	}

	cfg.next = areas.Next
	cfg.previous = areas.Previous
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandMapb(cfg *config, args ...string) error {
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, err := cfg.pokeClient.ListLocationAreas(cfg.previous)
	if err != nil {
		return err
	}
	cfg.next = areas.Next
	cfg.previous = areas.Previous
	for _, area := range areas.Results {
		fmt.Println(area.Name)
	}
	return nil
}

func commandExplore(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: explore <area-name>")
	}

	areaName := args[0]
	areaInfo, err := cfg.pokeClient.GetLocationArea(areaName)
	if err != nil {
		return err
	}

	fmt.Printf("exploring %s...\nFound Pokemon:\n", areaName)
	for _, pokemonEncounter := range areaInfo.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name)
	}

	return nil
}
