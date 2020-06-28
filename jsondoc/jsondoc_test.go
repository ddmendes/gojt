package jsondoc

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ddmendes/gojt/jsondoc/node"
	"github.com/ddmendes/gojt/jsondoc/pathiterator"
)

type nodeDouble struct {
	document              interface{}
	getCallCount          int
	getKeysCallCount      int
	getInterfaceCallCount int
}

type iteratorDouble struct {
	nextResponses  []bool
	valueResponses []string
	nextIndex      int
	nextCallCount  int
	valueCallCount int
}

var nodeDbl *nodeDouble = &nodeDouble{}
var iterDbl *iteratorDouble = &iteratorDouble{
	nextResponses:  []bool{true, false},
	valueResponses: []string{"token"},
}

const sampleInput string = "{\"key\":\"foo\",\"value\":\"bar\"}"

func getSampleFile() (*os.File, error) {
	temp, err := ioutil.TempFile(os.TempDir(), "test_input")
	if err != nil {
		return &os.File{}, err
	}

	if _, err := temp.Write([]byte(sampleInput)); err != nil {
		return temp, err
	}

	if _, err := temp.Seek(0, 0); err != nil {
		return temp, err
	}

	return temp, nil
}

func TestNewJSONReader(t *testing.T) {
	want, err := getSampleFile()
	defer os.Remove(want.Name())
	if err != nil {
		t.Fatal(err)
	}

	jsonReader := NewJSONReader(want)
	got := jsonReader.reader
	if got != want {
		t.Error("Should wrap given file")
	}
}

func TestReadJSON(t *testing.T) {
	// Mock Node
	origToNode := toNode
	toNode = toNodeDouble
	defer func() { toNode = origToNode }()

	// Read from temp file
	temp, err := getSampleFile()
	defer os.Remove(temp.Name())
	if err != nil {
		t.Fatal(err)
	}

	sampleDoc := map[string]string{
		"key":   "foo",
		"value": "bar",
	}

	jsonReader := JSONReader{reader: temp}
	jsonDoc, err := jsonReader.ReadJSON()

	doc, ok := jsonDoc.Value.GetInterface().(map[string]interface{})
	if !ok {
		t.Error("Wrong underlying interface")
	}

	for key, interfaceValue := range doc {
		got, ok := interfaceValue.(string)
		if !ok {
			t.Error("Generic interface is not string")
		}

		want, ok := sampleDoc[key]
		if !ok {
			t.Errorf("Unexpected key: %v", key)
		} else if got != want {
			t.Errorf("Want: %v; Got: %v", sampleDoc, doc)
		}
	}
}

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
	origToStringIterator := toStringIterator
	toStringIterator = toIteratorDouble
	defer func() { toStringIterator = origToStringIterator }()
	defer iterDbl.reset()

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

func toNodeDouble(i interface{}) node.Node {
	nodeDbl.document = i
	return nodeDbl
}

func toIteratorDouble(_ string) pathiterator.StringIterator {
	return iterDbl
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
	return n.document
}

func (n *nodeDouble) reset() {
	n.document = nil
	n.getCallCount = 0
	n.getKeysCallCount = 0
	n.getInterfaceCallCount = 0
}

func (iter *iteratorDouble) Next() bool {
	iter.nextCallCount++
	res := iter.nextResponses[iter.nextIndex]
	if iter.nextIndex == iter.valueCallCount {
		iter.nextIndex++
	}
	return res
}

func (iter *iteratorDouble) Value() string {
	res := iter.valueResponses[iter.valueCallCount]
	iter.valueCallCount++
	return res
}

func (iter *iteratorDouble) reset() {
	iter.nextIndex = 0
	iter.nextCallCount = 0
	iter.valueCallCount = 0
}
