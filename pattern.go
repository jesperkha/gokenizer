package gokenizer

// For future additional metadata about the expression
type matchResult struct {
	matched bool
}

type matcherFunc func(s string) matchResult

type pattern struct {
	s  string
	mf matcherFunc
}

func createPattern(s string) (p pattern, err error) {

	p.s = s
	return p, err
}

func parseClass(s string) (mf matcherFunc, err error) {

	return mf, err
}
