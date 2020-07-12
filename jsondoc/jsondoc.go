package jsondoc

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/ddmendes/gojt/jsondoc/node"
	"github.com/ddmendes/gojt/jsondoc/pathiterator"
)

const inputBufferSize int = 125

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

// JSONReader reads a JSON document from a Reader.
type JSONReader struct {
	reader io.Reader
}

// NewJSONReader creates a Reader capable of reading JSON documents
func NewJSONReader(reader io.Reader) JSONReader {
	return JSONReader{reader: reader}
}

// ReadJSON reads a document from JSONReader and wraps it into a JSONDoc.
func (r JSONReader) ReadJSON() (JSONDoc, error) {
	bufReader := bufio.NewReader(io.Reader(r.reader))
	data, err := readRawData(bufReader)
	if err != nil {
		return Wrap(nil), err
	}

	var document interface{}
	if err := json.Unmarshal(data, &document); err != nil {
		return Wrap(nil), err
	}

	return Wrap(document), nil
}

// Wrap anything into a JSONDoc
func Wrap(document interface{}) JSONDoc {
	return JSONDoc{
		Value: toNode(document),
	}
}

func readRawData(reader *bufio.Reader) ([]byte, error) {
	output := make([]byte, 0, inputBufferSize)
	for {
		input, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return output, nil
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
	node, err := sequentialGet(jsondoc.Value, path)
	if err != nil {
		return jsondoc, err
	}
	return JSONDoc{
		Value: node,
	}, nil
}

func sequentialGet(json node.Node, path string) (node.Node, error) {
	iter := toStringIterator(path)
	var token string
	var err error

	for iter.Next() {
		token, err = iter.Value()
		if err != nil {
			return nil, err
		}
		json, err = json.Get(token)
		if err != nil {
			return nil, err
		}
	}

	return json, nil
}

// runGetPipeline get object for a given path using the pipeline model.
// Using benchmark tests this feature has proven to be slower than the
// get. So code will keep using the sequential get.
// CAUTION! No errors are handled in this method as development will
// not use it.
func runGetPipeline(json node.Node, path string) (node.Node, error) {
	tokenChan := make(chan string)
	nodeChan := make(chan node.Node)
	// errChan := make(chan error)

	go func(p string, tokChan chan<- string) {
		iter := toStringIterator(p)
		for iter.Next() {
			value, _ := iter.Value()
			tokChan <- value
		}
		close(tokChan)
	}(path, tokenChan)

	go func(n node.Node, nChan chan<- node.Node, tokChan <-chan string) {
		for token := range tokChan {
			n, _ = n.Get(token)
		}
		nChan <- n
		close(nChan)
	}(json, nodeChan, tokenChan)

	nodeResult := <-nodeChan
	return nodeResult, nil
}
