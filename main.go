package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"lazygit-newtest/components"
	"lazygit-newtest/forms"
	"lazygit-newtest/styles"
)

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
		forms.FileItem{Path: "1/Config.md", Status: "M"},
		forms.FileItem{Path: "2/commands/git.go", Status: "M"},
		forms.FileItem{Path: "3/commands/git_test.go", Status: "M"},
		forms.FileItem{Path: "4/config/app_config.go", Status: "M"},
		forms.FileItem{Path: "5/gui/gui.go", Status: "M"},
		forms.FileItem{Path: "6/Config.md", Status: "M"},
		forms.FileItem{Path: "7/commands/git.go", Status: "M"},
		forms.FileItem{Path: "8/commands/git_test.go", Status: "M"},
		forms.FileItem{Path: "9/config/app_config.go", Status: "M"},
		forms.FileItem{Path: "10/gui/gui.go", Status: "M"},
		forms.FileItem{Path: "11/Config.md", Status: "M"},
		forms.FileItem{Path: "12/commands/git.go", Status: "M"},
		forms.FileItem{Path: "13/commands/git_test.go", Status: "M"},
		forms.FileItem{Path: "14/config/app_config.go", Status: "M"},
		forms.FileItem{Path: "15/gui/gui.go", Status: "M"},
		forms.FileItem{Path: "16/Config.md", Status: "M"},
		forms.FileItem{Path: "17/commands/git.go", Status: "M"},
		forms.FileItem{Path: "18/commands/git_test.go", Status: "M"},
		forms.FileItem{Path: "19/config/app_config.go", Status: "M"},
		forms.FileItem{Path: "20/gui/gui.go", Status: "M"},
	}

	l := forms.NewFileList(items)

	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.CharLimit = 50
	ti.Width = styles.FixedBoxWidth - 8
	ti.Prompt = ""
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorPink)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorGray)

	modalItems := []list.Item{
		forms.ModalItem{Label: "Option A: Export"},
		forms.ModalItem{Label: "Option B: Import"},
		forms.ModalItem{Label: "Option C: Settings"},
		forms.ModalItem{Label: "Option D: Help"},
		forms.ModalItem{Label: "Option E: About"},
	}
	ml := forms.NewModalList(modalItems)

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
				if fi, ok := item.(forms.FileItem); ok {
					if strings.Contains(strings.ToLower(fi.Path), strings.ToLower(filterValue)) {
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

func (m model) View() string {
	if m.quitting {
		return ""
	}

	listBorder := styles.UnfocusedBorderStyle
	searchBorder := styles.UnfocusedBorderStyle
	rightBorder := styles.UnfocusedBorderStyle
	topPanelBorder := styles.UnfocusedBorderStyle

	switch m.focusIndex {
	case 0:
		topPanelBorder = styles.FocusedBorderStyle
	case 1:
		listBorder = styles.FocusedBorderStyle
	case 2:
		searchBorder = styles.FocusedBorderStyle
	}

	var leftCol strings.Builder
	leftCol.WriteString(styles.HeaderStyle.Render("lazygit -> newtest"))
	leftCol.WriteString("\n")

	leftCol.WriteString(components.RenderButtonPanel(topPanelBorder, styles.FixedBoxWidth, "Static text", "(X)", m.focusIndex == 0))
	leftCol.WriteString("\n")

	listView := m.list.View()
	listLines := strings.Split(strings.TrimRight(listView, "\n"), "\n")
	counterText := components.FormatCounter(m.list.Index(), len(m.list.Items()))
	leftCol.WriteString(components.RenderPanel(listBorder, styles.FixedBoxWidth, styles.PanelHeight+2, "Files", listLines, counterText))
	leftCol.WriteString("\n")

	searchLines := []string{}
	inputPrefix := "🔍 "
	inputValue := m.textInput.View()
	searchLines = append(searchLines, inputPrefix+inputValue)
	leftCol.WriteString(components.RenderPanel(searchBorder, styles.FixedBoxWidth, 3, "Search", searchLines, ""))

	var rightCol strings.Builder
	rightCol.WriteString("\n")

	previewContent := m.textInput.Value()
	if previewContent == "" {
		previewContent = styles.RightPanelStyle.Render("(waiting for input...)")
	}
	rightCol.WriteString(components.RenderPanel(rightBorder, styles.FixedBoxWidth, styles.PanelHeight, "Preview", []string{previewContent}, ""))
	rightCol.WriteString("\n")

	infoLines := components.FormatInfoLines(len(m.originalItems), len(m.list.Items()), m.list.Index()+1, m.buttonPressed)
	rightCol.WriteString(components.RenderPanel(rightBorder, styles.FixedBoxWidth, styles.PanelHeight, "Info", infoLines, ""))

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftCol.String(), rightCol.String())
	focusLabel := components.FormatFocusLabel(m.focusIndex)
	footer := fmt.Sprintf("\nFocus: %s • tab: switch • ↑/k ↓/j: nav • q: quit", focusLabel)

	if m.modalOpen {
		modalList := m.modalList.View()
		buttons := []string{"[1] OK", "[2] Apply", "[3] Cancel", "[4] Reset"}
		modal := components.RenderModal(styles.ModalWidth, modalList, buttons, m.modalFocus)
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
