package forms

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"lazygit-newtest/styles"
)

// FileItem represents a file in the list
type FileItem struct {
	Path   string
	Status string
}

func (i FileItem) Title() string       { return i.Path }
func (i FileItem) Description() string { return "" }
func (i FileItem) FilterValue() string { return i.Path }

// ModalItem represents an item in the modal list
type ModalItem struct {
	Label string
}

func (i ModalItem) Title() string       { return i.Label }
func (i ModalItem) Description() string { return "" }
func (i ModalItem) FilterValue() string { return i.Label }

// CustomDelegate renders file items in the list
type CustomDelegate struct{}

func (d CustomDelegate) Height() int                               { return 1 }
func (d CustomDelegate) Spacing() int                              { return 0 }
func (d CustomDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d CustomDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(FileItem)
	if !ok {
		return
	}
	var line string
	if index == m.Index() {
		line = fmt.Sprintf(" %s %s",
			styles.SelectedStyle.Render("> "+item.Status),
			styles.SelectedStyle.Render(item.Path))
	} else {
		line = fmt.Sprintf(" %s %s",
			styles.StatusStyle.Render(item.Status),
			item.Path)
	}
	fmt.Fprint(w, line)
}

// ModalDelegate renders modal items
type ModalDelegate struct{}

func (d ModalDelegate) Height() int                               { return 1 }
func (d ModalDelegate) Spacing() int                              { return 0 }
func (d ModalDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d ModalDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(ModalItem)
	if !ok {
		return
	}
	var line string
	if index == m.Index() {
		line = styles.SelectedStyle.Render("► " + item.Label)
	} else {
		line = "  " + item.Label
	}
	fmt.Fprint(w, line)
}

// NewFileList creates a new list model for files
func NewFileList(items []list.Item) list.Model {
	l := list.New(items, CustomDelegate{}, styles.FixedBoxWidth-4, styles.PanelHeight)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)
	return l
}

// NewModalList creates a new list model for modal options
func NewModalList(items []list.Item) list.Model {
	l := list.New(items, ModalDelegate{}, styles.ModalWidth-4, 8)
	l.SetShowStatusBar(false)
	l.SetShowFilter(false)
	l.SetShowHelp(false)
	l.SetShowPagination(false)
	l.SetShowTitle(false)
	return l
}
