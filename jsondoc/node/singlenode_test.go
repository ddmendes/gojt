package node

import "testing"

func TestToMultiNode(t *testing.T) {
	want := SingleNode{
		Elem: []interface{}{
			map[string]string{"key": "foo", "value": "bar"},
			map[string]string{"key": "jhon", "value": "doe"},
		},
	}

	multiNode, err := want.ToMultiNode()
	if err != nil {
		t.Fatal(err)
	}

	elemArray := want.Elem.([]interface{})
	if len(elemArray) != len(multiNode) {
		t.Errorf("Wrong length of underlying array. Want %v but got %v", len(elemArray), len(multiNode))
	}
}
