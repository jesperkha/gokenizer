package test

import (
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/jesperkha/gokenizer"
)

func TestInvalidUse(t *testing.T) {
	tokr := gokenizer.New()
	if err := tokr.Run(""); err != nil {
		t.Error(err)
	}

	// no pattern
	tokr.Pattern("", func(t gokenizer.Token) error {
		return nil
	})

	// no callback
	tokr.Pattern("foo", nil)

	// invalid class name
	tokr.Class("{foo}", "bar")

	if err := tokr.Run("123"); err == nil {
		t.Error("expected error")
	}
}

func TestEmpty(t *testing.T) {
	tokr := gokenizer.New()
	if err := tokr.Run(""); err != nil {
		t.Error(err)
	}
}

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

func makeClassTester(className string, inputs, expected []string) func(*testing.T) {
	return func(t *testing.T) {
		tokr := gokenizer.New()
		word := ""

		tokr.Pattern(fmt.Sprintf("{%s}", className), func(tok gokenizer.Token) error {
			word = tok.Lexeme
			return nil
		})

		for i, input := range inputs {
			if err := tokr.Run(input); err != nil {
				t.Error(err)
			}

			if word != expected[i] {
				t.Errorf("expected '%s', got '%s'", expected[i], word)
			}

			word = ""
		}
	}
}

func TestClasses(t *testing.T) {
	tests := []func(*testing.T){
		makeClassTester(
			"word",
			[]string{"foo", "foo bar", "123foo!"},
			[]string{"foo", "bar", "foo"},
		),
		makeClassTester(
			"var",
			[]string{"$foo", "foo_bar"},
			[]string{"$foo", "foo_bar"},
		),
		makeClassTester(
			"number",
			[]string{"1", "1234", "123foo!"},
			[]string{"1", "1234", "123"},
		),
		makeClassTester(
			"string",
			[]string{"\"hello\"", "\"foo", "foo\"", "foo\"bar\"faz"},
			[]string{"\"hello\"", "", "", "\"bar\""},
		),
		makeClassTester(
			"hex",
			[]string{"abc", "#FF01AB", "golang"},
			[]string{"abc", "#FF01AB", "a"},
		),
		makeClassTester(
			"base64",
			[]string{"(aGVsbG8gd29ybGQ=)"},
			[]string{"aGVsbG8gd29ybGQ="},
		),
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

	tokr.ClassFunc("math", func(b byte) bool {
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

	tokr.Class("variable", "{word}")
	tokr.Class("onetwothree", "123")
	tokr.Class("declaration", "var {variable} = {onetwothree}")

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

	tokr.Class("some", "{number}", "{word}{symbol}", "hello")

	tokr.Pattern("{some}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		if w := t.Get("some").Get("number").Lexeme; w != "" && w != "123" {
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

	tokr.Class("foo", "{word}{ws}{symbol}{ws}{word}")

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

func TestClassOptional(t *testing.T) {
	input := "foo =bar;"
	expect := []string{input}
	output := []string{}

	tokr := gokenizer.New()

	tokr.ClassOptional("semicolon", ";")
	tokr.ClassOptional("space", " ")

	tokr.Pattern("{word}{space}={space}{word}{semicolon}", func(t gokenizer.Token) error {
		output = append(output, t.Lexeme)
		return nil
	})

	if err := tokr.Run(input); err != nil {
		t.Fatal(err)
	}

	if slices.Compare(expect, output) != 0 {
		t.Errorf("expected '%s', got '%s'", strings.Join(expect, "|"), strings.Join(output, "|"))
	}
}

func TestMatches(t *testing.T) {
	tokr := gokenizer.New()

	tokr.Class("username", "{word}{number}")

	if ok, err := tokr.Matches("bob123", "{username}"); !ok || err != nil {
		t.Errorf("Expected match, got non match and err: %s", err.Error())
	}

	if ok, err := tokr.Matches("bob123foo", "{username}"); ok && err == nil {
		t.Errorf("expected non-match")
	}
}
