package table

import ()

// Table holds table data that can be sorted and printed.
type Table struct {
	headers            []interface{}
	columnSortPriority []int
	data               [][]interface{}
}

// Headers adds column headers.
func (t *Table) Headers(data ...interface{}) {
	t.headers = data
}

// Add appends a new row of data.
func (t *Table) Add(data ...interface{}) {
	t.data = append(t.data, data)
}

// Sort sets the columns to sort by when formating.
func (t *Table) Sort(columnSortPriority ...int) {
	t.columnSortPriority = columnSortPriority
}

// String formats the table data to a string.
func (t *Table) String() string {
	columnWidths := calcColumnWidths(append(t.data, t.headers))
	sortOrder := calcSortOrder(t.data, t.columnSortPriority)
	formattedTable := formatHorizontalSeparator(columnWidths)
	formattedTable += formatRow(t.headers, columnWidths)
	formattedTable += formatHorizontalSeparator(columnWidths)
	for i := 0; i < len(sortOrder); i++ {
		formattedTable += formatRow(t.data[sortOrder[i]], columnWidths)
	}
	formattedTable += formatHorizontalSeparator(columnWidths)
	return formattedTable
}

func calcColumnWidths(data [][]interface{}) []int {
	columnWidths := []int{}
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i]); j++ {
			width := len(formatCell(data[i][j]))
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
