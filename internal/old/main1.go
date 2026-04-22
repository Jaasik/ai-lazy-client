package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const fixedBoxWidth = 44

var (
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	headerStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	borderStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type fileItem struct {
	path   string
	status string
}

func (i fileItem) Title() string       { return i.path }
func (i fileItem) Description() string { return "" }
func (i fileItem) FilterValue() string { return i.path }

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
		line = fmt.Sprintf(" %s %s",
			selectedStyle.Render("> "+item.status),
			selectedStyle.Render(item.path))
	} else {
		line = fmt.Sprintf(" %s %s",
			statusStyle.Render(item.status),
			item.path)
	}
	fmt.Fprint(w, line)
}

type model struct {
	list     list.Model
	width    int
	height   int
	quitting bool
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

	l := list.New(items, customDelegate{}, fixedBoxWidth-5, 10)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)

	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	var b strings.Builder

	b.WriteString(headerStyle.Render("lazygit -> newtest"))
	b.WriteString("\n\n")

	listView := m.list.View()
	lines := strings.Split(strings.TrimRight(listView, "\n"), "\n")

	// Top border
	topPrefix := "┌──Files"
	topBorder := topPrefix + strings.Repeat("─", fixedBoxWidth-len(topPrefix)-3) + "┐"
	b.WriteString(borderStyle.Render(topBorder))
	b.WriteString("\n")

	// Middle lines
	for _, line := range lines {
		lineWidth := lipgloss.Width(line)
		padding := fixedBoxWidth - 4 - lineWidth
		if padding < 0 {
			padding = 0
		}
		leftBorder := borderStyle.Render("│")
		rightBorder := borderStyle.Render("│")
		middleLine := leftBorder + " " + line + strings.Repeat(" ", padding-8) + " " + rightBorder
		b.WriteString(middleLine)
		b.WriteString("\n")
	}

	// Bottom border
	counterText := fmt.Sprintf("%d of %d", m.list.Index()+1, len(m.list.Items()))
	bottomPadding := fixedBoxWidth - len(counterText) - 2
	bottomBorder := "└" + strings.Repeat("─", bottomPadding-9) + counterText + "─┘"
	b.WriteString(borderStyle.Render(bottomBorder))

	b.WriteString("\n\n")
	b.WriteString("↑/k up • ↓/j down • q quit")

	return b.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
