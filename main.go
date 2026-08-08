package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"os"
	"slices"
	"time"
)

type config struct {
	next       string
	previous   string
	pokeClient pokeAPIClient
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type pokeAPIClient struct {
	httpClient http.Client
}

func NewPokeAPIClient(timeout time.Duration) pokeAPIClient {
	return pokeAPIClient{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
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
	pokeClient := NewPokeAPIClient(5 * time.Second)
	cfg := config{
		pokeClient: pokeClient,
	}

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

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationAreasResponse struct {
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

func commandMap(cfg *config) error {
	locationAreaURL := "https://pokeapi.co/api/v2/location-area"
	if cfg.next != "" {
		locationAreaURL = cfg.next
	}

	areas, err := cfg.pokeClient.fetchLocationAreas(locationAreaURL)
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

func commandMapb(cfg *config) error {
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	locationAreaURL := cfg.previous

	areas, err := cfg.pokeClient.fetchLocationAreas(locationAreaURL)
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

func (c *pokeAPIClient) fetchLocationAreas(url string) (locationAreasResponse, error) {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return locationAreasResponse{}, err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		return locationAreasResponse{}, fmt.Errorf("bad status code: %d", res.StatusCode)
	}

	var areas locationAreasResponse
	if err := json.NewDecoder(res.Body).Decode(&areas); err != nil {
		return locationAreasResponse{}, err
	}
	return areas, nil
}
