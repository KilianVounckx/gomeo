package main

import (
	"strings"
)

func repeat(values []Value, number int) []Value {
	if number < 0 {
		return nil
	}

	res := make([]Value, len(values)*number)
	for i, value := range values {
		for j := 0; j < number; j++ {
			res[i+j*len(values)] = value
		}
	}
	return res
}

func remove(values []Value, index int) []Value {
	if index < 0 || index > len(values)-1 {
		return nil
	}

	res := make([]Value, len(values)-1)
	for i := 0; i < index; i++ {
		res[i] = values[i]
	}
	for i := index + 1; i < len(values); i++ {
		res[i-1] = values[i]
	}
	return res
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func stringWithArrows(text string, start, end *Position) string {
	if len(text) == 0 {
		return ""
	}
	if len(text) == 1 {
		return text + "\n^"
	}

	res := ""

	index_start := max(strings.IndexRune(reverse(text[0:start.index]), '\n'), 0)
	if index_start != 0 {
		index_start = len(text) - index_start + 1
	}
	index_end := strings.IndexRune(text[index_start+1:len(text)], '\n')
	if index_end < 0 {
		index_end = len(text)
	}

	line_count := end.line - start.line + 1
	for i := 0; i < line_count; i++ {
		if index_end <= index_start {
			continue
		}
		line := text[index_start:index_end]

		var column_start int
		if i == 0 {
			column_start = start.column
		} else {
			column_start = 0
		}

		var column_end int
		if i == line_count-1 {
			column_end = end.column
		} else {
			column_end = len(line) - 1
		}

		res += line + "\n"
		res += strings.Repeat(" ", column_start) + strings.Repeat("^", column_end-column_start)

		index_start = index_end
		if index_start < len(text) {
			index_end = strings.IndexRune(text[index_start+1:len(text)], '\n')
		} else {
			index_end = len(text)
		}
		if index_end < 0 {
			index_end = len(text)
		}
	}

	return strings.Replace(res, "\t", "", 0)
}
