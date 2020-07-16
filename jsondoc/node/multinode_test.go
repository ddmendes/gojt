package node

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"testing"
)

func readJSONFile(path string) (MultiNode, error) {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var genericInterface interface{}
	if err := json.Unmarshal(jsonFile, &genericInterface); err != nil {
		return nil, err
	}

	interfaceArray, ok := genericInterface.([]interface{})
	if !ok {
		return nil, errors.New("JSON is not an array")
	}

	genericNode := ToMultiNode(interfaceArray)
	multiNode, ok := genericNode.(MultiNode)
	if !ok {
		return nil, errors.New("Node is not MultiNode")
	}

	return multiNode, nil
}

func BenchmarkSequentialGet(b *testing.B) {
	node, err := readJSONFile("./testdata/array.json")
	if err != nil {
		b.Fatal(err)
	}
	path := "name"
	b.ResetTimer()
	for i := 0; i < 10; i++ {
		node.sequentialGet(path)
	}
}

func BenchmarkConcurrentGet(b *testing.B) {
	node, err := readJSONFile("./testdata/array.json")
	if err != nil {
		b.Fatal(err)
	}
	path := "name"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.concurrentGet(path)
	}
}
