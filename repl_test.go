package main
import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input string
		expected []string
	}{
		{
			input: "  hello   world   ",
			expected: []string{"hello", "world"},
		},
		{
			input: "ola meus ConteRRANESO",
			expected: []string{"ola", "meus", "conterraneso"},
		},
	}
	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("Os tamnhos dos resultados diferem: got %d, want %d", len(actual), len(c.expected))
			continue
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("As palavras sao diferentes: got %s, want %s", word, expectedWord)
				continue
			}
		}
	}
}

