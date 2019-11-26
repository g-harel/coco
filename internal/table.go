package internal

import (
	"fmt"
	"strings"
)

// Must be single character.
const tableVerticalSeparator = "|"
const tableHorizontalSeparator = "-"
const tableIntersectionSeparator = "+"

// Table holds table data that can be sorted and printed.
type Table struct {
	headers []interface{}
	data    [][]interface{}
}

// Headers adds column headers.
func (t *Table) Headers(data ...interface{}) {
	t.headers = data
}

// Add appends a new row of data.
func (t *Table) Add(data ...interface{}) {
	t.data = append(t.data, data)
}

// Format formats the table data to a string.
func (t *Table) Format(columnSortPriority ...int) string {
	columnWidths := tableColumnWidths(append(t.data, t.headers))
	sortOrder := tableSortOrder(t.data, columnSortPriority)
	formattedTable := tableFormatHorizontalSeparator(columnWidths)
	formattedTable += tableFormatRow(t.headers, columnWidths)
	formattedTable += tableFormatHorizontalSeparator(columnWidths)
	for i := 0; i < len(sortOrder); i++ {
		formattedTable += tableFormatRow(t.data[sortOrder[i]], columnWidths)
	}
	formattedTable += tableFormatHorizontalSeparator(columnWidths)
	return formattedTable
}

func tableColumnWidths(data [][]interface{}) []int {
	columnWidths := []int{}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			width := len(tableFormatCell(data[i][j]))
			if len(columnWidths) <= j {
				columnWidths = append(columnWidths, width)
				continue
			}
			if width > columnWidths[j] {
				columnWidths[j] = width
			}
		}
	}
	return columnWidths
}

// Columns are prioritized by the order their index appears in
// "columnSortPriority". Column indecies that are not provided are prioritized
// the lowest and from first to last. Integers are compared as numbers, but all
// other data types are compared as strings. Nil values are given the lowest
// possible order.
func tableSortOrder(data [][]interface{}, columnSortPriority []int) []int {
	// TODO sort
	order := []int{}
	for i := 0; i < len(data); i++ {
		order = append(order, i)
	}
	return order
}

func tableFormatHorizontalSeparator(columnWidths []int) string {
	row := []string{}
	for i := 0; i < len(columnWidths); i++ {
		row = append(row, strings.Repeat(tableHorizontalSeparator, columnWidths[i] + 2))
	}
	pre := tableIntersectionSeparator
	post := tableIntersectionSeparator + "\n"
	return pre + strings.Join(row, tableIntersectionSeparator) + post;
}

// Numbers are right-aligned, all other data types are left-aligned.
func tableFormatRow(data []interface{}, columnWidths []int) string {
	row := []string{}
	for i := 0; i < len(columnWidths); i++ {
		value := ""
		number := false
		if i < len(data) {
			if _, ok := data[i].(int); ok {
				number = true
			}
			value = tableFormatCell(data[i])
		}
		if number {
			row = append(row, fmt.Sprintf("%*v", columnWidths[i], value))
		} else {
			row = append(row, fmt.Sprintf("%-*v", columnWidths[i], value))
		}
	}
	pre := tableVerticalSeparator + " "
	mid := " " + tableVerticalSeparator + " "
	post := " " + tableVerticalSeparator + "\n"
	return pre + strings.Join(row, mid) + post
}

// Integers are formatted with commas for each thousand. Nil is formatted as an
// empty string.
func tableFormatCell(data interface{}) string {
	if data == nil {
		return ""
	}
	number, ok := data.(int)
	if !ok {
		return fmt.Sprintf("%v", data)
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
