package main

import (
	"testing"
)

func TestApi(t *testing.T) {
	data, err := getData("https://pokeapi.co/api/v2/location/1")
	if err != nil {
		t.Errorf("Nao era eperado nenhum tipo de erro %v", err)
	}
	response := convertToJSON(data)

	word := response["areas"].([]any)[0].(map[string]any)["name"].(string)
	expectedWord := "canalave-city-area"
	if word != expectedWord {
		t.Errorf("O resultado nao foi o esperado era %s e recebeu %s", word, expectedWord)
	}
}
