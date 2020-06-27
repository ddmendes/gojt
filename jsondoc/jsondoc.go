package jsondoc

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/ddmendes/gojt/jsondoc/node"
	"github.com/ddmendes/gojt/jsondoc/pathiterator"
)

var (
	toNode           func(interface{}) node.Node
	toStringIterator func(string) pathiterator.StringIterator
)

func init() {
	toNode = node.ToSingleNode
	toStringIterator = pathiterator.NewPathIterator
}

// JSONDoc represents a generic JSON Document
type JSONDoc struct {
	Value node.Node
}

// Wrap anything into a JSONDoc
func Wrap(document interface{}) JSONDoc {
	return JSONDoc{
		Value: toNode(document),
	}
}

// ReadPipedDoc Reads a JSON document piped to os.Stdin
func ReadPipedDoc(jsondoc *JSONDoc) error {
	pipedData, err := getPipedData()
	if err != nil {
		*jsondoc = JSONDoc{}
		return err
	}

	input, err := unmarshalJSON(pipedData)
	if err != nil {
		*jsondoc = JSONDoc{}
		return err
	}

	*jsondoc = JSONDoc{Value: node.SingleNode{Elem: input}}
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
	return jsondoc.Value.GetKeys()
}

// Marshal returns the JSON encoding of JSONDoc. Set beautify to true
// to get a indented json document.
func (jsondoc JSONDoc) Marshal(beautify bool) ([]byte, error) {
	if beautify {
		return json.MarshalIndent(jsondoc.Value.GetInterface(), "", "  ")
	}
	return json.Marshal(jsondoc.Value.GetInterface())
}

// Get the JSON object on a given path
func (jsondoc JSONDoc) Get(path string) (JSONDoc, error) {
	actualItem := jsondoc.Value
	iter := toStringIterator(path)
	var err error

	for iter.Next() {
		token := iter.Value()
		actualItem, err = actualItem.Get(token)
		if err != nil {
			return jsondoc, err
		}
	}

	return JSONDoc{
		Value: actualItem,
	}, nil
}
