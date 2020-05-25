package jsondoc_test

import (
	"testing"

	"github.com/ddmendes/gojt/jsondoc"
)

func TestNext(t *testing.T) {
	type TestTable struct {
		iterator jsondoc.PathIterator
		want     bool
	}
	cases := []TestTable{
		{jsondoc.NewPathIterator(""), false},
		{jsondoc.NewPathIterator("."), true},
		{jsondoc.NewPathIterator(".hello.world.from.path.iterator"), true},
		{jsondoc.NewPathIterator("  .hello.world.from.path.iterator"), true},
	}
	for _, testCase := range cases {
		got := testCase.iterator.Next()
		if got != testCase.want {
			t.Errorf("Want %v. Got %v", got, testCase.want)
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
