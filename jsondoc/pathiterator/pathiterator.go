package pathiterator

import (
	"fmt"
	"regexp"
	"strings"
)

var tokenRegex *regexp.Regexp

func init() {
	var err error
	tokenRegex, err = regexp.Compile("^(?:\\]\\.|\\.|\\[|\\])(\\])?([^\\.^\\[\\]]+)?")
	if err != nil {
		panic(err)
	}
}

// PathIterator type controls tokenization of JSON paths
type PathIterator struct {
	path      []byte
	token     []byte
	readIndex int
	err       error
}

// NewPathIterator creates a new PathIterator for a given path string.
func NewPathIterator(path string) StringIterator {
	return &PathIterator{
		path:      []byte(strings.TrimSpace(path)),
		token:     nil,
		readIndex: 0,
		err:       nil,
	}
}

// Next moves iterator to next token.
// Returns true if there is a valid token to be read.
// Returns false if iterator is exhausted
func (pIter *PathIterator) Next() bool {
	tokens := tokenRegex.FindSubmatch(pIter.getSubpath())
	switch {
	case tokens == nil:
		return false
	case len(tokens[1]) > 0 && len(tokens[2]) > 0:
		// If both capture groups validate then path is malformed.
		// Value function will detect the error as len(token) == 0 and len(path) > 0.
		// Setting readIndex to len(path) makes next call to Next return false.
		pIter.readIndex = len(pIter.path)
		pIter.err = fmt.Errorf("Malformed path \"%v\", syntax error near \"%v\"", string(pIter.path), string(tokens[0]))
		return true
	case len(tokens[1]) > 0:
		// Captures the forEach token []
		pIter.token = tokens[1]
	case len(tokens[2]) > 0:
		// Captures a common named token
		pIter.token = tokens[2]
	case len(pIter.path) == 1 && pIter.path[0] == byte('.'):
		pIter.token = pIter.path
		pIter.readIndex++
		return true
	default:
		return false
	}

	pIter.readIndex += len(tokens[0])
	return true
}

func (pIter PathIterator) getSubpath() []byte {
	if pIter.readIndex < len(pIter.path) {
		return pIter.path[pIter.readIndex:]
	}
	return []byte{}
}

// Value gets actual token from iterator.
func (pIter PathIterator) Value() (string, error) {
	if pIter.err != nil {
		return "", pIter.err
	}
	return string(pIter.token), nil
}
