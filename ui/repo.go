package ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/rrenode/gonodedb/model"
)

func BuildRepoPanel(repo model.Repo) string {
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	lines := []string{
		fmt.Sprintf(" %s", cyan("Name:   ")+repo.Name),
		fmt.Sprintf(" %s", yellow("Alias:  ")+repo.Alias),
		fmt.Sprintf(" %s", green("ID:     ")+repo.ID),
		fmt.Sprintf(" %s", magenta("Path:   ")+repo.Path),
	}

	width := maxLineWidth(lines) + 2
	borderTop := "┌" + strings.Repeat("─", width) + "┐"
	borderBottom := "└" + strings.Repeat("─", width) + "┘"

	boxed := []string{borderTop}
	for _, line := range lines {
		boxed = append(boxed, fmt.Sprintf("│ %-*s │", width, line))
	}
	boxed = append(boxed, borderBottom)

	return strings.Join(boxed, "\n")
}

func maxLineWidth(lines []string) int {
	max := 0
	for _, line := range lines {
		plain := stripANSI(line)
		if len(plain) > max {
			max = len(plain)
		}
	}
	return max
}

// Removes ANSI color codes so spacing is calculated correctly
func stripANSI(s string) string {
	// Very simple ANSI strip: remove \x1b[...m
	var result []rune
	skip := false
	for i := 0; i < len(s); i++ {
		if s[i] == 0x1b && i+1 < len(s) && s[i+1] == '[' {
			skip = true
			i++
			continue
		}
		if skip && (s[i] >= 'a' && s[i] <= 'z' || s[i] >= 'A' && s[i] <= 'Z') {
			skip = false
			continue
		}
		if !skip {
			result = append(result, rune(s[i]))
		}
	}
	return string(result)
}
