package main

import (
	"os"
	"testing"
)

func TestStructConversion(t *testing.T) {
	data, err := os.ReadFile("./testsData/pikachu.json")
	if err != nil {
		t.Fatal(err)
	}
	pokeStruct, erro := ConvertToStruct(data)
	if erro != nil {
		t.Errorf("%v", erro)
	}

	if pokeStruct.Forms[0].Name != "pikachu" {
		t.Errorf("Esperava o nome do pokemon no caso pikachu")
	}
}
