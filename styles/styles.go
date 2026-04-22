package styles

import "github.com/charmbracelet/lipgloss"

// Color palette
const (
	ColorGreen      = lipgloss.Color("42")
	ColorPink       = lipgloss.Color("205")
	ColorCyan       = lipgloss.Color("51")
	ColorBlue       = lipgloss.Color("27")
	ColorPurple     = lipgloss.Color("90")
	ColorGray       = lipgloss.Color("240")
	ColorLightGray  = lipgloss.Color("252")
	ColorWhite      = lipgloss.Color("15")
	ColorDarkGreen  = lipgloss.Color("40")
)

// Dimensions
const (
	FixedBoxWidth = 44
	PanelHeight   = 10
	ModalWidth    = 50
)

// Base styles
var (
	// Status and selection
	StatusStyle   = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)
	SelectedStyle = lipgloss.NewStyle().Foreground(ColorGreen).Bold(true)

	// Headers
	HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(ColorPink)

	// Borders
	UnfocusedBorderStyle = lipgloss.NewStyle().Foreground(ColorGray)
	FocusedBorderStyle   = lipgloss.NewStyle().Foreground(ColorDarkGreen)

	// Right panel
	RightPanelStyle = lipgloss.NewStyle().Foreground(ColorLightGray)

	// Buttons
	ButtonStyle        = lipgloss.NewStyle().Foreground(ColorWhite).Background(ColorGray).Bold(true).Padding(0, 2)
	ButtonFocusedStyle = lipgloss.NewStyle().Foreground(ColorWhite).Background(ColorPurple).Bold(true).Padding(0, 2)
	ModalButtonStyle   = lipgloss.NewStyle().Foreground(ColorWhite).Background(ColorBlue).Bold(true).Padding(0, 2)
)

// Panel border configuration
type PanelBorder struct {
	Top         string
	Bottom      string
	Left        string
	Right       string
	TopLeft     string
	TopRight    string
	BottomLeft  string
	BottomRight string
}

var DefaultBorder = PanelBorder{
	Top: "─", Bottom: "─", Left: "│", Right: "│",
	TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
}
