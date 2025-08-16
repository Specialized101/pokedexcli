package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func startRepl() {
	input := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := input.Scan()
		if !ok {
			log.Fatal("Error while scanning user input")
			return
		}
		text := cleanInput(input.Text())
		if len(text) == 0 {
			continue
		}
		fmt.Printf("Your command was: %s\n", text[0])
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(text)))
}
