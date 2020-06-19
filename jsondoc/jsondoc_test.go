package jsondoc_test

import (
	"testing"

	"github.com/ddmendes/gojt/jsondoc"
	"github.com/ddmendes/gojt/jsondoc/node"
)

var testJSONDocKeys = []string{
	"strElem",
	"boolElem",
	"nilElem",
	"numberElem",
	"numArrElem",
}

func loadTestJSONDoc() jsondoc.JSONDoc {
	return jsondoc.JSONDoc{
		Value: node.SingleNode{
			Elem: map[string]interface{}{
				"strElem":    "foobar",
				"boolElem":   true,
				"nilElem":    nil,
				"numberElem": float64(3.1415),
				"numArrElem": []interface{}{
					float64(1),
					float64(2),
					float64(3),
					float64(4),
				},
			},
		},
	}
}

func TestGetKeys(t *testing.T) {
	document := loadTestJSONDoc()
	got, err := document.GetKeys()
	want := testJSONDocKeys

	if err != nil {
		t.Fatal("Failed to get keys from JSON")
	}

	if len(got) != len(want) {
		t.Errorf("Want %v but got %v", want, got)
	}

	var found bool
	for wantElem := range want {
		found = false
		for gotElem := range got {
			if wantElem == gotElem {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Key %v not found in %v", wantElem, got)
		}
	}
}

func TestGet(t *testing.T) {
	testJSON := loadTestJSONDoc()

	type TestCase struct {
		doc  jsondoc.JSONDoc
		path string
		want interface{}
	}

	testCases := []TestCase{
		{testJSON, ".strElem", "foobar"},
		{testJSON, ".boolElem", true},
		{testJSON, ".nilElem", nil},
		{testJSON, ".numberElem", float64(3.1415)},
		{testJSON, ".numArrElem.0", float64(1)},
		{testJSON, ".numArrElem[0]", float64(1)},
		{testJSON, ".numArrElem.1", float64(2)},
		{testJSON, ".numArrElem[1]", float64(2)},
		{testJSON, ".numArrElem.2", float64(3)},
		{testJSON, ".numArrElem[2]", float64(3)},
		{testJSON, ".numArrElem.3", float64(4)},
		{testJSON, ".numArrElem[3]", float64(4)},
	}

	for _, testCase := range testCases {
		child, err := testCase.doc.Get(testCase.path)
		if err != nil {
			t.Errorf("Failed to get path \"%v\" from JSON %v (error: %v)", testCase.path, testCase.doc, err)
		}
		got := child.Value.GetInterface()
		if testCase.want != got {
			t.Errorf("Want %v but got %v on path %v of JSON %v", testCase.want, got, testCase.path, testCase.doc)
		}
	}
}
