package node

import "errors"

// MultiNode represents multiple json components wrapped.
type MultiNode []Node

// ToMultiNode wraps an []interface{} with a MultiNode component.
func ToMultiNode(iArr []interface{}) Node {
	singleNodes := make([]Node, len(iArr))
	for i, v := range iArr {
		singleNodes[i] = SingleNode{Elem: v}
	}

	return MultiNode(singleNodes)
}

// Get a child Node for each underlying Node.
func (n MultiNode) Get(token string) (Node, error) {
	var err error
	for i, v := range n {
		n[i], err = v.Get(token)
		if err != nil {
			return nil, err
		}
	}
	return n, nil
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
