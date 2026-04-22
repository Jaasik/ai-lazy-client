package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lazygit-newtest/styles"
)

// RenderPanel renders a panel with border, title, content and optional footer
func RenderPanel(borderStyle lipgloss.Style, width, height int, title string, contentLines []string, footerText string) string {
	innerWidth := width - 4
	var b strings.Builder

	// Top border with title
	topLine := "┌─" + title
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

	// Content lines
	maxLines := height - 2
	for i := 0; i < maxLines; i++ {
		var line string
		if i < len(contentLines) {
			line = contentLines[i]
		}
		lineWidth := lipgloss.Width(line)
		pad := innerWidth - lineWidth
		if pad < 0 {
			pad = 0
			if lineWidth > innerWidth {
				line = line[:innerWidth]
			}
		}
		renderedLine := borderStyle.Render("│") + " " + line + strings.Repeat(" ", pad) + " " + borderStyle.Render("│")
		b.WriteString(renderedLine)
		b.WriteString("\n")
	}

	// Bottom border with optional footer
	bottomLine := "└" + strings.Repeat("─", width-2) + "┘"
	if footerText != "" {
		footerPos := width - lipgloss.Width(footerText) - 2
		if footerPos < 0 {
			footerPos = 0
		}
		bottomLine = "└" + strings.Repeat("─", footerPos) + footerText + "─┘"
	}
	b.WriteString(borderStyle.Render(bottomLine))
	return b.String()
}

// RenderButtonPanel renders a panel with a button
func RenderButtonPanel(borderStyle lipgloss.Style, width int, staticText string, buttonText string, focused bool) string {
	innerWidth := width - 4
	var b strings.Builder

	// Top border
	topLine := "┌─Action"
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

	// Button style
	var btnStyle lipgloss.Style
	if focused {
		btnStyle = styles.ButtonFocusedStyle
	} else {
		btnStyle = styles.ButtonStyle
	}
	renderedBtn := btnStyle.Render(buttonText)

	// Layout static text and button
	leftW := lipgloss.Width(staticText)
	rightW := lipgloss.Width(renderedBtn)
	spaceW := innerWidth - leftW - rightW
	if spaceW < 1 {
		spaceW = 1
	}

	line := borderStyle.Render("│") + " " + staticText + strings.Repeat(" ", spaceW) + renderedBtn + " " + borderStyle.Render("│")
	b.WriteString(line)
	b.WriteString("\n")

	// Bottom border
	bottomLine := "└" + strings.Repeat("─", width-2) + "┘"
	b.WriteString(borderStyle.Render(bottomLine))
	return b.String()
}

// RenderModal renders a modal dialog with list and buttons
func RenderModal(width int, listContent string, buttons []string, focusedBtn int) string {
	innerW := width - 2

	// Title
	title := lipgloss.NewStyle().
		Foreground(styles.ColorCyan).
		Bold(true).
		Width(innerW).
		Align(lipgloss.Center).
		Render("Modal Options")

	// List content (fixed 5 lines)
	listLines := strings.Split(strings.TrimRight(listContent, "\n"), "\n")
	for len(listLines) < 5 {
		listLines = append(listLines, "")
	}
	if len(listLines) > 5 {
		listLines = listLines[:5]
	}
	listBlock := lipgloss.NewStyle().Width(innerW).Render(strings.Join(listLines, "\n"))

	// Separator
	sep := strings.Repeat("─", innerW)

	// Buttons panel
	btnStyle := styles.ModalButtonStyle
	b1 := btnStyle.Render(buttons[0])
	b2 := btnStyle.Render(buttons[1])
	b3 := btnStyle.Render(buttons[2])
	b4 := btnStyle.Render(buttons[3])

	row1 := b1 + "  " + b2
	row2 := b3 + "  " + b4

	// Center rows
	pad1 := (innerW - lipgloss.Width(row1)) / 2
	if pad1 < 0 {
		pad1 = 0
	}
	r1 := strings.Repeat(" ", pad1) + row1 + strings.Repeat(" ", innerW-lipgloss.Width(row1)-pad1)

	pad2 := (innerW - lipgloss.Width(row2)) / 2
	if pad2 < 0 {
		pad2 = 0
	}
	r2 := strings.Repeat(" ", pad2) + row2 + strings.Repeat(" ", innerW-lipgloss.Width(row2)-pad2)

	// Button panel with vertical spacing
	emptyLine := strings.Repeat(" ", innerW)
	btnPanel := strings.Join([]string{
		emptyLine,
		r1,
		emptyLine,
		r2,
		emptyLine,
	}, "\n")

	// Assemble content
	content := strings.Join([]string{title, listBlock, sep, btnPanel}, "\n")

	// Pad each line to exact width
	var padded []string
	for _, line := range strings.Split(content, "\n") {
		w := lipgloss.Width(line)
		if w < innerW {
			line += strings.Repeat(" ", innerW-w)
		} else if w > innerW {
			line = line[:innerW]
		}
		padded = append(padded, line)
	}
	finalContent := strings.Join(padded, "\n")

	// Apply pink border
	boxBorder := lipgloss.Border{
		Top: "─", Bottom: "─", Left: "│", Right: "│",
		TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
	}

	return lipgloss.NewStyle().
		Border(boxBorder).
		BorderForeground(styles.ColorPink).
		Width(width).
		Render(finalContent)
}

// FormatCounter formats the list counter text
func FormatCounter(current, total int) string {
	return fmt.Sprintf("%d of %d", current+1, total)
}

// FormatInfoLines formats information lines for the info panel
func FormatInfoLines(totalFiles, filtered, selected int, buttonPressed bool) []string {
	btnStatus := "Button: inactive"
	if buttonPressed {
		btnStatus = "Button: ACTIVE"
	}
	return []string{
		fmt.Sprintf("Total files: %d", totalFiles),
		fmt.Sprintf("Filtered: %d", filtered),
		fmt.Sprintf("Selected: %d", selected),
		btnStatus,
		"",
		"Shortcuts:",
		"  tab/shift+tab - focus",
		"  enter/space   - open modal",
		"  esc           - clear/close",
		"  q             - quit",
	}
}

// FormatFocusLabel returns the current focus label
func FormatFocusLabel(focusIndex int) string {
	switch focusIndex {
	case 0:
		return "button"
	case 1:
		return "list"
	case 2:
		return "input"
	default:
		return "unknown"
	}
}
