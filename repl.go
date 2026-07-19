package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	nextUrl string
	prevUrl string
}

type locationAreaResp struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
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
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	return strings.Fields(lowerText)
}

func commandExit(*config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config) error {
	fmt.Println("")
	return nil
}

func commandMap(conf *config) error {

	if conf.nextUrl == "" {
		conf.nextUrl = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(conf.nextUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return fmt.Errorf("bad statuscode: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locationAreaJson locationAreaResp
	err = json.Unmarshal(byteData, &locationAreaJson)
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

func commandMapb(conf *config) error {

	if conf.prevUrl == "" {
		fmt.Println("first, use command map at least twice")
		return nil
	}

	res, err := http.Get(conf.prevUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode > 299 || res.StatusCode < 200 {
		return fmt.Errorf("bad statuscode: %d", res.StatusCode)
	}

	byteData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var locationAreaJson locationAreaResp
	err = json.Unmarshal(byteData, &locationAreaJson)
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
