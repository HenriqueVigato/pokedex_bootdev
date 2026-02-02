package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

func main() {
	const apiPokemonArea = "https://pokeapi.co/api/v2/location-area/"
	const apiPokemons = "https://pokeapi.co/api/v2/pokemon/"
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	cache := pokecache.NewCache(25000 * time.Millisecond)
	configs := &Config{
		Next:     apiPokemonArea,
		Previous: "",
		cache:    cache,
		pokedex:  make(map[string][]byte),
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

		cmd := commands[userCommand[0]]
		switch cmd.name {
		case "help", "exit", "map", "mapb", "pokedex":
			err := cmd.callback(configs, "")
			if err != nil {
				fmt.Println("Erro:", err)
			}
		case "explore":
			if len(userCommand) < 2 {
				fmt.Println("Favor informe uma area a ser explorada")
			} else {
				erro := cmd.callback(configs, apiPokemonArea+userCommand[1])
				if erro != nil {
					fmt.Println("Erro:", erro)
				}
			}
		case "catch":
			if len(userCommand) < 2 {
				fmt.Println("Favor informe um nome de pokemon a ser capturado")
			} else {
				erro := cmd.callback(configs, apiPokemons+userCommand[1])
				if erro != nil {
					fmt.Println("Erro:", erro)
				}
			}
		case "inspect":
			if len(userCommand) < 2 {
				fmt.Println("Favor informe um nome de pokemon a ser especionado")
			} else {
				erro := cmd.callback(configs, userCommand[1])
				if erro != nil {
					fmt.Println("Erro: ", erro)
				}
			}
		default:
			fmt.Println("Comando desconhecido. Digite 'help' para ver os comandos.")
		}
	}
}
