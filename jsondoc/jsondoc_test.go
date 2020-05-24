package jsondoc_test

import (
	"testing"

	"github.com/ddmendes/gojt/jsondoc"
)

var testJSONDocKeys = []string{
	"strElem",
	"boolElem",
	"nilElem",
	"numberElem",
}

func loadTestJSONDoc() jsondoc.JSONDoc {
	return jsondoc.JSONDoc{
		Value: map[string]interface{}{
			"strElem":    "foobar",
			"boolElem":   true,
			"nilElem":    nil,
			"numberElem": float64(3.1415),
		},
		Err: nil,
	}
}

/*
This test is not working yet.
Need to properly set FileMode flag ModeNamedPipe

func TestReadPipedDoc(t *testing.T) {
	feedContent := []byte("{\"foo\":{\"bar\":\"baz\"}, \"numbers\":[1,2,3,4]}")
	tmpFile, err := ioutil.TempFile("", "FakeStdin")
	if err != nil {
		t.Fatal("Failed to create temp file")
	}
	defer os.Remove(tmpFile.Name())

	info, err := tmpFile.Stat()
	if err != nil {
		t.Fatal("Failed to get tmpFile Stat")
	}

	if err := tmpFile.Chmod(os.FileMode(info.Mode() | os.ModeNamedPipe)); err != nil {
		t.Fatal("Failed to set temp file mode to named pipe")
	}

	if _, err := tmpFile.Write(feedContent); err != nil {
		t.Fatal("Failed to write data input on temp file")
	}

	if _, err := tmpFile.Seek(0, 0); err != nil {
		t.Fatal("Failed to seek temp file back to beginning")
	}

	stdin := os.Stdin
	os.Stdin = tmpFile
	defer func() { os.Stdin = stdin }()

	var got jsondoc.JSONDoc
	err = jsondoc.ReadPipedDoc(&got)

	if err != nil {
		t.Fatal(err)
	}

	failTest := func() { t.Error("Failed to read json from Stdin") }
	docMap, ok := got.Value.(map[string]interface{})
	if !ok {
		failTest()
	}
	foo, ok := docMap["foo"]
	if !ok {
		failTest()
	}

	fooMap, ok := foo.(map[string]interface{})
	if !ok {
		failTest()
	}

	bar, ok := fooMap["bar"]
	if !ok {
		failTest()
	}

	barStr, ok := bar.(string)
	if !ok {
		failTest()
	}

	if barStr != "baz" {
		failTest()
	}

	numbers, ok := docMap["numbers"]
	if !ok {
		failTest()
	}

	numbersSlice, ok := numbers.([]float64)
	if !ok {
		failTest()
	}

	if !reflect.DeepEqual(numbersSlice, []float64{1, 2, 3, 4}) {
		failTest()
	}
}
*/

func TestGetKeys(t *testing.T) {
	document := loadTestJSONDoc()
	got := document.GetKeys()
	want := testJSONDocKeys

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
