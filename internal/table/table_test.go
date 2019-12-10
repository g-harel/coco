package table

import (
	"testing"

	"github.com/g-harel/coco/internal/assert"
)

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
			assert.Equal(t, formatCell(tc.Input), tc.Expected)
		})
	}
}

func TestFormat(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Headers("TEST", "ABC", "1234")
		tb.Add(0, "a", 1234)
		tb.Add("aa aaaa aa aa a")
		tb.Sort(1, 1, 12)
		assert.Equal(t, tb.String(),
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
		tb.Add(nil, "b", 1, "b")
		tb.Add(nil, nil, 2)
		tb.Add(nil, "a", 1, "a")
		tb.Sort(2, 1)
		assert.Equal(t, tb.String(),
			""+
				"+---+---+---+---+\n"+
				"| A | B | C | D |\n"+
				"+---+---+---+---+\n"+
				"|   |   | 2 |   |\n"+
				"| 1 | 1 | 1 | 1 |\n"+
				"| 1 | 1 | 1 |   |\n"+
				"|   | a | 1 | a |\n"+
				"|   | b | 1 | b |\n"+
				"| 1 | 1 |   | 1 |\n"+
				"| 1 |   |   |   |\n"+
				"+---+---+---+---+\n")
	})
}
