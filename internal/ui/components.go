package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Компонент таблицы в стиле LazyGit
func NewTableComponent(columns []table.Column, rows []table.Row, focused bool) table.Model {
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(focused),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(LazyGitBgLight).
		BorderBottom(true).
		Foreground(LazyGitBlue).
		Bold(true).
		Background(LazyGitBgMain)

	s.Selected = CursorStyle
	s.Cell = lipgloss.NewStyle().Foreground(LazyGitText)

	t.SetStyles(s)

	return t
}

// Отрисовка заголовка в стиле LazyGit
func RenderHeader(activePanel PanelType) string {
	tabs := []string{"Tasks", "Projects", "Stats"}
	var renderedTabs []string

	for i, tab := range tabs {
		if PanelType(i) == activePanel {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(tab))
		} else {
			renderedTabs = append(renderedTabs, InactiveTabStyle.Render(tab))
		}
	}

	// Левая часть заголовка с названием
	title := HeaderStyle.Render(" LazyTasks ")

	// Правая часть с информацией (как в LazyGit показывают ветку)
	branchInfo := lipgloss.NewStyle().
		Background(LazyGitBgLight).
		Foreground(LazyGitYellow).
		Padding(0, 1).
		Render(" main ")

	// Собираем заголовок как в LazyGit
	headerLeft := lipgloss.JoinHorizontal(lipgloss.Top, title, branchInfo)
	headerRight := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)

	header := lipgloss.JoinHorizontal(lipgloss.Top, headerLeft, headerRight)

	return header
}

// Отрисовка футера в стиле LazyGit
func RenderFooter() string {
	var helpItems []string

	// Группы горячих клавиш как в LazyGit
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("↑/k"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("up"),
		),
	)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("↓/j"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("down"),
		),
	)
	helpItems = append(helpItems, SeparatorStyle)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("space"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("toggle"),
		),
	)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("a"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("add"),
		),
	)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("d"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("delete"),
		),
	)
	helpItems = append(helpItems, SeparatorStyle)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("tab"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("switch"),
		),
	)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("?"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("help"),
		),
	)
	helpItems = append(helpItems,
		fmt.Sprintf("%s %s",
			HelpKeyStyle.Render("q"),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("quit"),
		),
	)

	return FooterStyle.Render(strings.Join(helpItems, " "))
}

// Отрисовка информационной панели в стиле LazyGit
func RenderInfoPanel(selectedTask *Task) string {
	if selectedTask == nil {
		return InfoPanelStyle.Render(
			lipgloss.JoinVertical(lipgloss.Top,
				lipgloss.NewStyle().
					Foreground(LazyGitBlue).
					Bold(true).
					Render("No selection"),
				"",
				lipgloss.NewStyle().
					Foreground(LazyGitTextDim).
					Italic(true).
					Render("Select a task to see details"),
			),
		)
	}

	// Статус с соответствующей иконкой
	var statusIcon string
	switch selectedTask.Status {
	case "done":
		statusIcon = StatusDoneStyle
	case "in_progress":
		statusIcon = StatusInProgressStyle
	default:
		statusIcon = StatusPendingStyle
	}

	// Приоритет с цветом
	var priorityIcon string
	switch selectedTask.Priority {
	case "high":
		priorityIcon = PriorityHighStyle
	case "low":
		priorityIcon = PriorityLowStyle
	default:
		priorityIcon = PriorityMediumStyle
	}

	info := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(LazyGitBlue).
			Bold(true).
			Render("Task Details"),
		"",
		fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("Title:"),
			lipgloss.NewStyle().Foreground(LazyGitTextBold).Render(selectedTask.Title),
		),
		fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("Status:"),
			statusIcon,
		),
		fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("Priority:"),
			priorityIcon,
		),
		fmt.Sprintf("%s %s",
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render("Project:"),
			lipgloss.NewStyle().Foreground(LazyGitCyan).Render(selectedTask.Project),
		),
		"",
		lipgloss.NewStyle().
			Foreground(LazyGitTextDim).
			Italic(true).
			Render(selectedTask.Description),
	)

	return InfoPanelStyle.Render(info)
}

// Отрисовка модального окна помощи в стиле LazyGit
func RenderHelpModal() string {
	helpContent := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().
			Foreground(LazyGitBlue).
			Bold(true).
			Render("LazyTasks Help"),
		"",
		lipgloss.NewStyle().
			Foreground(LazyGitYellow).
			Bold(true).
			Render("Navigation:"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("↑/k"), "Move up"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("↓/j"), "Move down"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("←/h"), "Move left"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("→/l"), "Move right"),
		"",
		lipgloss.NewStyle().
			Foreground(LazyGitYellow).
			Bold(true).
			Render("Actions:"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("enter"), "Select/Open"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("space"), "Toggle status"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("a"), "Add new task"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("e"), "Edit task"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("d"), "Delete task"),
		"",
		lipgloss.NewStyle().
			Foreground(LazyGitYellow).
			Bold(true).
			Render("Global:"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("tab"), "Switch panels"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("/"), "Filter/Search"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("?"), "Show help"),
		fmt.Sprintf("  %s %s", HelpKeyStyle.Render("q"), "Quit"),
		"",
		lipgloss.NewStyle().
			Foreground(LazyGitTextDim).
			Italic(true).
			Render("Press any key to close"),
	)

	return HelpModalStyle.Render(helpContent)
}

// Визуализация прогресса в стиле LazyGit
func RenderProgressBar(current, total int, width int) string {
	if total == 0 {
		total = 1
	}

	percentage := float64(current) / float64(total)
	filled := int(float64(width) * percentage)

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)

	return lipgloss.NewStyle().
		Foreground(LazyGitGreen).
		Render(bar)
}
