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
func (t *Tokenizer) Class(name string, check func(byte) bool) {
	if _, err := t.getClass(name); err == nil {
		t.err = fmt.Errorf("class '%s' already defined", name)
		return
	}

	t.classes[name] = func(iter *stringiter.StringIter) matchResult {
		word := ""

		for !iter.Eof() && check(iter.Peek()) {
			word += iter.Consume()
		}

		return matchResult{
			class:   name,
			matched: len(word) > 0,
			lexeme:  word,
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
	result := matchResult{}

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
		// Todo: more info in Token
		Lexeme: result.lexeme,
		values: result.values,
	}

	err := t.callbacks[callbackIdx](token)
	return err
}

// For future additional metadata about the expression
type matchResult struct {
	matched bool
	class   string
	lexeme  string
	values  map[string][]string
}

// Matches with the given string. The implementation is dynamically created in createPattern.
type matcherFunc func(iter *stringiter.StringIter) matchResult

// Returns a function that matches the string literal s.
func literalMatcherFunc(s string) matcherFunc {
	return func(iter *stringiter.StringIter) matchResult {
		iter.PeekN(uint(len(s)))
		lexeme := iter.Consume()

		return matchResult{
			class:   "static",
			matched: lexeme == s,
			lexeme:  lexeme,
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

	f := func(iter *stringiter.StringIter) (res matchResult) {
		iter.Push()
		values := make(map[string][]string)

		for idx, mf := range funcs {
			tempResult := mf(iter)
			if !tempResult.matched {
				return res
			}

			// Add map from previous match func
			if tempResult.class != "" {
				for k, v := range tempResult.values {
					values[k] = append(tempResult.values[k], v...)
				}
			}

			if className := classNames[idx]; className != "static" {
				values[className] = append(values[className], tempResult.lexeme)
			}
		}

		length := iter.Pop()
		iter.PeekN(uint(length))
		matchedString := iter.Consume()

		return matchResult{
			matched: true,
			class:   class,
			lexeme:  matchedString,
			values:  values,
		}
	}

	return f, err
}
