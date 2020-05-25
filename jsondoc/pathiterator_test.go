package jsondoc_test

import (
	"testing"

	"github.com/ddmendes/gojt/jsondoc"
)

func TestNext(t *testing.T) {
	type TestCase struct {
		iterator   jsondoc.PathIterator
		tokenCount int
	}

	testCases := []TestCase{
		{jsondoc.NewPathIterator(""), 0},
		{jsondoc.NewPathIterator("."), 1},
		{jsondoc.NewPathIterator(".hello.world.from.path.iterator"), 5},
		{jsondoc.NewPathIterator("  .hello.world.from.path.iterator"), 5},
	}

	for _, testCase := range testCases {
		for i := testCase.tokenCount; i > 0; i-- {
			want := true
			got := testCase.iterator.Next()
			if got != want {
				t.Errorf("Want %v. Got %v.\nTest case %v.", want, got, testCase)
			}
		}
		want := false
		got := testCase.iterator.Next()
		if got != want {
			t.Errorf("Want %v. Got %v.\nTest case %v.", want, got, testCase)
		}
	}
}

func TestValue(t *testing.T) {
	path := ".hello.world.from.path.iterator"
	tokens := []string{
		"hello",
		"world",
		"from",
		"path",
		"iterator",
	}

	pathIterator := jsondoc.NewPathIterator(path)
	for i, want := range tokens {
		ok := pathIterator.Next()
		got := pathIterator.Value()

		if !ok {
			t.Errorf("PathIterator exhausted on step %d but still want token %v for path %v.", i+1, want, path)
		}

		if got != want {
			t.Errorf("PathIterator failed on step %d for path %v. Want %v. Got %v.", i+1, path, want, got)
		}
	}
}
