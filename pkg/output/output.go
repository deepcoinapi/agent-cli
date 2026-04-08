// Package output provides JSON and table formatting helpers.
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// JSON pretty-prints v to stdout.
func JSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	enc.Encode(v)
}

// Table prints rows as an aligned table.
// columns are the map keys to display; headers are the column labels.
func Table(rows []map[string]any, columns, headers []string) {
	if len(rows) == 0 {
		fmt.Println("No data.")
		return
	}
	if columns == nil {
		// collect keys from first row
		for k := range rows[0] {
			columns = append(columns, k)
		}
		headers = columns
	}

	// compute column widths
	widths := make([]int, len(columns))
	for i, h := range headers {
		widths[i] = len(h)
	}
	cells := make([][]string, len(rows))
	for r, row := range rows {
		cells[r] = make([]string, len(columns))
		for c, col := range columns {
			cells[r][c] = fmt.Sprintf("%v", row[col])
			if len(cells[r][c]) > widths[c] {
				widths[c] = len(cells[r][c])
			}
		}
	}

	// print header
	parts := make([]string, len(headers))
	for i, h := range headers {
		parts[i] = fmt.Sprintf("%-*s", widths[i], h)
	}
	line := strings.Join(parts, "  ")
	fmt.Println(line)
	// separator
	sepParts := make([]string, len(headers))
	for i := range headers {
		sepParts[i] = strings.Repeat("-", widths[i])
	}
	fmt.Println(strings.Join(sepParts, "  "))

	// print rows
	for _, row := range cells {
		for i := range row {
			parts[i] = fmt.Sprintf("%-*s", widths[i], row[i])
		}
		fmt.Println(strings.Join(parts, "  "))
	}
}
