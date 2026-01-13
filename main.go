package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cache := pokecache.NewCache(10 * time.Millisecond)
	configs := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
	}

	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("")
	fmt.Println("help to see the commands")

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
		userCommand := cleanInput(message)

		if len(userCommand) == 0 {
			fmt.Println("Comando vazio")
			continue
		}

		if cmd, exist := commands[userCommand[0]]; exist {

			err := cmd.callback(configs)
			if err != nil {
				fmt.Println("Erro:", err)
			}
		} else {
			fmt.Println("Comando desconhecido. Digite 'help' para ver os comandos.")
		}

	}
}
