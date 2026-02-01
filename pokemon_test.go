package main

import (
	"encoding/json"
	"os"
	"testing"
)

func TestStructConversion(t *testing.T) {
	data, err := os.ReadFile("./testsData/pikachu.json")
	if err != nil {
		t.Fatal(err)
	}
	var pokemon Pokemon
	err = json.Unmarshal(data, &pokemon)
	if err != nil {
		t.Fatal(err)
	}

	if pokemon.Forms[0].Name != "pikachu" {
		t.Errorf("Esperava o nome do pokemon no caso pikachu")
	}
}
