package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		enteredText := scanner.Text()
		formattedText := strings.ToLower(enteredText)
		words := strings.Fields(formattedText)
		firstWord := words[0]
		fmt.Printf("Your command was: %s \n", firstWord)
	}

}
