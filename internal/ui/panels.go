package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Типы панелей (как в LazyGit)
type PanelType int

const (
	TasksPanel PanelType = iota
	ProjectsPanel
	StatsPanel
)

// Модель задачи
type Task struct {
	ID          int
	Title       string
	Description string
	Status      string // pending, in_progress, done
	Priority    string // low, medium, high
	Project     string
}

// Модель проекта
type Project struct {
	ID    int
	Name  string
	Tasks []int // ID задач
}

// Данные для отображения
var (
	SampleTasks = []Task{
		{ID: 1, Title: "Implement main UI", Description: "Create the main layout", Status: "in_progress", Priority: "high", Project: "LazyTUI"},
		{ID: 2, Title: "Add task management", Description: "CRUD operations for tasks", Status: "pending", Priority: "medium", Project: "LazyTUI"},
		{ID: 3, Title: "Create statistics panel", Description: "Show project stats", Status: "pending", Priority: "low", Project: "LazyTUI"},
		{ID: 4, Title: "Write documentation", Description: "README and comments", Status: "done", Priority: "medium", Project: "Docs"},
		{ID: 5, Title: "Setup CI/CD", Description: "GitHub Actions", Status: "pending", Priority: "high", Project: "DevOps"},
	}

	SampleProjects = []Project{
		{ID: 1, Name: "LazyTUI", Tasks: []int{1, 2, 3}},
		{ID: 2, Name: "Docs", Tasks: []int{4}},
		{ID: 3, Name: "DevOps", Tasks: []int{5}},
	}
)

// Форматирование задачи для отображения
func (t Task) Format() string {
	var statusIcon string
	switch t.Status {
	case "done":
		statusIcon = StatusDoneStyle
	case "in_progress":
		statusIcon = StatusInProgressStyle
	default:
		statusIcon = StatusPendingStyle
	}

	var priorityIcon string
	priorityStyle := lipgloss.NewStyle().Foreground(LazyGitGreen)
	switch t.Priority {
	case "high":
		priorityIcon = "‼"
		priorityStyle = lipgloss.NewStyle().Foreground(LazyGitRed).Bold(true)
	case "medium":
		priorityIcon = "!"
		priorityStyle = lipgloss.NewStyle().Foreground(LazyGitYellow)
	default:
		priorityIcon = "•"
	}

	return fmt.Sprintf("%s %s %s %s",
		statusIcon,
		lipgloss.NewStyle().Foreground(LazyGitTextBold).Render(t.Title),
		priorityStyle.Render(priorityIcon),
		lipgloss.NewStyle().Foreground(LazyGitTextDim).Render(t.Project),
	)
}
