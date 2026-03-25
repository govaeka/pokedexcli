package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Welcome to the Pokedex!")
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		enteredText := scanner.Text()
		words := cleanInput(enteredText)
		if len(words) == 0 {
			continue
		}
		firstWord := words[0]
		// onderstaand blok ivm maps: theorie herhalen
		if cmd, ok := commandLookUp[firstWord]; ok {
			err := cmd.callback()
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}

}
