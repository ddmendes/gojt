package pathiterator

// StringIterator defines an interface to Iterate over string items
type StringIterator interface {
	// Next tells whether or not there are still items to be read from Iterator.
	Next() bool
	// Value returns the next string values to read.
	Value() (string, error)
}
