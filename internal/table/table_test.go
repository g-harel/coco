package table

import (
	"testing"

	"github.com/g-harel/coco/internal/assert"
)

func TestFormat(t *testing.T) {
	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Title("test")
		tb.Headers("TEST", "ABC", "1234")
		tb.Add(0, "a", 1234)
		tb.Add("aa aaaa aa aa a")
		tb.Sort(1, 1, 12)
		assert.Equal(t, tb.String(),
			""+
				"+------+\n"+
				"| test |\n"+
				"+-----------------+-----+-------+\n"+
				"| TEST            | ABC | 1234  |\n"+
				"+-----------------+-----+-------+\n"+
				"|               0 | a   | 1,234 |\n"+
				"| aa aaaa aa aa a |     |       |\n"+
				"+-----------------+-----+-------+\n")
	})

	t.Run("", func(t *testing.T) {
		tb := Table{}
		tb.Title("a")
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
				"+---+\n"+
				"| a |\n"+
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
