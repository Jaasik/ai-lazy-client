package ui

// Глобальные горячие клавиши в стиле LazyGit
const (
	KeyQuit    = "q"
	KeyUp      = "k"
	KeyDown    = "j"
	KeyLeft    = "h"
	KeyRight   = "l"
	KeyEnter   = "enter"
	KeyTab     = "tab"
	KeySpace   = " "
	KeyDelete  = "d"
	KeyEdit    = "e"
	KeyAdd     = "a"
	KeyHelp    = "?"
	KeyRefresh = "r"
	KeyFilter  = "/"
	KeyEscape  = "esc"
)

// Карта подсказок для футера
var HelpText = map[string]string{
	"Navigation": "↑/k ↓/j →/l ←/h",
	"Actions":    "⏎ enter • space • a(add) • d(delete) • e(edit)",
	"Global":     "tab • ?(help) • q(quit)",
}
