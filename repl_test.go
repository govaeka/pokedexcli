package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			input:    " hello world ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO   WORLD",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		exp := c.expected
		if len(actual) != len(exp) {
			t.Errorf("Failed: Expected length %d but got %d", len(exp), len(actual))
		}

		for i := range actual {
			word := actual[i]
			expWord := exp[i]
			if word != expWord {
				t.Errorf("Failed: Expected %s but got %s", expWord, word)
			}
		}
	}

}
