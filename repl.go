package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/govaeka/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

type config struct {
	nextUrl string
	prevUrl string
	cache   *pokecache.Cache
	caught  map[string]pokecache.Pokemon
}

type locationAreaResp struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name              string `json:"name"`
	URL               string `json:"url"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

var commandLookUp = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Display a help message",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Display next 20 location areas",
		callback:    commandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Display previous 20 location areas",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Display the Pokemon in the area",
		callback:    commandExplore,
	},
	"catch": {
		name:        "catch <pokemon_name>",
		description: "Attempt to catch a pokemon",
		callback:    commandCatch,
	},
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	return strings.Fields(lowerText)
}

func commandCatch(conf *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide a pokemon name")
	}
	name := args[0]
	fmt.Printf("Throwing a Pokeball at %s...", name)

	return nil
}

func commandExit(conf *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandExplore(conf *config, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("please provide an area name")
	}

	area := args[0]
	fmt.Printf("Exploring %s...\n", area)

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", area)

	slice, found := conf.cache.Get(url)

	var byteData []byte

	if found {
		byteData = slice
	} else {
		res, err := http.Get(url)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		byteData, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		conf.cache.Add(url, byteData)
	}

	var jsonData locationArea
	err := json.Unmarshal(byteData, &jsonData)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range jsonData.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandHelp(conf *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")

	return nil
}

func commandMap(conf *config, args []string) error {

	if conf.nextUrl == "" {
		conf.nextUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	slice, found := conf.cache.Get(conf.nextUrl)

	var byteData []byte

	if found {
		byteData = slice
	} else {
		res, err := http.Get(conf.nextUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 || res.StatusCode < 200 {
			return fmt.Errorf("bad statuscode: %d", res.StatusCode)
		}

		byteData, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		conf.cache.Add(conf.nextUrl, byteData)
	}

	var locationAreaJson locationAreaResp
	err := json.Unmarshal(byteData, &locationAreaJson)
	if err != nil {
		return err
	}

	conf.nextUrl = locationAreaJson.Next
	conf.prevUrl = locationAreaJson.Previous

	for _, area := range locationAreaJson.Results {
		fmt.Println(area.Name)
	}

	return nil
}

func commandMapb(conf *config, args []string) error {

	if conf.prevUrl == "" {
		fmt.Println("first, use command map at least twice")
		return nil
	}

	slice, found := conf.cache.Get(conf.prevUrl)

	var byteData []byte

	if found {
		byteData = slice
	} else {
		res, err := http.Get(conf.prevUrl)
		if err != nil {
			return err
		}
		defer res.Body.Close()

		if res.StatusCode > 299 || res.StatusCode < 200 {
			return fmt.Errorf("bad statuscode: %d", res.StatusCode)
		}

		byteData, err = io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		conf.cache.Add(conf.prevUrl, byteData)
	}

	var locationAreaJson locationAreaResp
	err := json.Unmarshal(byteData, &locationAreaJson)
	if err != nil {
		return err
	}

	conf.nextUrl = locationAreaJson.Next
	conf.prevUrl = locationAreaJson.Previous

	for _, area := range locationAreaJson.Results {
		fmt.Println(area.Name)
	}

	return nil
}

// test poke
