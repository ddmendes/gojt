package jsondoc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// JSONDoc represents a generic JSON Document
type JSONDoc struct {
	Value interface{}
	Err   error
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

	fmt.Printf("NamedPipe: %v\n", int(info.Mode()&os.ModeNamedPipe))
	fmt.Printf("Size: %v\n", info.Size())
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

// GetKeys return an slice of strings with JSONDoc keys
func (jsondoc JSONDoc) GetKeys() []string {
	m, ok := jsondoc.Value.(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("%v has no keys", jsondoc.Value))
	}

	var keys = make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
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
