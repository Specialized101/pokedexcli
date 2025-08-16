package main

const LOCATION_AREA_URL = "https://pokeapi.co/api/v2/location-area/"

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
