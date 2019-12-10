package table

import (
	"sort"
)

// Columns are prioritized by the order their index appears in
// "columnSortPriority". Column indecies that are not provided are prioritized
// the lowest and from first to last. Integers are compared as numbers, but all
// other data types are compared as strings. Nil values are given the lowest
// possible order.
func calcSortOrder(data [][]interface{}, columnSortPriority []int) []int {
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
		iFormatted := formatCell(iValue)
		jFormatted := formatCell(jValue)
		return iFormatted < jFormatted
	}
	return true
}

func (s tableDataSorter) Swap(i, j int) {
	s.sortOrder[i], s.sortOrder[j] = s.sortOrder[j], s.sortOrder[i]
}
