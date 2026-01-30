package main

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	pokecache "github.com/HenriqueVigato/pokedex_bootdev/internal"
)

const pokemonArea = "https://pokeapi.co/api/v2/location-area/"

var (
	cache   = pokecache.NewCache(25000 * time.Millisecond)
	configs = &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
	}
)

func resetData() {
	cache = pokecache.NewCache(25000 * time.Millisecond)
	configs = &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
	}
}

func capturaOutput(commands map[string]cliCommand, command, moreInfo string) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	commands[command].callback(configs, moreInfo)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	return output
}

func TestHelpCommands(t *testing.T) {
	resetData()
	commands := getCommands()

	output := capturaOutput(commands, "help", "")

	if !strings.Contains(output, "Displays the names of locations") {
		t.Errorf("Esperava encontar 'Displays the names of locations' dentre as respostas")
	}
}

func TestMapCommands(t *testing.T) {
	commands := getCommands()
	output := capturaOutput(commands, "map", "")

	if !strings.Contains(output, "canalave-city-area") {
		t.Errorf("Esperava encontrar 'Canalave-city-area, mas nao foi encontrado")
	}
}

func TestMapbCommands(t *testing.T) {
	resetData()
	commands := getCommands()
	_ = capturaOutput(commands, "map", "")
	_ = capturaOutput(commands, "map", "")
	mapWasCalled := strings.Contains(capturaOutput(commands, "map", ""), "ravaged-path-area")

	if !mapWasCalled {
		t.Errorf("Erro em chamara a funcao map para preparar o teste de mapb")
	}

	output := capturaOutput(commands, "mapb", "")

	if !strings.Contains(output, "mt-coronet-b1f") {
		t.Errorf("Esperava encontrar a area de mt-coronet-b1f")
	}
}

func TestExploreCommands(t *testing.T) {
	urlArea := "https://pokeapi.co/api/v2/location-area/ravaged-path-area"
	resetData()
	commands := getCommands()
	output := capturaOutput(commands, "explore", urlArea)

	if !strings.Contains(output, "zubat") {
		t.Errorf("Esperava encontrar o pokemon Zubat")
	}
}

func TestCatchCommands(t *testing.T) {
	resetData()
	commands := getCommands()
	escaped := 0
	captured := 0

	for i := 0; i < 10; i++ {
		output := capturaOutput(commands, "catch", "pikachu")

		if strings.Contains(output, "escaped") {
			escaped++
		}
		if strings.Contains(output, "caught") {
			captured++
		}
	}
	if captured <= 0 {
		t.Errorf("Deveria capturar pelo menos 1 pokemon em 10 tentativas")
	}
	if escaped <= 0 {
		t.Errorf("Deveria ter pelo menos um que conseguiu escapar")
	}
}
