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

	// Modal state
	modalOpen   bool
	modalList   list.Model
	modalFocus  int // 0 = список, 1-4 = кнопки
	optionsMode bool  // true = options module with 2 frames, false = simple modal

	// Preview panel scroll state
	previewScroll int
}

func initialModel() model {
	// Load files from the files directory
	filesDir := forms.GetProjectFilesDir()
	items, err := forms.LoadFilesFromDir(filesDir)
	if err != nil || len(items) == 0 {
		// Fallback to default items if directory is empty or doesn't exist
		items = []list.Item{
			forms.FileItem{Path: "1/Config.md", Status: "M"},
			forms.FileItem{Path: "2/commands/git.go", Status: "M"},
			forms.FileItem{Path: "3/commands/git_test.go", Status: "M"},
			forms.FileItem{Path: "4/config/app_config.go", Status: "M"},
			forms.FileItem{Path: "5/gui/gui.go", Status: "M"},
		}
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
		optionsMode:   false,
		previewScroll: 0,
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
				m.optionsMode = false
				return m, nil
			case "tab":
				if m.optionsMode {
					// Options mode: 0 = first frame list, 1-4 = buttons
					m.modalFocus = (m.modalFocus + 1) % 5
				} else {
					// Simple modal mode: 0 = list, 1-4 = buttons
					m.modalFocus = (m.modalFocus + 1) % 5
				}
				return m, nil
			case "shift+tab":
				if m.optionsMode {
					m.modalFocus = (m.modalFocus - 1 + 5) % 5
				} else {
					m.modalFocus = (m.modalFocus - 1 + 5) % 5
				}
				return m, nil
			case "up", "k":
				if m.optionsMode {
					if m.modalFocus == 0 {
						// Navigate in first frame list (just visual, no actual list)
						return m, nil
					} else if m.modalFocus >= 1 && m.modalFocus <= 4 {
						// Navigate buttons vertically
						if m.modalFocus > 2 {
							m.modalFocus -= 2
						} else {
							m.modalFocus += 2
						}
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
						if m.modalFocus < 1 {
							m.modalFocus = 1
						}
					}
				} else {
					if m.modalFocus == 0 {
						m.modalList, cmd = m.modalList.Update(msg)
						cmds = append(cmds, cmd)
					} else if m.modalFocus >= 1 && m.modalFocus <= 4 {
						if m.modalFocus > 2 {
							m.modalFocus -= 2
						} else {
							m.modalFocus += 2
						}
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
						if m.modalFocus < 1 {
							m.modalFocus = 1
						}
					}
				}
				return m, tea.Batch(cmds...)
			case "down", "j":
				if m.optionsMode {
					if m.modalFocus == 0 {
						return m, nil
					} else if m.modalFocus >= 1 && m.modalFocus <= 4 {
						if m.modalFocus <= 2 {
							m.modalFocus += 2
						} else {
							m.modalFocus -= 2
						}
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
						if m.modalFocus < 1 {
							m.modalFocus = 1
						}
					}
				} else {
					if m.modalFocus == 0 {
						m.modalList, cmd = m.modalList.Update(msg)
						cmds = append(cmds, cmd)
					} else if m.modalFocus >= 1 && m.modalFocus <= 4 {
						if m.modalFocus <= 2 {
							m.modalFocus += 2
						} else {
							m.modalFocus -= 2
						}
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
						if m.modalFocus < 1 {
							m.modalFocus = 1
						}
					}
				}
				return m, tea.Batch(cmds...)
			case "left", "h":
				if m.modalFocus >= 1 && m.modalFocus <= 4 {
					if m.modalFocus%2 == 0 {
						m.modalFocus--
					} else {
						m.modalFocus++
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
					}
				}
				return m, nil
			case "right", "l":
				if m.modalFocus >= 1 && m.modalFocus <= 4 {
					if m.modalFocus%2 == 1 {
						m.modalFocus++
						if m.modalFocus > 4 {
							m.modalFocus = 4
						}
					} else {
						m.modalFocus--
					}
				}
				return m, nil
			case "enter", " ":
				if m.modalFocus >= 1 && m.modalFocus <= 4 {
					return m, nil
				}
			}
		}
		if !m.optionsMode {
			m.modalList, cmd = m.modalList.Update(msg)
			cmds = append(cmds, cmd)
		}
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
				m.optionsMode = true // Use new options module with 2 frames
				return m, nil
			}
		case "up", "k":
			if m.focusIndex == 1 {
				m.list, cmd = m.list.Update(msg)
				cmds = append(cmds, cmd)
				return m, tea.Batch(cmds...)
			}
			// Scroll preview panel with up arrow when not in other focus areas
			if m.focusIndex == 2 && m.previewScroll > 0 {
				m.previewScroll--
				return m, nil
			}
		case "down", "j":
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

		// Scroll preview with up/down when typing long text
		if msg, ok := msg.(tea.KeyMsg); ok {
			if msg.String() == "up" && m.previewScroll > 0 {
				m.previewScroll--
			}
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
		rightBorder = styles.FocusedBorderStyle
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

	// Render Preview panel with scrollable multi-line text
	previewContent := m.textInput.Value()
	if previewContent == "" {
		previewContent = "(waiting for input...)"
	}
	rightCol.WriteString(components.RenderScrollablePanel(rightBorder, styles.FixedBoxWidth, styles.PanelHeight, "Preview", previewContent, m.previewScroll))
	rightCol.WriteString("\n")

	infoLines := components.FormatInfoLines(len(m.originalItems), len(m.list.Items()), m.list.Index()+1)
	rightCol.WriteString(components.RenderPanel(rightBorder, styles.FixedBoxWidth, styles.PanelHeight, "Info", infoLines, ""))

	mainView := lipgloss.JoinHorizontal(lipgloss.Top, leftCol.String(), rightCol.String())
	focusLabel := components.FormatFocusLabel(m.focusIndex)
	footer := fmt.Sprintf("\nFocus: %s • tab: switch • ↑/k ↓/j: nav • q: quit", focusLabel)

	if m.modalOpen {
		buttons := []string{"[1] OK", "[2] Apply", "[3] Cancel", "[4] Reset"}
		var modal string
		if m.optionsMode {
			// Use new options module with 2 separate frames and 4 buttons in 2 rows
			modal = components.RenderDualFrameOptions(styles.ModalWidth, "Options", "Actions", buttons, m.modalFocus)
		} else {
			modalList := m.modalList.View()
			modal = components.RenderModal(styles.ModalWidth, modalList, buttons, m.modalFocus)
		}
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
