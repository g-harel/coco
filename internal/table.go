package internal

import (
	"fmt"
	"sort"
	"strings"
)

// Must be single character.
const tableVerticalSeparator = "|"
const tableHorizontalSeparator = "-"
const tableIntersectionSeparator = "+"

// Table holds table data that can be sorted and printed.
type Table struct {
	headers            []interface{}
	columnSortPriority []int
	data               [][]interface{}
}

// Headers adds column headers.
func (t *Table) Headers(data ...interface{}) {
	ExecSafe(func() {
		t.headers = data
	})
}

// Add appends a new row of data.
func (t *Table) Add(data ...interface{}) {
	ExecSafe(func() {
		t.data = append(t.data, data)
	})
}

// Sort sets the columns to sort by when formating.
func (t *Table) Sort(columnSortPriority ...int) {
	ExecSafe(func() {
		t.columnSortPriority = columnSortPriority
	})
}

// String formats the table data to a string.
func (t *Table) String() string {
	columnWidths := tableColumnWidths(append(t.data, t.headers))
	sortOrder := tableSortOrder(t.data, t.columnSortPriority)
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

func tableFormatHorizontalSeparator(columnWidths []int) string {
	row := []string{}
	for i := 0; i < len(columnWidths); i++ {
		row = append(row, strings.Repeat(tableHorizontalSeparator, columnWidths[i]+2))
	}
	pre := tableIntersectionSeparator
	post := tableIntersectionSeparator + "\n"
	return pre + strings.Join(row, tableIntersectionSeparator) + post
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

// Columns are prioritized by the order their index appears in
// "columnSortPriority". Column indecies that are not provided are prioritized
// the lowest and from first to last. Integers are compared as numbers, but all
// other data types are compared as strings. Nil values are given the lowest
// possible order.
func tableSortOrder(data [][]interface{}, columnSortPriority []int) []int {
	// Calculate number of columns as the max row length in the data.
	columnCount := 0
	for i := 0; i < len(data); i++ {
		if len(data[i]) > columnCount {
			columnCount = len(data[i])
		}
	}
	// Clean up sort priority by removing duplicates, adding missing indecies,
	// and removing out-of-bounds indecies.
	actualSortPriority := []int{}
	seenColumnIndecies := map[int]bool{}
	for i := 0; i < len(columnSortPriority); i++ {
		columnIndex := columnSortPriority[i]
		if columnIndex >= columnCount {
			continue
		}
		if seenColumnIndecies[columnIndex] {
			continue
		}
		actualSortPriority = append(actualSortPriority, columnIndex)
		seenColumnIndecies[columnIndex] = true
	}
	for i := 0; i < columnCount; i++ {
		if seenColumnIndecies[i] {
			continue
		}
		actualSortPriority = append(actualSortPriority, i)
	}
	// Compute sorted order.
	initialSortOrder := []int{}
	for i := 0; i < len(data); i++ {
		initialSortOrder = append(initialSortOrder, i)
	}
	sorter := tableDataSorter{
		columnSortPriority: actualSortPriority,
		sortOrder:          initialSortOrder,
		data:               data,
	}
	sort.Sort(sorter)
	return sorter.sortOrder
}

type tableDataSorter struct {
	columnSortPriority []int
	sortOrder          []int
	data               [][]interface{}
}

func (s tableDataSorter) Len() int {
	return len(s.sortOrder)
}

func (s tableDataSorter) Less(i, j int) bool {
	iRow := s.data[s.sortOrder[i]]
	jRow := s.data[s.sortOrder[j]]
	for i := 0; i < len(s.columnSortPriority); i++ {
		columnIndex := s.columnSortPriority[i]
		if len(iRow) <= columnIndex && len(jRow) < columnIndex {
			continue
		}
		if len(iRow) <= columnIndex {
			return false
		}
		if len(jRow) <= columnIndex {
			return true
		}
		iValue := iRow[columnIndex]
		jValue := jRow[columnIndex]
		if iValue == jValue {
			continue
		}
		if iValue == nil {
			return false
		}
		if jValue == nil {
			return true
		}
		iNum, iNumOk := iValue.(int)
		jNum, jNumOk := jValue.(int)
		if iNumOk && jNumOk {
			return iNum > jNum
		}
		iFormatted := tableFormatCell(iValue)
		jFormatted := tableFormatCell(jValue)
		return iFormatted > jFormatted
	}
	return true
}

func (s tableDataSorter) Swap(i, j int) {
	s.sortOrder[i], s.sortOrder[j] = s.sortOrder[j], s.sortOrder[i]
}
