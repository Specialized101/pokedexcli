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
			input:    "   hello  world   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "I Use Vim BTW",
			expected: []string{"i", "use", "vim", "btw"},
		},
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "Pikachu has 123 ATK Power",
			expected: []string{"pikachu", "has", "123", "atk", "power"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		lenActual, lenExpected := len(actual), len(c.expected)
		if lenActual != lenExpected {
			t.Errorf("input slice size missmatch, expected %v but received %v", lenExpected, actual)
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("wrong output, expected %v but received %v", expectedWord, word)
				t.Fail()
			}
		}
	}
}
