package jsondoc

import (
	"testing"

	"github.com/ddmendes/gojt/jsondoc/node"
)

type nodeDouble struct {
	getCallCount          int
	getKeysCallCount      int
	getInterfaceCallCount int
}

var nodeDbl *nodeDouble = &nodeDouble{}

func TestGetKeys(t *testing.T) {
	jsondoc := JSONDoc{
		Value: nodeDbl,
	}
	defer nodeDbl.reset()

	_, err := jsondoc.GetKeys()
	if err != nil {
		t.Error("GetKeys returned an error")
	}

	if nodeDbl.getKeysCallCount == 0 {
		t.Error("Should have called underlying node GetKeys")
	}
}

func TestGet(t *testing.T) {
	jsondoc := JSONDoc{
		Value: nodeDbl,
	}
	defer nodeDbl.reset()

	_, err := jsondoc.Get("./token")
	if err != nil {
		t.Error("Get returned an error")
	}

	if nodeDbl.getCallCount == 0 {
		t.Error("Should have called underlying node Get")
	}
}

func TestMarshal(t *testing.T) {
	jsondoc := JSONDoc{
		Value: nodeDbl,
	}

	_, err := jsondoc.Marshal(true)
	if err != nil {
		t.Error("Marshal(true) returned an error")
	} else if nodeDbl.getInterfaceCallCount == 0 {
		t.Error("Should have called GetIterface() into Marshal(true)")
	}
	nodeDbl.reset()

	_, err = jsondoc.Marshal(false)
	if err != nil {
		t.Error("Marshal(false) returned an error")
	} else if nodeDbl.getInterfaceCallCount == 0 {
		t.Error("Should have called GetIterface() into Marshal(false)")
	}
	nodeDbl.reset()
}

func toNodeDouble(_ interface{}) node.Node {
	return nodeDbl
}

func (n *nodeDouble) Get(token string) (node.Node, error) {
	n.getCallCount++
	return n, nil
}

func (n *nodeDouble) GetKeys() ([]string, error) {
	n.getKeysCallCount++
	return []string{}, nil
}

func (n *nodeDouble) GetInterface() interface{} {
	n.getInterfaceCallCount++
	return interface{}(0)
}

func (n *nodeDouble) reset() {
	n.getCallCount = 0
	n.getKeysCallCount = 0
	n.getInterfaceCallCount = 0
}
