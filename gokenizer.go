package gokenizer

import (
	"fmt"

	"github.com/jesperkha/gokenizer/stringiter"
)

type Tokenizer struct {
	err        error
	matchFuncs []matcherFunc
	callbacks  []func(Token) error
	classes    map[string]matcherFunc
}

// Matches with the given string. The implementation is dynamically created
// in createPattern.
type matcherFunc func(iter *stringiter.StringIter) Token

// Returns true if the character b is part of the class.
type CheckerFunc func(b byte) bool

func New() Tokenizer {
	return Tokenizer{
		classes: make(map[string]matcherFunc),
	}
}

// Pattern adds a new pattern to the tokenizer. If a match is found, the
// callback function f is called. The callback may return an error which
// will be returned by Run(). The patterns are matched by the order they
// are defined in.
func (t *Tokenizer) Pattern(pattern string, f func(Token) error) {
	mf, err := t.createMatcherFunc(pattern, "")
	if err != nil {
		t.err = err
	}
	t.matchFuncs = append(t.matchFuncs, mf)
	t.callbacks = append(t.callbacks, f)
}

// Class registers a new class with the given matcher function. The function
// should return true for any byte that is a legal character in the class.
// The class cannot override any existing names.
func (t *Tokenizer) Class(name string, check CheckerFunc) {
	if _, err := t.getClass(name); err == nil {
		t.err = fmt.Errorf("class '%s' already defined", name)
		return
	}

	t.classes[name] = checkFuncToMatchFunc(name, check)
}

// Convert boolean checker function to token matcher function.
func checkFuncToMatchFunc(class string, check CheckerFunc) matcherFunc {
	return func(iter *stringiter.StringIter) Token {
		pos := iter.Pos()
		word := ""

		for !iter.Eof() && check(iter.Peek()) {
			word += iter.Consume()
		}

		return Token{
			Pos:     pos,
			Lexeme:  word,
			Source:  iter.Source(),
			Length:  len(word),
			class:   class,
			matched: len(word) > 0,
		}
	}
}

// ClassFromPattern creates a new class that matches to the given pattern.
// The class cannot override any existing names.
func (t *Tokenizer) ClassFromPattern(name string, pattern string) {
	if _, err := t.getClass(name); err == nil {
		t.err = fmt.Errorf("class '%s' already defined", name)
		return
	}

	f, err := t.createMatcherFunc(pattern, name)
	if err != nil {
		t.err = err
		return
	}

	t.classes[name] = f
}

// Run tokenizer on given input string. Returns first error received by a
// pattern callback function. Patterns are matched by the order the are
// defined in.
func (t *Tokenizer) Run(s string) error {
	if t.err != nil {
		return t.err
	}

	iter := stringiter.New(s)
	for !iter.Eof() {
		if err := t.matchNext(&iter); err != nil {
			return err
		}
	}

	return nil
}

// Continue matching until one is found. Returns callbacks error.
func (t *Tokenizer) matchNext(iter *stringiter.StringIter) error {
	callbackIdx := 0
	result := Token{}
	pos := iter.Pos()

	for idx, mf := range t.matchFuncs {
		iter.Push()
		if result = mf(iter); result.matched {
			callbackIdx = idx
			break
		}

		iter.Pop()
	}

	if !result.matched {
		iter.Consume() // Next
		return nil
	}

	token := Token{
		Pos:    pos,
		Lexeme: result.Lexeme,
		Source: iter.Source(),
		Length: len(result.Lexeme),
		values: result.values,
	}

	err := t.callbacks[callbackIdx](token)
	return err
}

// Returns a function that matches the string literal s.
func literalMatcherFunc(s string) matcherFunc {
	return func(iter *stringiter.StringIter) Token {
		pos := iter.Pos()
		iter.PeekN(uint(len(s)))
		lexeme := iter.Consume()

		return Token{
			Lexeme:  lexeme,
			Pos:     pos,
			Length:  len(lexeme),
			Source:  iter.Source(),
			class:   "static",
			matched: lexeme == s,
		}
	}
}

func parseClass(iter *stringiter.StringIter) (name string, err error) {
	if iter.Consume() != "{" {
		return name, fmt.Errorf("expected { before class name")
	}

	if !iter.Seek('}') {
		return name, fmt.Errorf("expected } after class name")
	}

	name = iter.Consume()
	iter.Consume() // }
	return name, err
}

// Returns two equal length lists of matcher functions and their class names.
func (t *Tokenizer) parsePattern(pattern string) (funcs []matcherFunc, classNames []string, err error) {
	pIter := stringiter.New(pattern)

	for !pIter.Eof() {
		if pIter.Peek() == '{' {
			// Parse class name if we find a {
			pIter.Restore()
			className, err := parseClass(&pIter)
			if err != nil {
				return funcs, classNames, err
			}

			f, err := t.getClass(className)
			if err != nil {
				return funcs, classNames, err
			}

			funcs = append(funcs, f)
			classNames = append(classNames, className)
		} else if pIter.Seek('{') {
			// Parse static word if there are characters before a {
			staticWord := pIter.Consume()
			if staticWord == "" {
				return funcs, classNames, fmt.Errorf("parser error")
			}

			funcs = append(funcs, literalMatcherFunc(staticWord))
			classNames = append(classNames, "static")
		} else {
			// Otherwise the rest of the pattern string is a static word
			funcs = append(funcs, literalMatcherFunc(pIter.Remainder()))
			classNames = append(classNames, "static")
			break
		}
	}

	return funcs, classNames, err
}

// Returns class matchFunc from either global or local context
func (t *Tokenizer) getClass(name string) (mf matcherFunc, err error) {
	mf, ok := classes[name]
	if !ok {
		mf, ok = t.classes[name]
		if !ok {
			return mf, fmt.Errorf("unknown class '%s'", name)
		}
	}

	return mf, err
}

// Returns a function that matches based on the given pattern.
func (t *Tokenizer) createMatcherFunc(pattern string, class string) (mf matcherFunc, err error) {
	funcs, classNames, err := t.parsePattern(pattern)
	if err != nil {
		return mf, err
	}

	f := func(iter *stringiter.StringIter) (res Token) {
		pos := iter.Pos()
		iter.Push()
		values := make(map[string][]Token)

		for idx, mf := range funcs {
			pos := iter.Pos()
			tempResult := mf(iter)
			if !tempResult.matched {
				return res
			}

			tempResult.Pos = pos
			tempResult.Length = len(tempResult.Lexeme)
			tempResult.Source = iter.Source()

			if className := classNames[idx]; className != "static" {
				values[className] = append(values[className], tempResult)
			}
		}

		length := iter.Pop()
		iter.PeekN(uint(length))
		matchedString := iter.Consume()

		return Token{
			Lexeme:  matchedString,
			Length:  len(matchedString),
			Pos:     pos,
			Source:  iter.Source(),
			matched: true,
			class:   class,
			values:  values,
		}
	}

	return f, err
}
