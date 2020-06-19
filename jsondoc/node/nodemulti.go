package node

import "errors"

// MultiNode represents multiple json components wrapped
type MultiNode []Node

// Get a child Node
func (n MultiNode) Get(token string) (Node, error) {
	return MultiNode{}, nil
}

// GetKeys returns all the keys available for this Node.
func (n MultiNode) GetKeys() ([]string, error) {
	if len(n) > 0 {
		return n[0].GetKeys()
	}
	return []string{}, errors.New("Document is empty")
}

// GetInterface gets the interface{} value of this Node component
func (n MultiNode) GetInterface() interface{} {
	interfaces := make([]interface{}, len(n), len(n))
	for i, v := range n {
		interfaces[i] = v.GetInterface()
	}
	return interface{}(interfaces)
}
