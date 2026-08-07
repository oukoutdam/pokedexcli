package main

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

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
		if err := cmd.callback(); err != nil {
			fmt.Fprintln(os.Stderr, "error running command:", err)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")

	for _, name := range slices.Sorted(maps.Keys(cliCommands)) {
		fmt.Printf("%s: %s\n", name, cliCommands[name].description)
	}
	return nil
}
