package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func startRepl() {
	nextPage := getConfig()
	for {
		fmt.Print("Pokedex > ")
		input := bufio.NewScanner(os.Stdin)
		ok := input.Scan()
		if !ok {
			if input.Err() == nil {
				return
			}
			fmt.Printf("Error while scanning user input: %v", input.Err())
			return
		}
		text := cleanInput(input.Text())
		if len(text) == 0 {
			continue
		}
		command, exists := getCommands()[text[0]]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(nextPage()); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(text)))
}

func commandExit(c *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMapf(c *Config) error {
	if c == nil {
		return fmt.Errorf("error: c (*config) was nil")
	}
	req, err := http.NewRequest("GET", c.nextUrl, nil)
	if err != nil {
		return fmt.Errorf("error while creating the request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while doing the request: %w", err)
	}
	defer res.Body.Close()

	var data LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("error while decoding the response to json format: %w", err)
	}
	c.setNextUrl(data.NextUrl)
	c.setPrevUrl(data.PrevUrl)

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func commandMapb(c *Config) error {
	if c == nil {
		return fmt.Errorf("error: c (*config) was nil")
	}
	if c.prevUrl == "" {
		return fmt.Errorf("you're on the first page")
	}
	req, err := http.NewRequest("GET", c.prevUrl, nil)
	if err != nil {
		return fmt.Errorf("error while creating the request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error while doing the request: %w", err)
	}
	defer res.Body.Close()

	var data LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("error while decoding the response to json format: %w", err)
	}
	c.setNextUrl(data.NextUrl)
	c.setPrevUrl(data.PrevUrl)

	for _, location := range data.Results {
		fmt.Println(location.Name)
	}

	return nil
}
