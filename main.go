package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/oukoutdam/pokedexcli/internal/pokeapi"
)

type config struct {
	next       string
	previous   string
	pokeClient pokeapi.Client
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	pokeClient := pokeapi.NewClient(5 * time.Second)
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
