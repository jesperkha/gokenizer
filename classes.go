package gokenizer

import "github.com/jesperkha/gokenizer/stringiter"

var classes = map[string]matcherFunc{
	"word": func(iter *stringiter.StringIter) matchResult {
		return matchResult{}
	},

	"number": func(iter *stringiter.StringIter) matchResult {
		return matchResult{}
	},

	"symbol": func(iter *stringiter.StringIter) matchResult {
		return matchResult{}
	},
}
