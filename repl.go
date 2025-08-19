package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"

	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Specialized101/pokedexcli/internal/pokecache"
)

func startRepl() {
	config := getConfig()
	cache := pokecache.NewCache(10 * time.Second)
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
		var param string
		if len(text) > 1 {
			param = text[1]
		}
		config().setParam(param)
		if err := command.callback(config(), cache); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(strings.TrimSpace(text)))
}

func commandExit(c *Config, cache *pokecache.Cache) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, cache *pokecache.Cache) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Printf("Usage:\n\n")
	for _, c := range getCommands() {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}

func commandMapf(c *Config, cache *pokecache.Cache) error {
	if c == nil {
		return fmt.Errorf("error: c (*config) was nil")
	}
	cachedData, ok := cache.Get(c.nextUrl)
	if ok {
		fmt.Println("LOG: Retrieving data from cache...")
		displayLocationsFromCache(c, cachedData)
		return nil
	}
	fmt.Println("LOG: Data not in cache, making a new request ...")

	res, err := makeGetRequest(c.nextUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("LOG: Saving data to cache...")

	cache.Add(c.nextUrl, rawData)

	displayLocationsFromCache(c, rawData)

	return nil
}

func commandMapb(c *Config, cache *pokecache.Cache) error {
	if c == nil {
		return fmt.Errorf("error: c (*config) was nil")
	}
	if c.prevUrl == "" {
		return fmt.Errorf("you're on the first page")
	}
	cachedData, ok := cache.Get(c.prevUrl)
	if ok {
		fmt.Println("LOG: Retrieving data from cache...")
		displayLocationsFromCache(c, cachedData)
		return nil
	}
	fmt.Println("LOG: Data not in cache, making a new request ...")
	res, err := makeGetRequest(c.prevUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("LOG: Saving data to cache...")

	cache.Add(c.prevUrl, rawData)

	displayLocationsFromCache(c, rawData)

	return nil
}

func commandExplore(c *Config, cache *pokecache.Cache) error {
	if c.param == "" {
		return fmt.Errorf("usage: explore [location name]\nexample: explore mt-coronet-2f\n\nTry map to see the list of locations")
	}

	fullUrl := fmt.Sprintf("%s/%s", LOCATION_AREA_URL, c.param)
	cachedData, ok := cache.Get(fullUrl)
	if ok {
		fmt.Println("Log: Retrieving data from cache...")
		displayPokemonListFromCache(cachedData)
		return nil
	}
	res, err := makeGetRequest(fullUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Println("LOG: Saving data to cache...")
	cache.Add(fullUrl, rawData)

	displayPokemonListFromCache(rawData)
	return nil
}

func commandCatch(c *Config, cache *pokecache.Cache) error {
	if c.param == "" {
		return fmt.Errorf("usage: catch [pokemon name]\nexample: catch pikachu")
	}

	fullUrl := fmt.Sprintf("%s/%s", POKEMON_URL, c.param)
	cachedData, ok := cache.Get(fullUrl)
	if !ok {
		res, err := makeGetRequest(fullUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		cachedData, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		cache.Add(fullUrl, cachedData)
		// fmt.Println(cachedData)~
	}
	var pokemon Pokemon

	err := json.Unmarshal(cachedData, &pokemon)
	if err != nil {
		fmt.Println("here")
		return err
	}

	succeeded := attemptToCatchPokemon(pokemon)
	if succeeded {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		c.pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func attemptToCatchPokemon(p Pokemon) bool {
	fmt.Printf("Throwing a Pokeball at %s...\n", p.Name)
	baseExp := max(p.BaseExp, MAX_ATTEMPT)
	maxBex := max(p.BaseExp, MAX_BEX)

	attemps := MIN_ATTEMPT + ((baseExp-MAX_ATTEMPT)/(maxBex-MAX_ATTEMPT))*(MAX_ATTEMPT-MIN_ATTEMPT)
	return rand.Intn(attemps) == 0
}

func displayPokemonListFromCache(data []byte) error {
	var pokemons PokemonList
	err := json.Unmarshal(data, &pokemons)
	if err != nil {
		return err
	}

	for _, p := range pokemons.Pokemons {
		fmt.Printf(" - %s\n", p.Pokemon.Name)
	}

	return nil
}

func displayLocationsFromCache(c *Config, data []byte) error {
	var la LocationAreaResponse
	err := json.Unmarshal(data, &la)
	if err != nil {
		return err
	}
	for _, loc := range la.Results {
		fmt.Println(loc.Name)
	}
	c.setNextUrl(la.NextUrl)
	c.setPrevUrl(la.PrevUrl)
	return nil
}

func makeGetRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("the location area does not exist")
	}
	return res, nil
}
