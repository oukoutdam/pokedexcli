package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
)

type config struct {
	Next     string
	Previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cfg := config{}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}

		userInput := scanner.Text()
		userInputSlice := cleanInput(userInput)
		if len(userInputSlice) == 0 {
			continue
		}
		cmd, ok := cliCommands[userInputSlice[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := cmd.callback(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, "error running command:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, name := range slices.Sorted(maps.Keys(cliCommands)) {
		fmt.Printf("%s: %s\n", name, cliCommands[name].description)
	}
	return nil
}

type locationAreas struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func commandMap(cfg *config) error {
	locationAreaURL := "https://pokeapi.co/api/v2/location-area"
	if cfg.Next != "" {
		locationAreaURL = cfg.Next
	}

	las, err := fetchLocationAreas(locationAreaURL)
	if err != nil {
		return err
	}

	cfg.Next = las.Next
	cfg.Previous = las.Previous
	for _, la := range las.Results {
		fmt.Println(la.Name)
	}
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreaURL := cfg.Previous

	las, err := fetchLocationAreas(locationAreaURL)
	if err != nil {
		return err
	}
	cfg.Next = las.Next
	cfg.Previous = las.Previous
	for _, la := range las.Results {
		fmt.Println(la.Name)
	}
	return nil
}

func fetchLocationAreas(url string) (locationAreas, error) {
	res, err := http.Get(url)
	if err != nil {
		return locationAreas{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return locationAreas{}, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var las locationAreas
	if err := json.NewDecoder(res.Body).Decode(&las); err != nil {
		return locationAreas{}, err
	}
	return las, nil
}
