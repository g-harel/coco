package internal

import (
	"fmt"
	"strings"
)

const columnSeparator = " | "

// Table holds table data that can be sorted and printed.
type Table struct {
	headers []string
	data    [][]interface{}
}

// Headers adds column headers.
func (t *Table) Headers(titles ...string) {
	t.headers = titles
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
	// Calculate column widths from max of headers and cells.
	columnWidths := []int{}
	for i := 0; i < len(t.headers); i++ {
		width := len(formatTableHeader(t.headers[i]))
		columnWidths = append(columnWidths, width)
	}
	for i := 0; i < len(t.data); i++ {
		for j := 0; j < len(t.data[i]); j++ {
			width := len(formatTableCell(t.data[i][j]))
			if len(columnWidths) < j {
				columnWidths = append(columnWidths, width)
				continue
			}
			if width > columnWidths[j] {
				columnWidths[j] = width
			}
		}
	}
	headers := []string{}
	for i := 0; i < len(columnWidths); i++ {
		value := ""
		if len(t.headers) > i {
			value = formatTableHeader(t.headers[i])
		}
		headers = append(headers, fmt.Sprintf("%-*v", columnWidths[i], value))
	}
	return strings.Join(headers, columnSeparator)
}

func formatTableHeader(value string) string {
	return strings.ToUpper(value)
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
