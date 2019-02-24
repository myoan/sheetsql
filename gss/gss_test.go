package gss

import (
	"testing"
)

func TestColumnToRange(t *testing.T) {
	validData := []struct {
		input1 []string
		input2 []string
		output []int
	}{
		{input1: []string{"a", "b", "c", "d", "e"}, input2: []string{"a", "b"}, output: []int{1, 2}},
		{input1: []string{"a", "b", "c", "d", "e"}, input2: []string{"a", "d"}, output: []int{1, 4}},
		{input1: []string{"a", "b", "c", "d", "e"}, input2: []string{"d", "a"}, output: []int{1, 4}},
		{input1: []string{"a", "b", "c", "d", "e"}, input2: []string{"d", "a", "c"}, output: []int{1, 4}},
		{input1: []string{"a", "b", "c", "d", "e"}, input2: []string{"*"}, output: []int{1, 5}},
		{input1: []string{"a", "c", "d"}, input2: []string{"d", "a", "c"}, output: []int{1, 3}},
		{input1: []string{"a", "c"}, input2: []string{"d", "a", "c"}, output: []int{1, 2}},
	}

	for _, tcase := range validData {
		min, max := columnToRange(tcase.input1, tcase.input2)
		if min != tcase.output[0] || max != tcase.output[1] {
			t.Errorf("columnToRange(%q, %q) = (%d, %d), want: %v", tcase.input1, tcase.input2, min, max, tcase.output)
		}
	}
}

func TestColumnAlphabet(t *testing.T) {
	validData := []struct {
		input  int
		output string
	}{
		{input: 0, output: ""},
		{input: 1, output: "A"},
		{input: 26, output: "Z"},
		{input: 27, output: "AA"},
		{input: 53, output: "BA"},
	}

	for _, tcase := range validData {
		out, _ := columnAlphabet(tcase.input)
		if out != tcase.output {
			t.Errorf("columnAlphabet(%q) = %q, want: %q", tcase.input, out, tcase.output)
		}
	}
}
