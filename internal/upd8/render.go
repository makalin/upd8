package upd8

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	ansiReset     = "\033[0m"
	ansiCyan      = "\033[36m"
	ansiGreen     = "\033[32m"
	ansiHiMagenta = "\033[95m"
	ansiHiYellow  = "\033[93m"
	ansiRed       = "\033[31m"
)

var ansiRegexp = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// Renderer prints scan results in a human-friendly way.
type Renderer struct {
	Writer       io.Writer
	EnableColor  bool
	ShowPackages bool
	Timestamp    bool
	EmptyMessage string
}

func (r Renderer) Render(results []Result) {
	if r.Writer == nil {
		return
	}

	if len(results) == 0 {
		msg := r.EmptyMessage
		if msg == "" {
			msg = "No supported package managers detected."
		}
		fmt.Fprintln(r.Writer, msg)
		return
	}

	if r.Timestamp {
		fmt.Fprintf(r.Writer, "\n[%s]\n", time.Now().Format(time.RFC3339))
	}

	headers := []string{"Manager", "Outdated", "Packages", "Update Command"}
	rows := make([][]string, 0, len(results)+1)
	rows = append(rows, headers)

	totalOutdated := 0

	for _, res := range results {
		if res.Err != nil {
			rows = append(rows, []string{
				colorize(r.EnableColor, res.Manager, ansiHiYellow),
				colorize(r.EnableColor, "error", ansiRed),
				truncate(res.Err.Error(), 60),
				res.UpdateCommand,
			})
			continue
		}

		count := len(res.Items)
		totalOutdated += count

		countText := fmt.Sprintf("%d", count)
		if count == 0 {
			countText = colorize(r.EnableColor, "0", ansiGreen)
		} else {
			countText = colorize(r.EnableColor, fmt.Sprintf("%d", count), ansiHiMagenta)
		}

		pkgList := "—"
		if count > 0 && r.ShowPackages {
			pkgNames := make([]string, 0, min(count, 3))
			for idx, item := range res.Items {
				if idx >= 3 {
					break
				}
				pkgNames = append(pkgNames, item.Name)
			}
			if len(res.Items) > 3 {
				pkgNames = append(pkgNames, fmt.Sprintf("+%d more", len(res.Items)-3))
			}
			pkgList = strings.Join(pkgNames, ", ")
		}

		rows = append(rows, []string{
			colorize(r.EnableColor, res.Manager, ansiCyan) + " " + stopwatch(res.DurationMs),
			countText,
			pkgList,
			res.UpdateCommand,
		})
	}

	widths := computeColumnWidths(rows)

	printRow(r.Writer, rows[0], widths)
	printSeparator(r.Writer, widths)

	for i := 1; i < len(rows); i++ {
		printRow(r.Writer, rows[i], widths)
	}

	if totalOutdated == 0 {
		fmt.Fprintln(r.Writer, "\n🎉 All supported package managers look up to date.")
	}
}

func printSeparator(w io.Writer, widths []int) {
	parts := make([]string, len(widths))
	for i, width := range widths {
		parts[i] = strings.Repeat("-", width)
	}
	fmt.Fprintln(w, strings.Join(parts, "  "))
}

func printRow(w io.Writer, row []string, widths []int) {
	cells := make([]string, len(row))
	for i, cell := range row {
		padding := widths[i] - displayWidth(cell)
		if padding < 0 {
			padding = 0
		}
		cells[i] = cell + strings.Repeat(" ", padding)
	}
	fmt.Fprintln(w, strings.Join(cells, "  "))
}

func computeColumnWidths(rows [][]string) []int {
	widths := make([]int, len(rows[0]))
	for _, row := range rows {
		for i, cell := range row {
			w := displayWidth(cell)
			if w > widths[i] {
				widths[i] = w
			}
		}
	}
	return widths
}

func displayWidth(s string) int {
	plain := ansiRegexp.ReplaceAllString(s, "")
	return utf8.RuneCountInString(plain)
}

func colorize(enabled bool, text, code string) string {
	if !enabled || text == "" {
		return text
	}
	return code + text + ansiReset
}

func stopwatch(durationMs int64) string {
	if durationMs <= 0 {
		return ""
	}
	return fmt.Sprintf("(%dms)", durationMs)
}

func truncate(input string, max int) string {
	if len(input) <= max {
		return input
	}
	if max <= 3 {
		return input[:max]
	}
	return input[:max-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
