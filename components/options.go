package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"lazygit-newtest/styles"
)

// RenderOptionsModal renders a modal with 2 separate frames and 4 buttons in 2 rows
func RenderOptionsModal(width int, listContent string, buttons []string, focusedBtn int) string {
	innerW := width - 2

	// Title
	title := lipgloss.NewStyle().
		Foreground(styles.ColorCyan).
		Bold(true).
		Width(innerW).
		Align(lipgloss.Center).
		Render("Options")

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

	// Buttons panel - render all buttons with consistent styling
	var renderedButtons []string
	for i, btnText := range buttons {
		var btnStyle lipgloss.Style
		if i+1 == focusedBtn {
			btnStyle = styles.ModalButtonStyle
		} else {
			btnStyle = styles.ButtonStyle
		}
		renderedButtons = append(renderedButtons, btnStyle.Render(btnText))
	}

	// Arrange buttons in 2 rows (2 buttons per row)
	row1 := renderedButtons[0] + "  " + renderedButtons[1]
	row2 := renderedButtons[2] + "  " + renderedButtons[3]

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

// RenderDualFrameOptions renders options module with 2 separate frames and 4 buttons in 2 rows
func RenderDualFrameOptions(width int, title1, title2 string, buttons []string, focusedElement int) string {
	// focusedElement: 0 = first frame list, 1-4 = buttons
	innerW := width - 4

	var b strings.Builder

	// First frame - top border
	topLine1 := "┌─" + title1
	padding1 := width - lipgloss.Width(topLine1) - 1
	topLine1 += strings.Repeat("─", padding1) + "┐"
	b.WriteString(styles.FocusedBorderStyle.Render(topLine1))
	b.WriteString("\n")

	// First frame - content area (3 lines for list items)
	for i := 0; i < 3; i++ {
		var line string
		if focusedElement == 0 && i == 0 {
			line = styles.SelectedStyle.Render("► Option " + string(rune('A'+i)))
		} else {
			line = "  Option " + string(rune('A'+i))
		}
		lineWidth := lipgloss.Width(line)
		pad := innerW - lineWidth
		if pad < 0 {
			pad = 0
		}
		renderedLine := styles.UnfocusedBorderStyle.Render("│") + " " + line + strings.Repeat(" ", pad) + " " + styles.UnfocusedBorderStyle.Render("│")
		b.WriteString(renderedLine)
		b.WriteString("\n")
	}

	// First frame - bottom border
	bottomLine1 := "└" + strings.Repeat("─", width-2) + "┘"
	b.WriteString(styles.UnfocusedBorderStyle.Render(bottomLine1))
	b.WriteString("\n")

	// Second frame - top border
	topLine2 := "┌─" + title2
	padding2 := width - lipgloss.Width(topLine2) - 1
	topLine2 += strings.Repeat("─", padding2) + "┐"
	b.WriteString(styles.UnfocusedBorderStyle.Render(topLine2))
	b.WriteString("\n")

	// Second frame - buttons area (2 rows of 2 buttons)
	// Row 1: buttons 0 and 1
	var btnStyles []lipgloss.Style
	for i := 0; i < 4; i++ {
		if i+1 == focusedElement {
			btnStyles = append(btnStyles, styles.ButtonFocusedStyle)
		} else {
			btnStyles = append(btnStyles, styles.ButtonStyle)
		}
	}

	renderedBtn0 := btnStyles[0].Render(buttons[0])
	renderedBtn1 := btnStyles[1].Render(buttons[1])
	renderedBtn2 := btnStyles[2].Render(buttons[2])
	renderedBtn3 := btnStyles[3].Render(buttons[3])

	// Row 1
	row1 := renderedBtn0 + "  " + renderedBtn1
	padRow1 := (innerW - lipgloss.Width(row1)) / 2
	if padRow1 < 0 {
		padRow1 = 0
	}
	row1Line := styles.UnfocusedBorderStyle.Render("│") + " " + strings.Repeat(" ", padRow1) + row1 + strings.Repeat(" ", innerW-lipgloss.Width(row1)-padRow1) + " " + styles.UnfocusedBorderStyle.Render("│")
	b.WriteString(row1Line)
	b.WriteString("\n")

	// Empty separator line
	emptyLine := styles.UnfocusedBorderStyle.Render("│") + " " + strings.Repeat(" ", innerW) + " " + styles.UnfocusedBorderStyle.Render("│")
	b.WriteString(emptyLine)
	b.WriteString("\n")

	// Row 2
	row2 := renderedBtn2 + "  " + renderedBtn3
	padRow2 := (innerW - lipgloss.Width(row2)) / 2
	if padRow2 < 0 {
		padRow2 = 0
	}
	row2Line := styles.UnfocusedBorderStyle.Render("│") + " " + strings.Repeat(" ", padRow2) + row2 + strings.Repeat(" ", innerW-lipgloss.Width(row2)-padRow2) + " " + styles.UnfocusedBorderStyle.Render("│")
	b.WriteString(row2Line)
	b.WriteString("\n")

	// Second frame - bottom border
	bottomLine2 := "└" + strings.Repeat("─", width-2) + "┘"
	b.WriteString(styles.UnfocusedBorderStyle.Render(bottomLine2))

	return b.String()
}
