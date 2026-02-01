package main

import (
	"encoding/json"
	"fmt"
)

type Pokemon struct {
	Forms []struct {
		Name string `json:"name"`
	}
	Heigth int `json:"height"`
	Weight int `json:"weight"`
	Stats  []struct {
		BaseStat int `json:"base_stat"`
		Stat     struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}

func ConvertToStruct(pokeData []byte) (*Pokemon, error) {
	var pokemon Pokemon
	if erro := json.Unmarshal(pokeData, &pokemon); erro != nil {
		fmt.Println(erro)
		return nil, fmt.Errorf("failed to Unmarshal pokemon data %w", erro)
	}

	return &pokemon, nil
}
