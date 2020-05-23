package jsdoc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// JSDoc represents a generic JSON Document
type JSDoc struct {
	Value interface{}
	Err   error
}

// ReadPipedDoc Reads a JSON document piped to os.Stdin
func ReadPipedDoc(jsdoc *JSDoc) error {
	pipedData, err := getPipedData()
	if err != nil {
		*jsdoc = JSDoc{nil, err}
		return err
	}

	input, err := unmarshalJSON(pipedData)
	if err != nil {
		*jsdoc = JSDoc{nil, err}
		return err
	}

	*jsdoc = JSDoc{input, nil}
	return nil
}

func getPipedData() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return []byte{}, err
	}

	if info.Mode()&os.ModeNamedPipe != 0 && info.Size() > 0 {
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

// GetKeys return an slice of strings with JSDoc keys
func (jsdoc JSDoc) GetKeys() []string {
	m, ok := jsdoc.Value.(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("%v has no keys", jsdoc.Value))
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

// Marshal returns the JSON encoding of JSDoc. Set beautify to true
// to get a indented json document.
func (jsdoc JSDoc) Marshal(beautify bool) ([]byte, error) {
	if jsdoc.Err != nil {
		return nil, jsdoc.Err
	}

	if beautify {
		return json.MarshalIndent(jsdoc.Value, "", "  ")
	}
	return json.Marshal(jsdoc.Value)
}
