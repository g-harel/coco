package internal

import (
	"strings"
	"testing"
)

func assertEqual(t *testing.T, actual, expected string) {
	actual = strings.ReplaceAll(actual, "\n", "\\n")
	expected = strings.ReplaceAll(expected, "\n", "\\n")
	if actual != expected {
		t.Errorf("actual/expected don't match: \n  actual: '%v'\nexpected: '%v'", actual, expected)
	}
}

func TestTableFormatCell(t *testing.T) {
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
			assertEqual(t, tableFormatCell(tc.Input), tc.Expected)
		})
	}
}

func TestFormat(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Headers("TEST", "ABC", "1234")
		tb.Add(0, "a", 1234)
		tb.Add("aa aaaa aa aa a")
		assertEqual(t, tb.Format(1, 1, 12),
			""+
				"+-----------------+-----+-------+\n"+
				"| TEST            | ABC | 1234  |\n"+
				"+-----------------+-----+-------+\n"+
				"|               0 | a   | 1,234 |\n"+
				"| aa aaaa aa aa a |     |       |\n"+
				"+-----------------+-----+-------+\n")
	})

	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Headers("A", "B", "C", "D")
		tb.Add(1)
		tb.Add(1, 1, 1)
		tb.Add(1, 1, 1, 1)
		tb.Add(1, 1, nil, 1)
		tb.Add(nil, nil, 2)
		tb.Add(nil, "a", 1, "a")
		assertEqual(t, tb.Format(2, 1),
			""+
				"+---+---+---+---+\n"+
				"| A | B | C | D |\n"+
				"+---+---+---+---+\n"+
				"|   |   | 2 |   |\n"+
				"|   | a | 1 | a |\n"+
				"| 1 | 1 | 1 | 1 |\n"+
				"| 1 | 1 | 1 |   |\n"+
				"| 1 | 1 |   | 1 |\n"+
				"| 1 |   |   |   |\n"+
				"+---+---+---+---+\n")
	})
}
