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
