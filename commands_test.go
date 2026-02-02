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

var (
	cache   = pokecache.NewCache(25000 * time.Millisecond)
	configs = &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
		pokedex:  make(map[string][]byte),
	}
)

func resetData() {
	cache = pokecache.NewCache(25000 * time.Millisecond)
	configs = &Config{
		Next:     "https://pokeapi.co/api/v2/location-area/",
		Previous: "",
		cache:    cache,
		pokedex:  make(map[string][]byte),
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
	expectedCommands := []string{"exit", "help", "map", "mapb", "explore", "catch", "inspect"}

	if !strings.Contains(output, "Displays the names of locations") {
		t.Errorf("Esperava encontar 'Displays the names of locations' dentre as respostas")
	}

	for _, command := range expectedCommands {
		if !strings.Contains(output, command) {
			t.Errorf("Esperava encontrar o comando: %s", command)
		}
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

	for range int(100) {
		output := capturaOutput(commands, "catch", "https://pokeapi.co/api/v2/pokemon/pikachu")

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
	if _, exist := configs.pokedex["pikachu"]; !exist {
		t.Errorf("Pikachu nao consta na pokedex")
	} else {
		catched := convertToJSON(configs.pokedex["pikachu"])
		t.Logf("Pikachu consta na pokedex: %v", catched["forms"].([]any)[0].(map[string]any)["name"].(string))
	}
}

func TestInspectCommands(t *testing.T) {
	commands := getCommands()
	output := capturaOutput(commands, "inspect", "pikachu")
	wrongOutput := capturaOutput(commands, "inspect", "mew")

	if !strings.Contains(wrongOutput, "tenha sido capturado") {
		t.Logf("Resposta obtida: %s", wrongOutput)
		t.Errorf("Deveria constar um mensagem para digitar nome de apenas pokemons capturados")
	}

	if !strings.Contains(output, "Name: pikachu") {
		t.Logf("Output: \n%s", output)
		t.Errorf("Deveria  conter o nome do pokemon no output")
	}

	if !strings.Contains(output, "speed: 90") {
		t.Logf("Output: %s", output)
		t.Errorf("Deveria conter no output a velocidade do pokemon")
	}
}

func TestPokedexCommands(t *testing.T) {
	commands := getCommands()
	output := capturaOutput(commands, "pokedex", "")

	if !strings.Contains(output, "pikachu") {
		t.Logf("Output: %s", output)
		t.Errorf("Deveria listar os pokemons capturados")
	}

	resetData()

	cleanOutput := capturaOutput(commands, "pokedex", "")

	if !strings.Contains(cleanOutput, "nao capturou nenhum pokemon ainda") {
		t.Logf("Output: \n %s", cleanOutput)
		t.Errorf("Deveria exibir uma mensagem de erro que nao tem nenhum pokemon na pokedex")
	}
}
