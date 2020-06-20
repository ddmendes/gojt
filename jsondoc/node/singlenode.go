package node

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errInvalidPath = errors.New("Invalid path")
)

// SingleNode is a JSON node for a single component
type SingleNode struct {
	Elem interface{}
}

// Get a child Node
func (n SingleNode) Get(token string) (Node, error) {
	elem := n.Elem
	if token == "." {
		return n, nil
	} else if token == "]" {
		return n.ToMultiNode()
	}

	switch e := elem.(type) {
	case map[string]interface{}:
		child, err := getFromMap(e, token)
		return SingleNode{Elem: child}, err
	case []interface{}:
		index, err := strconv.Atoi(token)
		if err != nil {
			return nil, errInvalidPath
		}
		child, err := getFromArray(e, index)
		return SingleNode{Elem: child}, err
	default:
		return nil, errInvalidPath
	}
}

// GetKeys returns all the keys available for this Node.
func (n SingleNode) GetKeys() ([]string, error) {
	m, ok := n.Elem.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%v has no keys", n.Elem)
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys, nil
}

func getFromMap(mapElem map[string]interface{}, token string) (interface{}, error) {
	child, ok := mapElem[token]
	if !ok {
		return nil, errInvalidPath
	}
	return child, nil
}

func getFromArray(arrElem []interface{}, index int) (interface{}, error) {
	if index < len(arrElem) {
		return arrElem[index], nil
	}
	return nil, errInvalidPath
}

// GetInterface gets the interface{} value of this Node component
func (n SingleNode) GetInterface() interface{} {
	return n.Elem
}

// ToMultiNode wraps an array Node into a MultiNode component.
// If actual Node is not an array an error is returned.
func (n SingleNode) ToMultiNode() (MultiNode, error) {
	arr, ok := n.Elem.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Used [] token over a not array element: %v", n.Elem)
	}

	singleNodes := make([]Node, len(arr))
	for i, v := range arr {
		singleNodes[i] = SingleNode{Elem: v}
	}

	return MultiNode(singleNodes), nil
}
