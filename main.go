package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the names of locations areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the Previous page of map",
			callback:    commandMapb,
		},
	}
	return commands
}

func commandExit(*Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config) error {
	command := getCommands()
	fmt.Printf("Comandos disponiveis: \n\n")
	for _, v := range command {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func getDataJSON(url string) (map[string]any, error) {
	res, err := http.Get(url)
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("response failed with a StatusCode: %v", res.StatusCode)
	}
	if err != nil {
		return nil, fmt.Errorf("error Creating a request to the api %v", err)
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("error readind the body response %v", err)
	}

	var dataJSON map[string]any
	if err := json.Unmarshal(body, &dataJSON); err != nil {
		return nil, fmt.Errorf("error during Unmarshal the response body %v", err)
	}

	return dataJSON, nil
}

func commandMap(c *Config) error {
	locations, err := getDataJSON(c.Next)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	c.Next = locations["next"].(string)
	if locations["Previous"] == nil {
		c.Previous = ""
	} else {
		c.Previous = locations["previous"].(string)
	}

	for _, v := range locations["results"].([]any) {
		fmt.Println(v.(map[string]any)["name"])
	}

	return nil
}

func commandMapb(c *Config) error {
	var locations map[string]any

	if locations["Previous"] == nil {
		fmt.Println("You are in the first page!!")
		return nil
	}

	return nil
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	commands := getCommands()
	configs := &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
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
