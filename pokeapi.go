package main

const LOCATION_AREA_URL = "https://pokeapi.co/api/v2/location-area"
const POKEMON_URL = "https://pokeapi.co/api/v2/pokemon"
const MIN_ATTEMPT = 2
const MAX_ATTEMPT = 50
const MAX_BEX = 300

type LocationAreaResponse struct {
	Count   int            `json:"count"`
	NextUrl string         `json:"next"`
	PrevUrl string         `json:"previous"`
	Results []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type PokemonList struct {
	Pokemons []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Name    string `json:"name"`
	Height  int    `json:"height"`
	Weight  int    `json:"weight"`
	BaseExp int    `json:"base_experience"`
	Stats   []struct {
		Base     int `json:"base_stat"`
		StatData struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
