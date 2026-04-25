package component

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/mattn/go-runewidth"
)

// NOTE: full credit to broderick's overlay solution:
// ----------------------------------------------------------------
// https://gist.github.com/Broderick-Westrope/b89b14770c09dda928c4a108f437b927
// ----------------------------------------------------------------

// At its core this allows overlaying a string on top of another.
// My use case is for building TUI applications using "github.com/charmbracelet/bubbletea" wherein I like to have modal windows presented on top of the main window.
// This code was derived from the following source, but has the following changes:
//	- Wrapping is not done for you. Instead, wrapping of the background and overlay strings must be done beforehand.
// 	- A helper function is included for overlaying at the center of the background string.
//	- A boolean `ignoreMarginWhitespace`. When false, margin whitespace in the overlay string will overwrite the background string.
//		When true, margin whitespace in the overlay string will be ignored such that these cells of the background string are preserved. See the tests for examples.
//
// CREDIT: https://gist.github.com/ras0q/9bf5d81544b22302393f61206892e2cd

// OverlayCenter writes the overlay string onto the background string such that the middle of the
// overlay string will be at the middle of the overlay will be at the middle of the background.
func OverlayCenter(bg string, overlay string, ignoreMarginWhitespace bool) (string, error) {
	row := lipgloss.Height(bg) / 2
	row -= lipgloss.Height(overlay) / 2
	col := lipgloss.Width(bg) / 2
	col -= lipgloss.Width(overlay) / 2
	return Overlay(bg, overlay, row, col, ignoreMarginWhitespace)
}

// Overlay writes the overlay string onto the background string at the specified row and column.
// In this case, the row and column are zero indexed.
func Overlay(bg, overlay string, row, col int, ignoreMarginWhitespace bool) (string, error) {
	bgLines := strings.Split(bg, "\n")
	overlayLines := strings.Split(overlay, "\n")

	for i, overlayLine := range overlayLines {
		targetRow := i + row

		// Ensure the target row exists in the background lines
		for len(bgLines) <= targetRow {
			bgLines = append(bgLines, "")
		}

		bgLine := bgLines[targetRow]
		bgLineWidth := ansi.StringWidth(bgLine)

		if bgLineWidth < col {
			bgLine += strings.Repeat(" ", col-bgLineWidth) // Add padding
		}

		// Handle ignoreMarginWhitespace
		if ignoreMarginWhitespace {
			// Process the overlay line to preserve leading and trailing whitespace
			overlayLine = removeMarginWhitespace(bgLine, overlayLine, col)
		}

		bgLeft := ansi.Truncate(bgLine, col, "")
		bgRight, err := truncateLeft(bgLine, col+ansi.StringWidth(overlayLine))
		if err != nil {
			return "", fmt.Errorf("failed to truncate line: %w", err)
		}

		bgLines[targetRow] = bgLeft + overlayLine + bgRight
	}

	return strings.Join(bgLines, "\n"), nil
}

// removeMarginWhitespace preserves the background where the overlay line has leading or trailing whitespace.
// This is done by detecting those empty cells in the overlay string and replacing them with the corresponding background cells.
func removeMarginWhitespace(bgLine, overlayLine string, col int) string {
	var result strings.Builder

	// Variables to track ANSI escape sequences
	inAnsi := false
	ansiSeq := strings.Builder{}

	// Strip ANSI codes to analyze whitespace
	overlayStripped := ansi.Strip(overlayLine)
	overlayRunes := []rune(overlayStripped)

	// Find first and last non-whitespace positions
	firstNonWhitespacePos := -1
	lastNonWhitespacePos := -1
	visualPos := 0
	overlayVisualWidths := make([]int, len(overlayRunes))

	for i, r := range overlayRunes {
		runeWidth := runewidth.RuneWidth(r)
		overlayVisualWidths[i] = runeWidth
		if !unicode.IsSpace(r) {
			if firstNonWhitespacePos == -1 {
				firstNonWhitespacePos = visualPos
			}
			lastNonWhitespacePos = visualPos + runeWidth - 1 // inclusive
		}
		visualPos += runeWidth
	}

	// If all characters are whitespace
	if firstNonWhitespacePos == -1 {
		firstNonWhitespacePos = 0
		lastNonWhitespacePos = -1
	}

	// Now, process the overlayLine, keeping track of visual positions
	visualPos = 0
	runeReader := strings.NewReader(overlayLine)

	for {
		r, _, err := runeReader.ReadRune()
		if err != nil {
			break
		}

		if r == '\x1b' {
			// Start of ANSI escape sequence
			inAnsi = true
			ansiSeq.WriteRune(r)
			continue
		}

		if inAnsi {
			ansiSeq.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				// End of ANSI escape sequence
				inAnsi = false
				result.WriteString(ansiSeq.String())
				ansiSeq.Reset()
			}
			continue
		}

		runeWidth := runewidth.RuneWidth(r)

		// Determine if current position is leading whitespace or trailing whitespace
		var isLeadingWhitespace, isTrailingWhitespace bool

		if visualPos < firstNonWhitespacePos {
			isLeadingWhitespace = true
		} else if visualPos > lastNonWhitespacePos {
			isTrailingWhitespace = true
		}

		if unicode.IsSpace(r) && (isLeadingWhitespace || isTrailingWhitespace) {
			// Preserve background character
			for k := 0; k < runeWidth; k++ {
				bgChar := getBgCharAt(bgLine, col+visualPos+k)
				result.WriteString(bgChar)
			}
		} else {
			// Include character from overlay (could be a non-whitespace or whitespace character in between)
			result.WriteRune(r)
		}

		visualPos += runeWidth
	}

	return result.String()
}

// getBgCharAt returns the character from the background line at the specified visual index.
func getBgCharAt(bgLine string, visualIndex int) string {
	var result strings.Builder
	displayWidth := 0
	inAnsi := false
	ansiSeq := strings.Builder{}

	runeReader := strings.NewReader(bgLine)
	for {
		r, _, err := runeReader.ReadRune()
		if err != nil {
			break
		}

		if r == '\x1b' {
			// Start of ANSI escape sequence
			inAnsi = true
			ansiSeq.WriteRune(r)
			continue
		}

		if inAnsi {
			ansiSeq.WriteRune(r)
			if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
				// End of ANSI escape sequence
				inAnsi = false
				result.WriteString(ansiSeq.String())
				ansiSeq.Reset()
			}
			continue
		}

		charWidth := runewidth.RuneWidth(r)
		if displayWidth+charWidth > visualIndex {
			// We have reached the desired index
			result.WriteRune(r)
			break
		}

		displayWidth += charWidth
	}

	// If no character found at the position, return a space
	if result.Len() == 0 {
		return " "
	}

	return result.String()
}

// truncateLeft removes characters from the beginning of a line, considering ANSI escape codes.
func truncateLeft(line string, padding int) (string, error) {
	if strings.Contains(line, "\n") {
		return "", fmt.Errorf("line must not contain newline")
	}

	wrapped := strings.Split(ansi.Hardwrap(line, padding, true), "\n")
	if len(wrapped) == 1 {
		return "", nil
	}

	var ansiStyle string
	// Regular expression to match ANSI escape codes.
	ansiStyles := regexp.MustCompile(`\x1b[[\d;]*m`).FindAllString(wrapped[0], -1)
	if len(ansiStyles) > 0 {
		ansiStyle = ansiStyles[len(ansiStyles)-1]
	}

	return ansiStyle + strings.Join(wrapped[1:], ""), nil
}
