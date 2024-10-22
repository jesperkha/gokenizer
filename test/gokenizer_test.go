package test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jesperkha/gokenizer"
)

func TestStaticPattern(t *testing.T) {
	word := "golang"
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

func TestClassParser(t *testing.T) {
	tokr := gokenizer.New()

	tokr.Pattern("foo{word}a{number}", func(t gokenizer.Token) error {
		return nil
	})

	if err := tokr.Run(""); err != nil {
		t.Fatal(err)
	}
}

func makeClassTester(className, input, expected string) func(*testing.T) {
	return func(t *testing.T) {
		tokr := gokenizer.New()
		word := ""

		tokr.Pattern(fmt.Sprintf("{%s}", className), func(tok gokenizer.Token) error {
			word = tok.Lexeme
			return nil
		})

		if err := tokr.Run(input); err != nil {
			t.Error(err)
		}

		if word != expected {
			t.Errorf("expected '%s', got '%s'", expected, word)
		}
	}
}

func TestClasses(t *testing.T) {
	input := "golang123!"

	tests := []func(*testing.T){
		makeClassTester("word", input, "golang"),
		makeClassTester("number", input, "123"),
		makeClassTester("symbol", input, "!"),
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i+1), tt)
	}
}

func makeTokenizerTester(input string, expected, patterns []string) func(t *testing.T) {
	return func(t *testing.T) {
		tokr := gokenizer.New()
		output := []string{}

		for _, p := range patterns {
			tokr.Pattern(p, func(t gokenizer.Token) error {
				output = append(output, t.Lexeme)
				return nil
			})
		}

		if err := tokr.Run(input); err != nil {
			t.Error(err)
		}

		if slices.Compare(output, expected) != 0 {
			t.Errorf("expected '%s', got '%s'", strings.Join(expected, "|"), strings.Join(output, "|"))
		}
	}
}

func TestTokenizer(t *testing.T) {
	tests := []func(*testing.T){
		// Basic test
		makeTokenizerTester(
			"Hello, world!",
			[]string{"Hello", ",", "world", "!"},
			[]string{"{word}", "{symbol}"},
		),
		// Multiple symbol test
		makeTokenizerTester(
			"a != b",
			[]string{"!", "="},
			[]string{"{symbol}"},
		),
		// Mix static and class
		makeTokenizerTester(
			"aQ:foo?!",
			[]string{"Q:foo?"},
			[]string{"Q{symbol}{word}?"},
		),
		// Test with braces
		makeTokenizerTester(
			"a{foo}",
			[]string{"{foo}"},
			[]string{"{lbrace}{word}{rbrace}"},
		),
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case_%d", i+1), tt)
	}
}

func TestValuesMap(t *testing.T) {
	input := "Hello, world!"
	tokr := gokenizer.New()

	testGet := func(t *testing.T, tok gokenizer.Token, class, expect string, idx int) {
		if w := tok.GetAt(class, idx); w != expect {
			t.Errorf("expected '%s', got '%s'", expect, w)
		}
	}

	tokr.Pattern("{word}{symbol} {word}{symbol}", func(tok gokenizer.Token) error {
		testGet(t, tok, "word", "Hello", 0)
		testGet(t, tok, "word", "world", 1)
		testGet(t, tok, "symbol", ",", 0)
		testGet(t, tok, "symbol", "!", 1)

		testGet(t, tok, "number", "", 0)
		testGet(t, tok, "word", "", 99)

		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Fatal(err)
	}
}
