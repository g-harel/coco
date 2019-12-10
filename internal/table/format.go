package table

import (
	"fmt"
	"strings"
)

// Must be single characters.
const verticalSeparator = "|"
const horizontalSeparator = "-"
const intersectionSeparator = "+"

func formatHorizontalSeparator(columnWidths []int) string {
	row := []string{}
	for i := 0; i < len(columnWidths); i++ {
		row = append(row, strings.Repeat(horizontalSeparator, columnWidths[i]+2))
	}
	pre := intersectionSeparator
	post := intersectionSeparator + "\n"
	return pre + strings.Join(row, intersectionSeparator) + post
}

// Numbers are right-aligned, all other data types are left-aligned.
func formatRow(data []interface{}, columnWidths []int) string {
	row := []string{}
	for i := 0; i < len(columnWidths); i++ {
		value := ""
		number := false
		if i < len(data) {
			if _, ok := data[i].(int); ok {
				number = true
			}
			value = formatCell(data[i])
		}
		if number {
			row = append(row, fmt.Sprintf("%*v", columnWidths[i], value))
		} else {
			row = append(row, fmt.Sprintf("%-*v", columnWidths[i], value))
		}
	}
	pre := verticalSeparator + " "
	mid := " " + verticalSeparator + " "
	post := " " + verticalSeparator + "\n"
	return pre + strings.Join(row, mid) + post
}

// Integers are formatted with commas for each thousand. Nil is formatted as an
// empty string.
func formatCell(data interface{}) string {
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
