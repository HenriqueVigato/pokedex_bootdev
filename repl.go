package main
import (
	"strings"
)

func cleanInput(text string) []string {
	var result []string
	result = strings.Fields(strings.ToLower(text))
	return result
}


