package jsondoc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

var (
	errInvalidPath = errors.New("Invalid path")
)

// JSONDoc represents a generic JSON Document
type JSONDoc struct {
	Value interface{}
	Err   error
}

// Wrap anything into a JSONDoc
func Wrap(document interface{}) JSONDoc {
	return JSONDoc{
		Value: document,
		Err:   nil,
	}
}

// ReadPipedDoc Reads a JSON document piped to os.Stdin
func ReadPipedDoc(jsondoc *JSONDoc) error {
	pipedData, err := getPipedData()
	if err != nil {
		*jsondoc = JSONDoc{nil, err}
		return err
	}

	input, err := unmarshalJSON(pipedData)
	if err != nil {
		*jsondoc = JSONDoc{nil, err}
		return err
	}

	*jsondoc = JSONDoc{input, nil}
	return nil
}

func getPipedData() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return []byte{}, err
	}

	if info.Mode()&os.ModeCharDevice == 0 {
		var output []byte
		reader := bufio.NewReader(os.Stdin)
		for {
			input, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			output = append(output, input)
		}
		return []byte(output), nil
	}
	return []byte{}, errors.New("No piped data")
}

func unmarshalJSON(input []byte) (interface{}, error) {
	var document interface{}
	err := json.Unmarshal(input, &document)
	return document, err
}

// GetKeys return an slice of strings with JSONDoc keys
func (jsondoc JSONDoc) GetKeys() ([]string, error) {
	m, ok := jsondoc.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("%v has no keys", jsondoc.Value)
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys, nil
}

// Marshal returns the JSON encoding of JSONDoc. Set beautify to true
// to get a indented json document.
func (jsondoc JSONDoc) Marshal(beautify bool) ([]byte, error) {
	if jsondoc.Err != nil {
		return nil, jsondoc.Err
	}

	if beautify {
		return json.MarshalIndent(jsondoc.Value, "", "  ")
	}
	return json.Marshal(jsondoc.Value)
}

// Get the JSON object on a given path
func (jsondoc JSONDoc) Get(path string) (JSONDoc, error) {
	actualItem := jsondoc.Value
	pathIterator := NewPathIterator(path)
	var err error

	for pathIterator.Next() {
		token := pathIterator.Value()
		actualItem, err = get(actualItem, token)
		if err != nil {
			return jsondoc, err
		}
	}

	return JSONDoc{
		Value: actualItem,
		Err:   nil,
	}, nil
}

func get(docInterface interface{}, token string) (interface{}, error) {
	if token == "." {
		return docInterface, nil
	}

	switch doc := docInterface.(type) {
	case map[string]interface{}:
		return getFromMap(doc, token)
	case []interface{}:
		index, err := strconv.Atoi(token)
		if err != nil {
			return nil, errInvalidPath
		}
		return getFromArray(doc, index)
	default:
		return nil, errInvalidPath
	}
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
