package internal

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

// Sort sorts the table data by column.
// First provided column index will have highest priority. Numbers are compared
// as numbers, but all other data types are ordered as strings. Nil values are
// given the lowest possible order. Column indecies that are not provided are
// prioritized last and in the same order.
func (t *Table) Sort(columns ...int) {

}

// String formats the table data to a string.
// Numbers are right-aligned and formatted with commas.
func (t *Table) String() {

}
