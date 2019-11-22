package internal

import "testing"

func TestFormatTableCell(t *testing.T) {
	t.Run("should return empty string for nil", func(t *testing.T) {
		expected := ""
		actual := formatTableCell(nil)
		if actual != expected {
			t.Errorf("actual/expected don't match: '%v' != '%v'", actual, expected)
		}
	})

	t.Run("should not modify string formatting", func(t *testing.T) {
		expected := "test"
		actual := formatTableCell(expected)
		if actual != expected {
			t.Errorf("actual/expected don't match: '%v' != '%v'", actual, expected)
		}
	})

	t.Run("should add commas to numbers", func(t *testing.T) {
		expected := "19,010,123"
		actual := formatTableCell(19010123)
		if actual != expected {
			t.Errorf("actual/expected don't match: '%v' != '%v'", actual, expected)
		}
	})
}
