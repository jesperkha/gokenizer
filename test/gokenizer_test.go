package test

import (
	"slices"
	"strings"
	"testing"

	"github.com/jesperkha/gokenizer"
)

func TestStaticPattern(t *testing.T) {
	word := "Hello"
	tokr := gokenizer.New()
	result := ""

	tokr.Pattern(word, func(tok gokenizer.Token) error {
		result = tok.Lexeme
		return nil
	})

	if err := tokr.Run(word); err != nil {
		t.Fatal(err)
	}

	if result != word {
		t.Fatalf("expected '%s', got '%s'", word, result)
	}
}

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

	if err := tokr.Run(input); err != nil {
		t.Fatal(err)
	}

	if slices.Compare(output, expectedOutput) != 0 {
		t.Fatalf("expected %s, got %s", strings.Join(expectedOutput, "|"), strings.Join(output, "|"))
	}
}
