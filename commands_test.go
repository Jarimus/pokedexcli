package main

import "testing"

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  Hello  World  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    " tHiS is foR You  ",
			expected: []string{"this", "is", "for", "you"},
		},
	}

	for _, c := range cases {
		actual := cleanInputString(c.input)

		// test output length
		if len(actual) != len(c.expected) {
			t.Errorf("Test FAIL!\nInput: %s\nWrong number of elements in the cleaned input. %d instead of %d.", c.input, len(actual), len(c.expected))
		}

		// test output words
		for i, word := range actual {
			if c.expected[i] != word {
				t.Errorf("Test FAIL!\nInput: %s\n%s instead of %s at index %d", c.input, word, c.expected[i], i)
			}
		}
	}
}
