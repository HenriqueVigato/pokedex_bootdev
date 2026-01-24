package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

func main() {
	const pokemonArea = "https://pokeapi.co/api/v2/location-area/"
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cache := pokecache.NewCache(25000 * time.Millisecond)
	configs := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
		pokedex:  make(map[string]any),
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
			if len(userCommand) > 1 && userCommand[0] == "explore" {
				erro := cmd.callback(configs, pokemonArea+userCommand[1])

				if erro != nil {
					fmt.Println("Erro:", erro)
				}
			} else {
				if userCommand[0] == "explore" {
					fmt.Println("Favor informe uma area.")
				} else {
					err := cmd.callback(configs, "")
					if err != nil {
						fmt.Println("Erro:", err)
					}
				}
			}
		} else {
			fmt.Println("Comando desconhecido. Digite 'help' para ver os comandos.")
		}
	}
}
