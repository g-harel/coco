package assert

import (
	"strings"
	"testing"
)

func Equal(t *testing.T, actual, expected string) {
	actual = strings.ReplaceAll(actual, "\n", "\\n")
	expected = strings.ReplaceAll(expected, "\n", "\\n")
	if actual != expected {
		t.Errorf("actual/expected don't match: \n  actual: '%v'\nexpected: '%v'", actual, expected)
	}
}
