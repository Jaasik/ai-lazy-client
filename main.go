package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	fixedBoxWidth = 44
	panelHeight   = 10
	modalWidth    = 50
)

var (
	statusStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	selectedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	headerStyle          = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	unfocusedBorderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedBorderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("40"))
	rightPanelStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	buttonStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("240")).Bold(true).Padding(0, 2)
	buttonFocusedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("90")).Bold(true).Padding(0, 2)
	// Единый стиль для всех кнопок модалки
	modalButtonStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("27")).Bold(true).Padding(0, 2)
)

type fileItem struct {
	path   string
	status string
}

func (i fileItem) Title() string       { return i.path }
func (i fileItem) Description() string { return "" }
func (i fileItem) FilterValue() string { return i.path }

type modalItem struct {
	label string
}

func (i modalItem) Title() string       { return i.label }
func (i modalItem) Description() string { return "" }
func (i modalItem) FilterValue() string { return i.label }

type customDelegate struct{}

func (d customDelegate) Height() int                               { return 1 }
func (d customDelegate) Spacing() int                              { return 0 }
func (d customDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d customDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(fileItem)
	if !ok {
		return
	}
	var line string
	if index == m.Index() {
		line = fmt.Sprintf(" %s %s", selectedStyle.Render("> "+item.status), selectedStyle.Render(item.path))
	} else {
		line = fmt.Sprintf(" %s %s", statusStyle.Render(item.status), item.path)
	}
	fmt.Fprint(w, line)
}

type modalDelegate struct{}

func (d modalDelegate) Height() int                               { return 1 }
func (d modalDelegate) Spacing() int                              { return 0 }
func (d modalDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d modalDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(modalItem)
	if !ok {
		return
	}
	var line string
	if index == m.Index() {
		line = selectedStyle.Render("► " + item.label)
	} else {
		line = "  " + item.label
	}
	fmt.Fprint(w, line)
}

type model struct {
	list          list.Model
	textInput     textinput.Model
	originalItems []list.Item
	width         int
	height        int
	quitting      bool
	focusIndex    int // 0 = кнопка, 1 = список, 2 = ввод
	buttonPressed bool

	// Modal state
	modalOpen  bool
	modalList  list.Model
	modalFocus int // 0 = список, 1-4 = кнопки
}

func initialModel() model {
	items := []list.Item{
		fileItem{path: "1/Config.md", status: "M"},
		fileItem{path: "2/commands/git.go", status: "M"},
		fileItem{path: "3/commands/git_test.go", status: "M"},
		fileItem{path: "4/config/app_config.go", status: "M"},
		fileItem{path: "5/gui/gui.go", status: "M"},
		fileItem{path: "6/Config.md", status: "M"},
		fileItem{path: "7/commands/git.go", status: "M"},
		fileItem{path: "8/commands/git_test.go", status: "M"},
		fileItem{path: "9/config/app_config.go", status: "M"},
		fileItem{path: "10/gui/gui.go", status: "M"},
		fileItem{path: "11/Config.md", status: "M"},
		fileItem{path: "12/commands/git.go", status: "M"},
		fileItem{path: "13/commands/git_test.go", status: "M"},
		fileItem{path: "14/config/app_config.go", status: "M"},
		fileItem{path: "15/gui/gui.go", status: "M"},
		fileItem{path: "16/Config.md", status: "M"},
		fileItem{path: "17/commands/git.go", status: "M"},
		fileItem{path: "18/commands/git_test.go", status: "M"},
		fileItem{path: "19/config/app_config.go", status: "M"},
		fileItem{path: "20/gui/gui.go", status: "M"},
	}

	l := list.New(items, customDelegate{}, fixedBoxWidth-4, panelHeight)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)

	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.CharLimit = 50
	ti.Width = fixedBoxWidth - 8
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	modalItems := []list.Item{
		modalItem{label: "Option A: Export"},
		modalItem{label: "Option B: Import"},
		modalItem{label: "Option C: Settings"},
		modalItem{label: "Option D: Help"},
		modalItem{label: "Option E: About"},
	}
	ml := list.New(modalItems, modalDelegate{}, modalWidth-4, 8)
	ml.SetShowStatusBar(false)
	ml.SetShowFilter(false)
	ml.SetShowHelp(false)
	ml.SetShowPagination(false)
	ml.SetShowTitle(false)

	return model{
		list:          l,
		textInput:     ti,
		originalItems: items,
		focusIndex:    1,
		modalList:     ml,
		modalFocus:    0,
	}
}

func (m model) Init() tea.Cmd { return textinput.Blink }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	if m.modalOpen {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q":
				m.modalOpen = false
				return m, nil
			case "tab":
				m.modalFocus = (m.modalFocus + 1) % 5
				return m, nil
			case "shift+tab":
				m.modalFocus = (m.modalFocus - 1 + 5) % 5
				return m, nil
			case "up", "down", "k", "j":
				if m.modalFocus == 0 {
					m.modalList, cmd = m.modalList.Update(msg)
					cmds = append(cmds, cmd)
				}
				return m, tea.Batch(cmds...)
			case "enter", " ":
				if m.modalFocus >= 1 && m.modalFocus <= 4 {
					return m, nil
				}
			}
		}
		m.modalList, cmd = m.modalList.Update(msg)
		cmds = append(cmds, cmd)
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "tab":
			m.focusIndex = (m.focusIndex + 1) % 3
			m.textInput.Blur()
			return m, nil
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + 3) % 3
			m.textInput.Blur()
			return m, nil
		case "enter", " ":
			if m.focusIndex == 0 {
				m.modalOpen = true
				m.buttonPressed = true
				return m, nil
			}
		case "up", "down", "k", "j":
			if m.focusIndex == 1 {
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
		}
	}

	if m.focusIndex == 2 {
		m.textInput.Focus()
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)

		filterValue := m.textInput.Value()
		if filterValue == "" {
			m.list.SetItems(m.originalItems)
		} else {
			var filtered []list.Item
			for _, item := range m.originalItems {
				if fi, ok := item.(fileItem); ok {
					if strings.Contains(strings.ToLower(fi.path), strings.ToLower(filterValue)) {
						filtered = append(filtered, item)
					}
				}
			}
			m.list.SetItems(filtered)
		}
	} else {
		m.textInput.Blur()
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func renderPanel(borderStyle lipgloss.Style, width, height int, title string, contentLines []string, footerText string) string {
	innerWidth := width - 4
	var b strings.Builder

	topLine := "┌─" + title
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

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

func renderButtonPanel(borderStyle lipgloss.Style, width int, staticText string, buttonText string, focused bool) string {
	innerWidth := width - 4
	var b strings.Builder

	topLine := "┌─Action"
	padding := width - lipgloss.Width(topLine) - 1
	topLine += strings.Repeat("─", padding) + "┐"
	b.WriteString(borderStyle.Render(topLine))
	b.WriteString("\n")

	var btnStyle lipgloss.Style
	if focused {
		btnStyle = buttonFocusedStyle
	} else {
		btnStyle = buttonStyle
	}
	renderedBtn := btnStyle.Render(buttonText)

	leftW := lipgloss.Width(staticText)
	rightW := lipgloss.Width(renderedBtn)
	spaceW := innerWidth - leftW - rightW
	if spaceW < 1 {
		spaceW = 1
	}

	line := borderStyle.Render("│") + " " + staticText + strings.Repeat(" ", spaceW) + renderedBtn + " " + borderStyle.Render("│")
	b.WriteString(line)
	b.WriteString("\n")

	bottomLine := "└" + strings.Repeat("─", width-2) + "┘"
	b.WriteString(borderStyle.Render(bottomLine))
	return b.String()
}

// renderModal переписан для гарантированной геометрии, одинаковых кнопок и чёткого разделения
func renderModal(width int, listContent string, buttons []string, focusedBtn int) string {
	pinkColor := lipgloss.Color("205")
	innerW := width - 2 // Ширина контента без учёта рамок

	// 1. Заголовок
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("51")).
		Bold(true).
		Width(innerW).
		Align(lipgloss.Center).
		Render("Modal Options")

	// 2. Список (фиксированно 5 строк)
	listLines := strings.Split(strings.TrimRight(listContent, "\n"), "\n")
	for len(listLines) < 5 {
		listLines = append(listLines, "")
	}
	if len(listLines) > 5 {
		listLines = listLines[:5]
	}
	listBlock := lipgloss.NewStyle().Width(innerW).Render(strings.Join(listLines, "\n"))

	// 3. Разделитель на всю ширину
	sep := strings.Repeat("─", innerW)

	// 4. Панель кнопок (отдельная, с отступами сверху/снизу и между рядами)
	btnStyle := modalButtonStyle // Единый стиль для всех
	b1 := btnStyle.Render(buttons[0])
	b2 := btnStyle.Render(buttons[1])
	b3 := btnStyle.Render(buttons[2])
	b4 := btnStyle.Render(buttons[3])

	row1 := b1 + "  " + b2
	row2 := b3 + "  " + b4

	// Центрирование рядов
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

	// Формирование панели кнопок с вертикальными отступами
	emptyLine := strings.Repeat(" ", innerW)
	btnPanel := strings.Join([]string{
		emptyLine, // Отступ сверху
		r1,
		emptyLine, // Расстояние между кнопками
		r2,
		emptyLine, // Отступ снизу
	}, "\n")

	// Сборка всего контента
	content := strings.Join([]string{title, listBlock, sep, btnPanel}, "\n")

	// Выравнивание каждой строки до exact innerW для корректного наложения рамки
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

	// Применение розовой рамки
	boxBorder := lipgloss.Border{
		Top: "─", Bottom: "─", Left: "│", Right: "│",
		TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
	}

	return lipgloss.NewStyle().
		Border(boxBorder).
		BorderForeground(pinkColor).
		Width(width).
		Render(finalContent)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	listBorder := unfocusedBorderStyle
	searchBorder := unfocusedBorderStyle
	rightBorder := unfocusedBorderStyle
	topPanelBorder := unfocusedBorderStyle

	switch m.focusIndex {
	case 0:
		topPanelBorder = focusedBorderStyle
	case 1:
		listBorder = focusedBorderStyle
	case 2:
		searchBorder = focusedBorderStyle
	}

	var leftCol strings.Builder
	leftCol.WriteString(headerStyle.Render("lazygit -> newtest"))
	leftCol.WriteString("\n")

	leftCol.WriteString(renderButtonPanel(topPanelBorder, fixedBoxWidth, "Static text", "(X)", m.focusIndex == 0))
	leftCol.WriteString("\n")

	listView := m.list.View()
	listLines := strings.Split(strings.TrimRight(listView, "\n"), "\n")
	counterText := fmt.Sprintf("%d of %d", m.list.Index()+1, len(m.list.Items()))
	leftCol.WriteString(renderPanel(listBorder, fixedBoxWidth, panelHeight+2, "Files", listLines, counterText))
	leftCol.WriteString("\n")

	searchLines := []string{}
	inputPrefix := "🔍 "
	inputValue := m.textInput.View()
	searchLines = append(searchLines, inputPrefix+inputValue)
	leftCol.WriteString(renderPanel(searchBorder, fixedBoxWidth, 3, "Search", searchLines, ""))

	var rightCol strings.Builder
	rightCol.WriteString("\n")

	previewContent := m.textInput.Value()
	if previewContent == "" {
		previewContent = rightPanelStyle.Render("(waiting for input...)")
	}
	rightCol.WriteString(renderPanel(rightBorder, fixedBoxWidth, panelHeight, "Preview", []string{previewContent}, ""))
	rightCol.WriteString("\n")

	btnStatus := "Button: inactive"
	if m.buttonPressed {
		btnStatus = "Button: ACTIVE"
	}
	infoLines := []string{
		fmt.Sprintf("Total files: %d", len(m.originalItems)),
		fmt.Sprintf("Filtered: %d", len(m.list.Items())),
		fmt.Sprintf("Selected: %d", m.list.Index()+1),
		btnStatus,
		"",
		"Shortcuts:",
		"  tab/shift+tab - focus",
		"  enter/space   - open modal",
		"  esc           - clear/close",
		"  q             - quit",
	}
	rightCol.WriteString(renderPanel(rightBorder, fixedBoxWidth, panelHeight, "Info", infoLines, ""))

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftCol.String(), rightCol.String())
	focusLabel := ""
	switch m.focusIndex {
	case 0:
		focusLabel = "button"
	case 1:
		focusLabel = "list"
	case 2:
		focusLabel = "input"
	}
	footer := fmt.Sprintf("\nFocus: %s • tab: switch • ↑/k ↓/j: nav • q: quit", focusLabel)

	if m.modalOpen {
		modalList := m.modalList.View()
		buttons := []string{"[1] OK", "[2] Apply", "[3] Cancel", "[4] Reset"}
		modal := renderModal(modalWidth, modalList, buttons, m.modalFocus)
		overlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
		return overlay
	}

	return mainView + footer
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
