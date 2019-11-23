package internal

import (
	"fmt"
	"strings"
)

// Table holds table data that can be sorted and printed.
type Table struct {
	titles []string
	data   [][]interface{}
}

// Titles adds table titles.
func (t *Table) Titles(titles ...string) {
	t.titles = titles
}

// Add appends a new row of data.
func (t *Table) Add(data ...interface{}) {
	t.data = append(t.data, data)
}

// Format formats the table data to a string.
// First provided column index will have highest priority. Intgers are compared
// as numbers, but all other data types are compared as strings. Nil values are
// given the lowest possible order. Column indecies that are not provided are
// prioritized last and in the same order. Numbers are right-aligned and
// formatted with commas.
func (t *Table) Format(columnSort ...int) string {
	return ""
}

func formatTableCell(value interface{}) string {
	if value == nil {
		return ""
	}
	number, ok := value.(int)
	if !ok {
		return fmt.Sprintf("%v", value)
	}
	if number == 0 {
		return "0"
	}
	sign := ""
	if number < 0 {
		number = -number
		sign = "-"
	}
	parts := []string{}
	for number > 0 {
		parts = append([]string{fmt.Sprintf("%03v", number%1000)}, parts...)
		number /= 1000
	}
	return sign + strings.TrimLeft(strings.Join(parts, ","), "0")
}
