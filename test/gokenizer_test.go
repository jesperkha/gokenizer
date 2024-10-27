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

	tokr.Pattern("foo{barl", func(t gokenizer.Token) error {
		return nil
	})

	if err := tokr.Run(""); err == nil {
		t.Fatal("expected error, got nil")
	}

	tokr.Pattern("}", func(t gokenizer.Token) error {
		return nil
	})

	if err := tokr.Run(""); err == nil {
		t.Fatal("expected error, got nil")
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

func TestTokenValues(t *testing.T) {
	tokr := gokenizer.New()

	input := "foo bar faz"
	expect := []string{"0", "4", "8"}
	output := []string{}

	tokr.Pattern("{word}", func(t gokenizer.Token) error {
		output = append(output, fmt.Sprint(t.Pos))
		if t.Source != input {
			return fmt.Errorf("expected source '%s', got '%s'", input, t.Source)
		}
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Fatal(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestValuesMap(t *testing.T) {
	input := "Hello, world!"
	tokr := gokenizer.New()

	testGet := func(t *testing.T, tok gokenizer.Token, class, expect string, idx int) {
		if w := tok.GetAt(class, idx); w.Lexeme != expect {
			t.Errorf("expected '%s', got '%s'", expect, w.Lexeme)
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

func TestUserClass(t *testing.T) {
	input := "a+b-c=3"
	expect := "a+b"

	tokr := gokenizer.New()

	tokr.Class("math", func(b byte) bool {
		return strings.Contains("+-/*=", string(b))
	})

	tokr.Pattern("a{math}c", func(tok gokenizer.Token) error {
		if tok.Lexeme != expect {
			return fmt.Errorf("expected '%s', got '%s'", expect, tok.Lexeme)
		}
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}
}

func TestUserPatternClass(t *testing.T) {
	input := "var foo = 123;"
	expect := strings.Split(input, "\n")
	output := []string{}

	tokr := gokenizer.New()

	tokr.ClassPattern("variable", "{word}")
	tokr.ClassPattern("onetwothree", "123")
	tokr.ClassPattern("declaration", "var {variable} = {onetwothree}")

	tokr.Pattern("{declaration};", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)

		if t.Lexeme != input {
			return fmt.Errorf("expected '%s', got '%s'", input, t.Lexeme)
		}
		if expect, got := input[:len(input)-1], t.Get("declaration"); got.Lexeme != expect {
			return fmt.Errorf("expected '%s', got '%s'", expect, got.Lexeme)
		}

		// Nested values
		if expect, got := "foo", t.Get("declaration").Get("variable"); got.Lexeme != expect {
			return fmt.Errorf("expected '%s', got '%s'", expect, got.Lexeme)
		}
		if expect, got := "123", t.Get("declaration").Get("onetwothree"); got.Lexeme != expect {
			return fmt.Errorf("expected '%s', got '%s'", expect, got.Lexeme)
		}
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestNestedParsing(t *testing.T) {
	input := "1,2,3\na,b,c\n7,8,9"
	expect := strings.Split(input, "\n")
	output := []string{}

	tokr := gokenizer.New()

	tokr.Pattern("{line}", func(t gokenizer.Token) error {
		return tokr.Run(t.Get("line").Lexeme)
	})

	tokr.Pattern("{number},{number},{number}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	tokr.Pattern("{word},{word},{word}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestClassAny(t *testing.T) {
	input := "123 foo! hello"
	expect := []string{"123", "foo!", "hello"}
	output := []string{}

	tokr := gokenizer.New()

	tokr.ClassAny("any", "{number}", "{word}{symbol}", "hello")

	tokr.Pattern("{any}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		if w := t.Get("any").Get("number").Lexeme; w != "" && w != "123" {
			return fmt.Errorf("expected '%s', got '%s'", "123", w)
		}
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestEmptyPattern(t *testing.T) {
	input := "a= b"
	expect := []string{"a= b"}
	output := []string{}

	tokr := gokenizer.New()

	// Whitespace
	tokr.ClassAny("ws", " ", "")

	tokr.ClassPattern("foo", "{word}{ws}{symbol}{ws}{word}")

	tokr.Pattern("{foo}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestDeepNesting(t *testing.T) {
	input := ""
	expect := []string{}
	output := []string{}

	tokr := gokenizer.New()

	tokr.Pattern("", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestClassOptional(t *testing.T) {
	input := "foo =bar;"
	expect := []string{input}
	output := []string{}

	tokr := gokenizer.New()

	tokr.ClassOptional("semicolon?", ";")
	tokr.ClassOptional("space?", " ")

	tokr.Pattern("{word}{space?}={space?}{word}{semicolon?}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Error(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}
