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

// RenderDualButtonPanel renders a panel with two identical buttons side by side
func RenderDualButtonPanel(borderStyle lipgloss.Style, width int, btn1Text, btn2Text string, focusedBtn int) string {
	innerWidth := width - 4
	var b strings.Builder

	// Top border
	topLine := "┌─Options"
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

	// Button styles
	var btn1Style, btn2Style lipgloss.Style
	if focusedBtn == 1 {
		btn1Style = styles.ButtonFocusedStyle
		btn2Style = styles.ButtonStyle
	} else if focusedBtn == 2 {
		btn1Style = styles.ButtonStyle
		btn2Style = styles.ButtonFocusedStyle
	} else {
		btn1Style = styles.ButtonStyle
		btn2Style = styles.ButtonStyle
	}

	renderedBtn1 := btn1Style.Render(btn1Text)
	renderedBtn2 := btn2Style.Render(btn2Text)

	// Calculate spacing to center buttons
	totalBtnWidth := lipgloss.Width(renderedBtn1) + lipgloss.Width(renderedBtn2)
	spaceBetween := 2 // fixed space between buttons
	remainingSpace := innerWidth - totalBtnWidth - spaceBetween
	if remainingSpace < 0 {
		remainingSpace = 0
	}
	leftPad := remainingSpace / 2
	rightPad := remainingSpace - leftPad

	line := borderStyle.Render("│") + " " + strings.Repeat(" ", leftPad) + renderedBtn1 + strings.Repeat(" ", spaceBetween) + renderedBtn2 + strings.Repeat(" ", rightPad) + " " + borderStyle.Render("│")
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

// WrapText wraps text to fit within a given width, preserving word boundaries
func WrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		wordLen := len(word)
		currentLen := currentLine.Len()

		if currentLen == 0 {
			currentLine.WriteString(word)
		} else if currentLen+1+wordLen <= maxWidth {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// RenderScrollablePanel renders a panel with scrollable multi-line content
func RenderScrollablePanel(borderStyle lipgloss.Style, width, height int, title string, content string, scrollTop int) string {
	innerWidth := width - 4
	innerHeight := height - 2

	var b strings.Builder

	// Top border with title
	topLine := "┌─" + title
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

	// Wrap content to fit width
	wrappedLines := WrapText(content, innerWidth)

	// Apply scroll offset
	if scrollTop >= len(wrappedLines) {
		scrollTop = len(wrappedLines) - 1
	}
	if scrollTop < 0 {
		scrollTop = 0
	}

	visibleEnd := scrollTop + innerHeight
	if visibleEnd > len(wrappedLines) {
		visibleEnd = len(wrappedLines)
	}

	// Render visible lines
	for i := scrollTop; i < visibleEnd; i++ {
		line := wrappedLines[i]
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

	// Fill empty lines if content is shorter than panel height
	for i := visibleEnd - scrollTop; i < innerHeight; i++ {
		emptyLine := borderStyle.Render("│") + " " + strings.Repeat(" ", innerWidth) + " " + borderStyle.Render("│")
		b.WriteString(emptyLine)
		b.WriteString("\n")
	}

	// Bottom border with scroll indicator
	bottomLine := "└" + strings.Repeat("─", width-2) + "┘"
	if len(wrappedLines) > innerHeight {
		scrollIndicator := fmt.Sprintf(" [%d/%d]", scrollTop+1, len(wrappedLines))
		footerPos := width - lipgloss.Width(scrollIndicator) - 2
		if footerPos < 0 {
			footerPos = 0
		}
		bottomLine = "└" + strings.Repeat("─", footerPos) + scrollIndicator + "─┘"
	}
	b.WriteString(borderStyle.Render(bottomLine))

	return b.String()
}
