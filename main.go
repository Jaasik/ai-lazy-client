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
	// Left panels (5 panels)
	leftPanels [5]panelState
	
	// Right panels (2 panels)
	rightPanels [2]panelState

	// File list
	fileList       list.Model
	fileItems      []list.Item
	selectedFile   string
	fileContent    string
	
	// Text input (bottom right)
	textInput textinput.Model

	// Window dimensions
	width    int
	height   int
	quitting bool

	// Focus management (0-6: 5 left + 2 right panels)
	focusIndex int

	// Modal state
	modalOpen  bool
	modalFocus int // 0-3: 4 buttons
}

type panelState struct {
	focused bool
	title   string
	content string
}

func initialModel() model {
	// Load files from the files directory
	filesDir := forms.GetProjectFilesDir()
	items, err := forms.LoadFilesFromDir(filesDir)
	if err != nil || len(items) == 0 {
		items = []list.Item{
			forms.FileItem{Path: "no files found", Status: " "},
		}
	}

	l := forms.NewFileList(items)
	
	// Select first file by default
	var selectedFile string
	var fileContent string
	if len(items) > 0 {
		if fi, ok := items[0].(forms.FileItem); ok {
			selectedFile = fi.Path
			fileContent = forms.ReadFileContent(filesDir, selectedFile)
		}
	}

	ti := textinput.New()
	ti.Placeholder = "введите текст..."
	ti.CharLimit = 200
	ti.Width = 40
	ti.Prompt = "> "
	ti.TextStyle = lipgloss.NewStyle().Foreground(styles.ColorCyan)
	ti.PromptStyle = lipgloss.NewStyle().Foreground(styles.ColorGreen)

	// Initialize left panels (5 panels)
	leftPanels := [5]panelState{
		{title: "Действие", content: "Нажмите Enter"},
		{title: "Файлы", content: ""},
		{title: "Информация", content: "Панель 3"},
		{title: "Настройки", content: "Панель 4"},
		{title: "Статус", content: "Панель 5"},
	}

	// Initialize right panels (2 panels)
	rightPanels := [2]panelState{
		{title: "Содержимое файла", content: ""},
		{title: "Ввод текста", content: ""},
	}

	return model{
		leftPanels:   leftPanels,
		rightPanels:  rightPanels,
		fileList:     l,
		fileItems:    items,
		selectedFile: selectedFile,
		fileContent:  fileContent,
		textInput:    ti,
		focusIndex:   0,
		modalFocus:   0,
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
			case "esc":
				m.modalOpen = false
				return m, nil
			case "tab":
				m.modalFocus = (m.modalFocus + 1) % 4
				return m, nil
			case "shift+tab":
				m.modalFocus = (m.modalFocus - 1 + 4) % 4
				return m, nil
			case "left", "h":
				if m.modalFocus%2 == 0 {
					m.modalFocus = (m.modalFocus + 1) % 4
				} else {
					m.modalFocus = (m.modalFocus - 1 + 4) % 4
				}
				return m, nil
			case "right", "l":
				if m.modalFocus%2 == 1 {
					m.modalFocus = (m.modalFocus + 1) % 4
				} else {
					m.modalFocus = (m.modalFocus - 1 + 4) % 4
				}
				return m, nil
			case "up", "k":
				if m.modalFocus >= 2 {
					m.modalFocus -= 2
				} else {
					m.modalFocus += 2
				}
				return m, nil
			case "down", "j":
				if m.modalFocus < 2 {
					m.modalFocus += 2
				} else {
					m.modalFocus -= 2
				}
				return m, nil
			case "enter", " ":
				// Button action
				return m, nil
			}
		}
		return m, nil
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
			// Cycle through all 7 panels (5 left + 2 right)
			m.focusIndex = (m.focusIndex + 1) % 7
			m.textInput.Blur()
			return m, nil
		case "shift+tab":
			m.focusIndex = (m.focusIndex - 1 + 7) % 7
			m.textInput.Blur()
			return m, nil
		case "enter":
			// Open modal when pressing Enter on top-left panel (index 0)
			if m.focusIndex == 0 {
				m.modalOpen = true
				m.modalFocus = 0
				return m, nil
			}
			// Handle Enter on file list (index 1)
			if m.focusIndex == 1 {
				// Select file and load content
				if idx := m.fileList.Index(); idx < len(m.fileItems) {
					if fi, ok := m.fileItems[idx].(forms.FileItem); ok {
						m.selectedFile = fi.Path
						filesDir := forms.GetProjectFilesDir()
						m.fileContent = forms.ReadFileContent(filesDir, fi.Path)
						m.rightPanels[0].content = m.fileContent
					}
				}
				return m, nil
			}
		case "up", "k":
			if m.focusIndex == 1 {
				m.fileList, cmd = m.fileList.Update(msg)
				cmds = append(cmds, cmd)
				// Update file content on navigation
				if idx := m.fileList.Index(); idx < len(m.fileItems) {
					if fi, ok := m.fileItems[idx].(forms.FileItem); ok {
						m.selectedFile = fi.Path
						filesDir := forms.GetProjectFilesDir()
						m.fileContent = forms.ReadFileContent(filesDir, fi.Path)
						m.rightPanels[0].content = m.fileContent
					}
				}
				return m, tea.Batch(cmds...)
			}
		case "down", "j":
			if m.focusIndex == 1 {
				m.fileList, cmd = m.fileList.Update(msg)
				cmds = append(cmds, cmd)
				// Update file content on navigation
				if idx := m.fileList.Index(); idx < len(m.fileItems) {
					if fi, ok := m.fileItems[idx].(forms.FileItem); ok {
						m.selectedFile = fi.Path
						filesDir := forms.GetProjectFilesDir()
						m.fileContent = forms.ReadFileContent(filesDir, fi.Path)
						m.rightPanels[0].content = m.fileContent
					}
				}
				return m, tea.Batch(cmds...)
			}
		}
	}

	// Handle text input for bottom-right panel (index 6)
	if m.focusIndex == 6 {
		m.textInput.Focus()
		m.textInput, cmd = m.textInput.Update(msg)
		cmds = append(cmds, cmd)
		m.rightPanels[1].content = m.textInput.Value()
	} else {
		m.textInput.Blur()
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	// Calculate panel dimensions
	panelWidth := (m.width - 10) / 7 // 7 panels total with spacing
	if panelWidth < 30 {
		panelWidth = 30
	}
	panelHeight := (m.height - 10) / 5 // 5 rows for left side

	if panelHeight < 5 {
		panelHeight = 5
	}

	// Build left column (5 panels)
	var leftCol strings.Builder
	leftCol.WriteString(styles.HeaderStyle.Render("  ГЛАВНЫЙ ЭКРАН"))
	leftCol.WriteString("\n\n")

	for i := 0; i < 5; i++ {
		borderStyle := styles.UnfocusedBorderStyle
		if m.focusIndex == i {
			borderStyle = styles.FocusedBorderStyle
		}

		var contentLines []string
		if i == 1 {
			// File list panel
			listView := m.fileList.View()
			contentLines = strings.Split(strings.TrimRight(listView, "\n"), "\n")
		} else {
			contentLines = []string{m.leftPanels[i].content}
		}

		leftCol.WriteString(components.RenderPanel(borderStyle, panelWidth, panelHeight, m.leftPanels[i].title, contentLines, ""))
		leftCol.WriteString("\n")
	}

	// Build right column (2 panels)
	var rightCol strings.Builder
	rightCol.WriteString("\n\n\n\n") // Align with left header

	for i := 0; i < 2; i++ {
		borderStyle := styles.UnfocusedBorderStyle
		if m.focusIndex == 5+i {
			borderStyle = styles.FocusedBorderStyle
		}

		var content string
		if i == 0 {
			// File content with colors
			content = m.formatFileContent(m.fileContent)
		} else {
			// Text input
			content = m.textInput.View()
		}

		rightCol.WriteString(components.RenderPanel(borderStyle, panelWidth, panelHeight*2+2, m.rightPanels[i].title, []string{content}, ""))
		if i < 1 {
			rightCol.WriteString("\n")
		}
	}

	// Join columns
	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftCol.String(), rightCol.String())

	// Footer
	focusLabels := []string{"Действие", "Файлы", "Инфо", "Настр.", "Статус", "Содерж.", "Ввод"}
	focusLabel := "неизвестно"
	if m.focusIndex >= 0 && m.focusIndex < len(focusLabels) {
		focusLabel = focusLabels[m.focusIndex]
	}
	footer := fmt.Sprintf("\n\nФокус: %s • Tab: переключить • Enter: выбрать/открыть • q: выход", focusLabel)

	// Modal overlay
	if m.modalOpen {
		buttons := []string{"[1] OK", "[2] Применить", "[3] Отмена", "[4] Сброс"}
		modal := components.RenderModalWithButtons(styles.ModalWidth, "Выбор действия", buttons, m.modalFocus)
		overlay := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, modal)
		return overlay
	}

	return mainView + footer
}

// formatFileContent formats file content with different colors
func (m model) formatFileContent(content string) string {
	if content == "" {
		return lipgloss.NewStyle().Foreground(styles.ColorGray).Render("(нет содержимого)")
	}

	lines := strings.Split(content, "\n")
	var formatted []string
	
	colors := []lipgloss.Color{
		styles.ColorCyan,
		styles.ColorGreen,
		styles.ColorPink,
		styles.ColorBlue,
		styles.ColorPurple,
	}

	for i, line := range lines {
		color := colors[i%len(colors)]
		formatted = append(formatted, lipgloss.NewStyle().Foreground(color).Render(line))
	}

	return strings.Join(formatted, "\n")
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
