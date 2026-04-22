package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Точная цветовая схема LazyGit (из их темы по умолчанию)
var (
	// Основные цвета из LazyGit
	LazyGitBlue   = lipgloss.Color("39")  // Основной синий для выделения
	LazyGitCyan   = lipgloss.Color("37")  // Бирюзовый для второстепенных элементов
	LazyGitGreen  = lipgloss.Color("34")  // Зеленый для успешных операций
	LazyGitYellow = lipgloss.Color("220") // Желтый для предупреждений
	LazyGitRed    = lipgloss.Color("196") // Красный для ошибок
	LazyGitPurple = lipgloss.Color("99")  // Фиолетовый для особых акцентов
	LazyGitOrange = lipgloss.Color("208") // Оранжевый для важного

	// Фоновые цвета LazyGit
	LazyGitBgMain   = lipgloss.Color("235") // Темно-серый фон (почти черный)
	LazyGitBgLight  = lipgloss.Color("236") // Светло-серый фон для панелей
	LazyGitBgHover  = lipgloss.Color("238") // Фон при наведении
	LazyGitBgActive = lipgloss.Color("240") // Фон активного элемента

	// Цвета текста
	LazyGitText     = lipgloss.Color("252") // Основной текст (светло-серый)
	LazyGitTextDim  = lipgloss.Color("245") // Приглушенный текст
	LazyGitTextBold = lipgloss.Color("231") // Яркий текст (почти белый)
)

// Стили компонентов (как в LazyGit)
var (
	// Общий стиль приложения
	AppStyle = lipgloss.NewStyle().
			Background(LazyGitBgMain).
			Foreground(LazyGitText).
			Padding(1, 2)

	// Заголовок (как статус бар в LazyGit)
	HeaderStyle = lipgloss.NewStyle().
			Background(LazyGitBlue).
			Foreground(lipgloss.Color("0")). // Черный текст на синем фоне
			Bold(true).
			Padding(0, 1)

	// Активная панель (как выделенный файл в LazyGit)
	ActivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(LazyGitBlue).
				BorderBackground(LazyGitBgMain).
				Background(LazyGitBgMain).
				Padding(0, 1).
				Margin(0, 1)

	// Неактивная панель
	InactivePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(LazyGitBgLight).
				BorderBackground(LazyGitBgMain).
				Background(LazyGitBgMain).
				Padding(0, 1).
				Margin(0, 1)

	// Стиль для выделенной строки (курсор) - как в LazyGit
	CursorStyle = lipgloss.NewStyle().
			Background(LazyGitBgActive).
			Foreground(LazyGitTextBold).
			Bold(true)

	// Статусная строка внизу (как в LazyGit)
	FooterStyle = lipgloss.NewStyle().
			Background(LazyGitBgLight).
			Foreground(LazyGitTextDim).
			Padding(0, 1)

	// Информационная панель (правая панель в LazyGit)
	InfoPanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(LazyGitBgLight).
			Background(LazyGitBgMain).
			Foreground(LazyGitText).
			Padding(1).
			Width(40)

	// Стили вкладок (как вкладки в LazyGit: Files, Commits, Branches)
	ActiveTabStyle = lipgloss.NewStyle().
			Background(LazyGitBlue).
			Foreground(lipgloss.Color("0")).
			Padding(0, 3).
			Bold(true)

	InactiveTabStyle = lipgloss.NewStyle().
				Background(LazyGitBgLight).
				Foreground(LazyGitTextDim).
				Padding(0, 3)

	// Стили для приоритетов задач (это строки, а не стили!)
	PriorityHighStyle   = lipgloss.NewStyle().Foreground(LazyGitRed).Bold(true).Render("‼")
	PriorityMediumStyle = lipgloss.NewStyle().Foreground(LazyGitYellow).Render("!")
	PriorityLowStyle    = lipgloss.NewStyle().Foreground(LazyGitGreen).Render("•")

	// Статусы задач (иконки как в LazyGit) - это строки!
	StatusDoneStyle       = lipgloss.NewStyle().Foreground(LazyGitGreen).Bold(true).Render("✓")
	StatusPendingStyle    = lipgloss.NewStyle().Foreground(LazyGitYellow).Render("●")
	StatusInProgressStyle = lipgloss.NewStyle().Foreground(LazyGitCyan).Bold(true).Render("◉")

	// Стиль для модального окна помощи
	HelpModalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(LazyGitBlue).
			Background(LazyGitBgMain).
			Foreground(LazyGitText).
			Padding(1, 2).
			Width(60)
)

// Дополнительные вспомогательные стили
var (
	// Стиль для ключей в хелпе (как в LazyGit)
	HelpKeyStyle = lipgloss.NewStyle().
			Background(LazyGitBgLight).
			Foreground(LazyGitBlue).
			Bold(true).
			Padding(0, 1)

	// Стиль для разделителей
	SeparatorStyle = lipgloss.NewStyle().
			Foreground(LazyGitTextDim).
			Render(" • ")
)
