package util

import (
	"fmt"
	"io"
	"strings"
)

type Columnizer interface {
	Append(args ...string)
	Print(w io.Writer)
}

func NewColumnizer() Columnizer {
	return &columnizer{
		columnSeparator: "  ",
	}
}

type columnizer struct {
	rows            [][]string
	maxColumns      int
	columnSeparator string
}

func (c *columnizer) Append(args ...string) {
	c.rows = append(c.rows, args)
	c.maxColumns = max(c.maxColumns, len(args))

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c *columnizer) Print(w io.Writer) {
	maxColumnWidths := make([]int, c.maxColumns)
	columnFormat := make([]string, c.maxColumns)

	for _, r := range c.rows {
		for i, c := range r {
			maxColumnWidths[i] = max(maxColumnWidths[i], len(c))
		}
	}

	for i, w := range maxColumnWidths {
		if i == len(maxColumnWidths)-1 {
			columnFormat[i] = "%s"
		} else {
			columnFormat[i] = fmt.Sprintf("%%-%ds", w)

		}
	}

	for _, r := range c.rows {
		columns := len(r)
		f := strings.Join(columnFormat[:columns], c.columnSeparator)
		args := make([]interface{}, columns)
		for i, v := range r {
			args[i] = v
		}

		_, _ = fmt.Fprintf(w, f, args...)
		_, _ = fmt.Fprintln(w)
	}
}
