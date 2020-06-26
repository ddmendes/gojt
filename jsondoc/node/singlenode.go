package node

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	errInvalidPath = errors.New("Invalid path")
)

// SingleNode is a JSON node for a single component.
type SingleNode struct {
	Elem interface{}
}

// ToSingleNode wraps an interface with a SingleNode component.
func ToSingleNode(i interface{}) Node {
	return SingleNode{
		Elem: i,
	}
}

// Get a child Node
func (n SingleNode) Get(token string) (Node, error) {
	elem := n.Elem
	if token == "." {
		return n, nil
	} else if token == "]" {
		arr, ok := n.Elem.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Used [] token over a not array element: %v", n.Elem)
		}
		return ToMultiNode(arr), nil
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
