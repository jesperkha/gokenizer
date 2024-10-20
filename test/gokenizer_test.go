package test

import (
	"slices"
	"strings"
	"testing"

	"github.com/jesperkha/gokenizer"
)

func TestBasicPatterns(t *testing.T) {
	input := "Hello, world!"
	expectedOutput := []string{"Hello", ",", "world", "!"}
	output := []string{}

	tokr := gokenizer.New()

	tokr.Pattern("{word}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	tokr.Pattern("{symbol}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	tokr.Run(input)

	if slices.Compare(output, expectedOutput) != 0 {
		t.Fatalf("expected %s, got %s", strings.Join(expectedOutput, "|"), strings.Join(output, "|"))
	}
}
