package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/govaeka/pokedexcli/internal/pokecache"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)
	var conf config
	conf.cache = pokecache.NewCache(5 * time.Minute)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		enteredText := scanner.Text()
		words := cleanInput(enteredText)
		if len(words) == 0 {
			continue
		}
		firstWord := words[0]
		args := words[1:]
		// onderstaand blok ivm maps: theorie herhalen
		if cmd, ok := commandLookUp[firstWord]; ok {
			err := cmd.callback(&conf, args)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}

}
