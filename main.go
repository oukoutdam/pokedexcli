package main

import (
	"bufio"
	"fmt"
	"os"
)

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
		fmt.Printf("Your command was: %s\n", userInputSlice[0])
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "error reading input:", err)
		os.Exit(1)
	}
}
