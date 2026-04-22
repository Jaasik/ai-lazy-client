package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	// Состояние
	activePanel   PanelType
	tasks         []Task
	projects      []Project
	tasksTable    table.Model
	projectsTable table.Model
	selectedTask  *Task
	selectedIndex int
	showHelp      bool
	width, height int
	filterMode    bool
	filterText    string
}

func NewModel() *Model {
	// Создаем колонки для таблицы задач
	taskColumns := []table.Column{
		{Title: "Status", Width: 8},
		{Title: "Task", Width: 40},
		{Title: "P", Width: 3},
		{Title: "Project", Width: 15},
	}

	// Конвертируем задачи в строки
	taskRows := []table.Row{}
	for _, task := range SampleTasks {
		statusIcon := "●"
		if task.Status == "done" {
			statusIcon = "✓"
		} else if task.Status == "in_progress" {
			statusIcon = "◉"
		}

		priorityIcon := "•"
		if task.Priority == "high" {
			priorityIcon = "‼"
		} else if task.Priority == "medium" {
			priorityIcon = "!"
		}

		taskRows = append(taskRows, table.Row{
			statusIcon,
			task.Title,
			priorityIcon,
			task.Project,
		})
	}

	// Создаем колонки для таблицы проектов
	projectColumns := []table.Column{
		{Title: "Project", Width: 30},
		{Title: "Tasks", Width: 10},
		{Title: "Progress", Width: 20},
	}

	projectRows := []table.Row{}
	for _, project := range SampleProjects {
		taskCount := len(project.Tasks)
		completed := 0
		for _, taskID := range project.Tasks {
			for _, task := range SampleTasks {
				if task.ID == taskID && task.Status == "done" {
					completed++
				}
			}
		}
		progress := fmt.Sprintf("%d/%d", completed, taskCount)
		projectRows = append(projectRows, table.Row{
			project.Name,
			fmt.Sprintf("%d", taskCount),
			progress,
		})
	}

	m := &Model{
		activePanel:   TasksPanel,
		tasks:         SampleTasks,
		projects:      SampleProjects,
		tasksTable:    NewTableComponent(taskColumns, taskRows, true),
		projectsTable: NewTableComponent(projectColumns, projectRows, false),
		selectedIndex: 0,
		showHelp:      false,
	}

	m.updateSelectedTask()
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	if m.showHelp {
		return m.handleHelpUpdate(msg)
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateTableWidths()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case KeyQuit, "ctrl+c":
			return m, tea.Quit

		case KeyHelp:
			m.showHelp = true
			return m, nil

		case KeyTab:
			m.switchPanel()
			return m, nil

		case KeyLeft:
			if m.activePanel > 0 {
				m.activePanel--
				m.updateFocus()
			}
			return m, nil

		case KeyRight:
			if m.activePanel < 2 {
				m.activePanel++
				m.updateFocus()
			}
			return m, nil

		case KeyUp, KeyDown:
			return m.handleNavigation(msg.String())

		case KeyEnter:
			return m.handleSelection()

		case KeySpace:
			return m.toggleTaskStatus()

		case KeyAdd:
			return m, m.addTask()

		case KeyDelete:
			return m, m.deleteTask()

		case KeyRefresh:
			return m, m.refreshData()
		}
	}

	// Обновляем активную таблицу
	if m.activePanel == TasksPanel {
		var cmd tea.Cmd
		m.tasksTable, cmd = m.tasksTable.Update(msg)
		cmds = append(cmds, cmd)
		m.updateSelectedTask()
	} else if m.activePanel == ProjectsPanel {
		var cmd tea.Cmd
		m.projectsTable, cmd = m.projectsTable.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.showHelp {
		return RenderHelpModal()
	}

	if m.width == 0 {
		return "Loading...\n"
	}

	// Верхняя панель с вкладками
	header := RenderHeader(m.activePanel)

	// Основной контент
	var content string
	if m.activePanel == TasksPanel {
		tasksStyle := ActivePanelStyle
		if m.activePanel != TasksPanel {
			tasksStyle = InactivePanelStyle
		}
		content = tasksStyle.Width(m.width - 50).Render(m.tasksTable.View())
	} else if m.activePanel == ProjectsPanel {
		projectsStyle := ActivePanelStyle
		if m.activePanel != ProjectsPanel {
			projectsStyle = InactivePanelStyle
		}
		content = projectsStyle.Width(m.width - 50).Render(m.projectsTable.View())
	} else if m.activePanel == StatsPanel {
		content = m.renderStatsPanel()
	}

	// Правая информационная панель
	infoPanel := RenderInfoPanel(m.selectedTask)

	// Собираем вместе
	mainContent := lipgloss.JoinHorizontal(lipgloss.Top, content, infoPanel)

	// Футер с подсказками
	footer := RenderFooter()

	// Полный интерфейс
	return lipgloss.JoinVertical(lipgloss.Top, header, mainContent, footer)
}

// Вспомогательные методы
func (m *Model) updateTableWidths() {
	tableWidth := m.width - 50 - 8 // вычитаем ширину info панели и отступы
	if tableWidth < 40 {
		tableWidth = 40
	}
	m.tasksTable.SetWidth(tableWidth)
	m.projectsTable.SetWidth(tableWidth)
}

func (m *Model) switchPanel() {
	m.activePanel = (m.activePanel + 1) % 3
	m.updateFocus()
}

func (m *Model) updateFocus() {
	m.tasksTable.Blur()
	m.projectsTable.Blur()

	if m.activePanel == TasksPanel {
		m.tasksTable.Focus()
	} else if m.activePanel == ProjectsPanel {
		m.projectsTable.Focus()
	}
}

func (m *Model) updateSelectedTask() {
	if len(m.tasks) > 0 && m.tasksTable.Cursor() < len(m.tasks) {
		m.selectedTask = &m.tasks[m.tasksTable.Cursor()]
	} else {
		m.selectedTask = nil
	}
}

func (m *Model) handleNavigation(key string) (tea.Model, tea.Cmd) {
	if m.activePanel == TasksPanel {
		if key == KeyUp {
			m.tasksTable.MoveUp(1)
		} else if key == KeyDown {
			m.tasksTable.MoveDown(1)
		}
		m.updateSelectedTask()
	} else if m.activePanel == ProjectsPanel {
		if key == KeyUp {
			m.projectsTable.MoveUp(1)
		} else if key == KeyDown {
			m.projectsTable.MoveDown(1)
		}
	}
	return m, nil
}

func (m *Model) handleSelection() (tea.Model, tea.Cmd) {
	if m.activePanel == TasksPanel && m.selectedTask != nil {
		// Показать детали задачи (можно открыть модальное окно)
		// Здесь можно добавить редактирование
	}
	return m, nil
}

func (m *Model) toggleTaskStatus() (tea.Model, tea.Cmd) {
	if m.activePanel == TasksPanel && m.selectedTask != nil {
		switch m.selectedTask.Status {
		case "pending":
			m.selectedTask.Status = "in_progress"
		case "in_progress":
			m.selectedTask.Status = "done"
		case "done":
			m.selectedTask.Status = "pending"
		}
		m.refreshTasksTable()
	}
	return m, nil
}

func (m *Model) addTask() tea.Cmd {
	// Здесь можно открыть форму для добавления задачи
	// Для простоты добавим тестовую задачу
	newTask := Task{
		ID:          len(m.tasks) + 1,
		Title:       "New Task",
		Description: "Task description",
		Status:      "pending",
		Priority:    "medium",
		Project:     "LazyTUI",
	}
	m.tasks = append(m.tasks, newTask)
	m.refreshTasksTable()
	return nil
}

func (m *Model) deleteTask() tea.Cmd {
	if m.activePanel == TasksPanel && m.selectedTask != nil {
		for i, task := range m.tasks {
			if task.ID == m.selectedTask.ID {
				m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
				break
			}
		}
		m.refreshTasksTable()
		m.updateSelectedTask()
	}
	return nil
}

func (m *Model) refreshData() tea.Cmd {
	m.refreshTasksTable()
	m.refreshProjectsTable()
	return nil
}

func (m *Model) refreshTasksTable() {
	rows := []table.Row{}
	for _, task := range m.tasks {
		var statusDisplay string
		if task.Status == "done" {
			statusDisplay = StatusDoneStyle
		} else if task.Status == "in_progress" {
			statusDisplay = StatusInProgressStyle
		} else {
			statusDisplay = StatusPendingStyle
		}

		var priorityIcon string
		if task.Priority == "high" {
			priorityIcon = "‼"
		} else if task.Priority == "medium" {
			priorityIcon = "!"
		} else {
			priorityIcon = "•"
		}

		rows = append(rows, table.Row{
			statusDisplay, //直接用字符串
			lipgloss.NewStyle().Foreground(LazyGitTextBold).Render(task.Title),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render(priorityIcon),
			lipgloss.NewStyle().Foreground(LazyGitCyan).Render(task.Project),
		})
	}
	m.tasksTable.SetRows(rows)
}

func (m *Model) refreshProjectsTable() {
	rows := []table.Row{}
	for _, project := range m.projects {
		taskCount := len(project.Tasks)
		completed := 0
		for _, taskID := range project.Tasks {
			for _, task := range m.tasks {
				if task.ID == taskID && task.Status == "done" {
					completed++
				}
			}
		}
		progress := fmt.Sprintf("%d/%d", completed, taskCount)

		// Определяем цвет прогресса
		progressStyle := lipgloss.NewStyle().Foreground(LazyGitGreen)
		if completed == 0 {
			progressStyle = lipgloss.NewStyle().Foreground(LazyGitRed)
		} else if completed < taskCount {
			progressStyle = lipgloss.NewStyle().Foreground(LazyGitYellow)
		}

		rows = append(rows, table.Row{
			lipgloss.NewStyle().Foreground(LazyGitCyan).Render(project.Name),
			lipgloss.NewStyle().Foreground(LazyGitTextDim).Render(fmt.Sprintf("%d", taskCount)),
			progressStyle.Render(progress),
		})
	}
	m.projectsTable.SetRows(rows)
}

func (m *Model) renderStatsPanel() string {
	total := len(m.tasks)
	if total == 0 {
		total = 1 // избегаем деления на ноль
	}
	completed := 0
	inProgress := 0
	pending := 0

	for _, task := range m.tasks {
		switch task.Status {
		case "done":
			completed++
		case "in_progress":
			inProgress++
		case "pending":
			pending++
		}
	}

	stats := lipgloss.JoinVertical(lipgloss.Top,
		lipgloss.NewStyle().Bold(true).Foreground(LazyGitBlue).Render("Statistics"),
		"",
		fmt.Sprintf("Total tasks:     %d", total),
		fmt.Sprintf("%s Completed:     %d", StatusDoneStyle, completed),
		fmt.Sprintf("%s In progress:   %d", StatusInProgressStyle, inProgress),
		fmt.Sprintf("%s Pending:       %d", StatusPendingStyle, pending),
		"",
		fmt.Sprintf("Completion:       %.1f%%", float64(completed)/float64(total)*100),
		"",
		"Progress:",
		RenderProgressBar(completed, total, 30),
	)

	return ActivePanelStyle.Width(m.width - 50).Render(stats)
}

func (m *Model) handleHelpUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg.(type) {
	case tea.KeyMsg:
		m.showHelp = false
		return m, nil
	}
	return m, nil
}
