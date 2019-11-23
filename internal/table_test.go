package internal

import "testing"

func assertEqual(t *testing.T, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("actual/expected don't match: \n  actual: '%v'\nexpected: '%v'", actual, expected)
	}
}

func TestFormatTableCell(t *testing.T) {
	tt := []struct {
		Description string
		Input       interface{}
		Expected    string
	}{
		{"should produce empty string for nil", nil, ""},
		{"should not change string formatting", "test", "test"},
		{"should correctly handle zero", 0, "0"},
		{"should correctly handle small numbers", 999, "999"},
		{"should add commas to large numbers", 19010123, "19,010,123"},
		{"should correctly handle small negative numbers", -12, "-12"},
		{"should correctly handle large negative numbers", -398111, "-398,111"},
	}

	for i := 0; i < len(tt); i++ {
		tc := tt[i]
		t.Run(tc.Description, func(t *testing.T) {
			assertEqual(t, formatTableCell(tc.Input), tc.Expected)
		})
	}
}

func TestFormat(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Headers("test", "abc", "1234")
		tb.Add(0, 1, 1234)
		assertEqual(t, tb.Format(), "")
	})
}
