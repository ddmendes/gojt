package jsondoc

import (
	"regexp"
	"strings"
)

// PathIterator type controls tokenization of JSON paths
type PathIterator struct {
	path  []byte
	token []byte
}

// NewPathIterator creates a new PathIterator for a given path string.
func NewPathIterator(path string) PathIterator {
	return PathIterator{
		path:  []byte(strings.TrimSpace(path)),
		token: nil,
	}
}

// Next moves iterator to next token.
// Returns true if there is a valid token to be read.
// Returns false if iterator is exhausted
func (pIter *PathIterator) Next() bool {
	tokenRegex, err := regexp.Compile("^(?:\\]\\.|\\.|\\[|\\])(\\])?([^\\.^\\[\\]]+)?")
	if err != nil {
		panic(err)
	}

	tokens := tokenRegex.FindSubmatch(pIter.path)
	switch {
	case tokens == nil:
		return false
	case len(tokens[1]) > 0:
		// Captures the forEach token []
		pIter.token = tokens[1]
	case len(tokens[2]) > 0:
		// Captures a common named token
		pIter.token = tokens[2]
	default:
		if len(pIter.path) == 1 && pIter.path[0] == byte('.') {
			pIter.token = pIter.path
			pIter.path = nil
			return true
		}
		return false
	}

	pIter.path = pIter.path[len(tokens[0]):]
	return true
}

// Value gets actual token from iterator.
func (pIter PathIterator) Value() string {
	return string(pIter.token)
}
