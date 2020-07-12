package pathiterator

import (
	"testing"
)

func TestNext(t *testing.T) {
	type TestCase struct {
		iterator   StringIterator
		tokenCount int
	}

	testCases := []TestCase{
		{NewPathIterator(""), 0},
		{NewPathIterator("."), 1},
		{NewPathIterator(".hello.world.from.path.iterator"), 5},
		{NewPathIterator("  .hello.world.from.path.iterator"), 5},
		{NewPathIterator(".mastermind[1].name"), 3},
		{NewPathIterator(".mastermind[].name"), 3},
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
	type TestCase struct {
		path   string
		tokens []string
	}
	testCases := []TestCase{
		{".", []string{"."}},
		{".hello.world.from.path.iterator", []string{"hello", "world", "from", "path", "iterator"}},
		{".mastermind[0].name", []string{"mastermind", "0", "name"}},
		{".mastermind[].birth", []string{"mastermind", "]", "birth"}},
	}

	for _, testCase := range testCases {
		pathIterator := NewPathIterator(testCase.path)
		for i, want := range testCase.tokens {
			ok := pathIterator.Next()
			got, _ := pathIterator.Value()

			if !ok {
				t.Errorf("PathIterator exhausted on step %d but still want token %v for path %v.", i+1, want, testCase.path)
			}

			if got != want {
				t.Errorf("PathIterator failed on step %d for path %v. Want %v. Got %v.", i+1, testCase.path, want, got)
			}
		}
	}
}

func TestValue_Should_Error_When_PathIsInvalid(t *testing.T) {
	invalidPaths := []string{".masterminds[]name"}

	for _, path := range invalidPaths {
		pathIterator := NewPathIterator(path)
		errored := false
		for pathIterator.Next() {
			_, err := pathIterator.Value()
			if err != nil {
				errored = true
			}
		}
		if !errored {
			t.Errorf("Expected path %v to fail but iterator exausted with no errors.", path)
		}
	}
}
