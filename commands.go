package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, string) error
}

type Config struct {
	Next     string
	Previous string
	cache    *pokecache.Cache
	pokedex  map[string][]byte
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
		"explore": {
			name:        "explore",
			description: "List of all the Pokemon located",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Try capture a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Gives detailed information of pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all the captured pokemon",
			callback:    commandPokedex,
		},
	}
	return commands
}

func commandExit(i *Config, n string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *Config, n string) error {
	command := getCommands()
	fmt.Printf("Comandos disponiveis: \n\n")
	for _, v := range command {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(c *Config, n string) error {
	result, ok := c.cache.Get(c.Next)
	if !ok {
		locations, err := getData(c.Next)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		c.cache.Add(c.Next, locations)
		updateNextPrevius(c, convertToJSON(locations))
		printMap(convertToJSON(locations)["results"].([]any))

	} else {
		jsonResult := convertToJSON(result)

		updateNextPrevius(c, jsonResult)
		printMap(jsonResult["results"].([]any))
	}
	return nil
}

func commandMapb(c *Config, n string) error {
	if c.Previous == "" {
		fmt.Println("There are not previous page")
		return nil
	}
	result, ok := c.cache.Get(c.Previous)
	if !ok {
		locations, err := getData(c.Previous)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		c.cache.Add(c.Previous, locations)
		updateNextPrevius(c, convertToJSON(locations))
		printMap(convertToJSON(locations)["results"].([]any))

	} else {
		jsonResult := convertToJSON(result)

		updateNextPrevius(c, jsonResult)
		printMap(jsonResult["results"].([]any))
	}
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

func commandExplore(c *Config, area string) error {
	result, ok := c.cache.Get(area)
	if !ok {
		locations, err := getData(area)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		c.cache.Add(area, locations)
		printPokemons(convertToJSON(locations)["pokemon_encounters"].([]any))
	} else {
		jsonResult := convertToJSON(result)
		printPokemons(jsonResult["pokemon_encounters"].([]any))
	}

	return nil
}

func printPokemons(locations []any) {
	for _, v := range locations {
		fmt.Println(v.(map[string]any)["pokemon"].(map[string]any)["name"].(string))
	}
}

func commandCatch(c *Config, pokemon string) error {
	pokemonName := path.Base(pokemon)
	var pokeDataJSON map[string]any
	var pokeByte []byte

	fmt.Printf("Throwing a Pokeball at %s... \n", pokemonName)
	result, ok := c.cache.Get(pokemon)

	if !ok {
		pokeData, err := getData(pokemon)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
		c.cache.Add(pokemon, pokeData)
		pokeByte = pokeData
		pokeDataJSON = convertToJSON(pokeData)
	} else {
		pokeDataJSON = convertToJSON(result)
		pokeByte = result
	}
	if tryCatchPokemon(int(pokeDataJSON["base_experience"].(float64))) {
		c.pokedex[pokemonName] = pokeByte
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func tryCatchPokemon(baseExperience int) bool {
	probability := baseExperience % 11
	roll := rand.Intn(10)

	if probability == roll {
		return true
	} else {
		return false
	}
}

func commandInspect(c *Config, pokemon string) error {
	pokemonData, ok := c.pokedex[pokemon]
	if !ok {
		fmt.Println("favor informe um pokemon que ja tenha sido capturado")
		return fmt.Errorf("Pokemon not found")
	}
	pokeStruck, err := ConvertToStruct(pokemonData)
	if err != nil {
		return fmt.Errorf("erro: %v", err)
	}
	fmt.Println(printPokemonStats(pokeStruck))

	return nil
}

func printPokemonStats(p *Pokemon) string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Name: %s\n", p.Forms[0].Name)
	fmt.Fprintf(&sb, "Height: %d\n", p.Heigth)
	fmt.Fprintf(&sb, "Weight: %d\n", p.Weight)

	sb.WriteString("Stats: \n")

	for _, stat := range p.Stats {
		fmt.Fprintf(&sb, "  -%s: ", stat.Stat.Name)
		fmt.Fprintf(&sb, "%d\n", stat.BaseStat)
	}

	sb.WriteString("Types:\n")

	for _, tipe := range p.Types {
		fmt.Fprintf(&sb, "  - %s", tipe.Type.Name)
	}

	return sb.String()
}

func commandPokedex(c *Config, n string) error {
	if len(c.pokedex) <= 0 {
		fmt.Println("Voce nao capturou nenhum pokemon ainda")
	} else {
		fmt.Println("Your Pokedex:")
		for key := range c.pokedex {
			fmt.Printf(" - %s\n", key)
		}
	}
	return nil
}
