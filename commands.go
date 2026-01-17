package main

import (
	"fmt"
	"os"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
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

func commandMap(c *Config) error {
	locations, err := convertToJSON(getData(c.Next))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	updateNextPrevius(c, locations)
	printMap(locations["results"].([]any))

	return nil
}

func commandMapb(c *Config) error {
	if c.Previous == "" {
		fmt.Println("There are not previous page")
		return nil
	}

	locations, err := convertToJSON(getData(c.Previous))
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	updateNextPrevius(c, locations)
	printMap(locations["results"].([]any))

	return nil
}

func printMap(locations []any) {
	for _, v := range locations {
		fmt.Println(v.(map[string]any)["name"])
	}
}

func updateNextPrevius(c *Config, locations map[string]any) {
	c.Next = locations["next"].(string)
	if locations["previous"] == nil {
		c.Previous = ""
	} else {
		c.Previous = locations["previous"].(string)
	}
}
