package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		success := scanner.Scan()
		if !success {
			fmt.Fprintln(os.Stderr, "Wasn't possible read a message")
			break
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "shouldn't see an error scannig a string")
		}

		message := scanner.Text()
		cleanedMessage := cleanInput(message)

		if len(cleanedMessage) == 0 {
			fmt.Println("Comando vazio")
			continue
		}

		firstMessage := cleanedMessage[0]

		fmt.Println("Your command was: ", firstMessage)
	}
}
