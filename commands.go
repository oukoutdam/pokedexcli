package main

import (
	"errors"
	"fmt"
	"maps"
	"math"
	"math/rand/v2"
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
		"catch": {
			name:        "catch",
			description: "catch <pokemon-name>: Try to catch a pokemon",
			callback:    commandCatch,
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

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return errors.New("usage: catch <pokemon-name>")
	}

	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := cfg.pokeClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	catchRate := calculateCatchRate(pokemon.BaseExperience)
	caught := attemptCatch(catchRate)
	if caught {
		if _, ok := cfg.pokedex[pokemon.Name]; ok {
			fmt.Printf("%s was caught! (already in your pokedex)\n", pokemon.Name)
		} else {
			fmt.Printf("%s was caught!\n", pokemon.Name)
			cfg.pokedex[pokemon.Name] = pokemon
		}
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func calculateCatchRate(baseExperience int) float64 {
	const minCatchRate, maxCatchRate = 0.05, 0.95
	const minExp, maxExp = 36.0, 608.0

	clampedExp := math.Min(math.Max(float64(baseExperience), minExp), maxExp)

	return maxCatchRate - ((clampedExp-minExp)/(maxExp-minExp))*(maxCatchRate-minCatchRate)
}

func attemptCatch(catchRate float64) bool {
	return rand.Float64() < catchRate
}
