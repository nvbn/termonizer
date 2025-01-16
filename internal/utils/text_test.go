package utils

import (
	"fmt"
	"testing"
)

func TestFindLineStart(t *testing.T) {
	type testData struct {
		text     string
		pos      int
		expected int
	}

	inputsExpecteds := []testData{
		{"", 0, 0},
		{"* first task", 7, 0},
		{"* first task\n* second task", 7, 0},
		{"* first task\n* second task", 13, 13},
		{"* first task\n* second task", 23, 13},
		{"* first task\n* second task", 35, 13},
	}

	for _, inputExpected := range inputsExpecteds {
		t.Run(fmt.Sprintf("%+v", inputExpected), func(t *testing.T) {
			actual := FindLineStart(inputExpected.text, inputExpected.pos)
			if actual != inputExpected.expected {
				t.Errorf("got %d, want %d", actual, inputExpected.expected)
			}
		})
	}
}

func TestFindLineEnd(t *testing.T) {
	type testData struct {
		text     string
		pos      int
		expected int
	}

	inputsExpecteds := []testData{
		{"", 0, 0},
		{"* first task", 7, 11},
		{"* first task\n* second task", 7, 12},
		{"* first task\n* second task", 13, 25},
		{"* first task\n* second task", 23, 25},
		{"* first task\n* second task", 35, 25},
		{"* first task\n* second task\n* third task", 13, 26},
	}

	for _, inputExpected := range inputsExpecteds {
		t.Run(fmt.Sprintf("%+v", inputExpected), func(t *testing.T) {
			actual := FindLineEnd(inputExpected.text, inputExpected.pos)
			if actual != inputExpected.expected {
				t.Errorf("got %d, want %d", actual, inputExpected.expected)
			}
		})
	}
}
