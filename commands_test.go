package main

import "testing"

func TestGetCommands(t *testing.T) {
	commands := getCommands()

	response := commands["help"].name
	expectResponse := "help"
	responseDescription := commands["help"].description
	expectResponseDescription := "Displays a help message"

	if response != expectResponse {
		t.Errorf("Resposta esperada '%s' e recebeu '%s'", expectResponse, response)
	}

	if responseDescription != expectResponseDescription {
		t.Errorf("Resposta esperada '%s' e recebeu '%s'", responseDescription, expectResponseDescription)
	}
}
