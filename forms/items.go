package forms

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// LoadFilesFromDir loads file names from a directory
func LoadFilesFromDir(dirPath string) ([]list.Item, error) {
	var items []list.Item

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			items = append(items, FileItem{Path: entry.Name(), Status: " "})
		}
	}

	return items, nil
}

// GetProjectFilesDir returns the path to the files directory in the project
func GetProjectFilesDir() string {
	// Get the directory where the executable is located
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to current working directory
		execPath, _ = os.Getwd()
	}
	execDir := filepath.Dir(execPath)

	// Check if files directory exists relative to executable
	filesDir := filepath.Join(execDir, "files")
	if _, err := os.Stat(filesDir); err == nil {
		return filesDir
	}

	// Check if files directory exists relative to current working directory
	cwd, _ := os.Getwd()
	filesDir = filepath.Join(cwd, "files")
	if _, err := os.Stat(filesDir); err == nil {
		return filesDir
	}

	// Default to relative path from where the program is run
	return "files"
}

// ReadFileContent reads and returns the content of a file
func ReadFileContent(dirPath, fileName string) string {
	fullPath := filepath.Join(dirPath, fileName)
	
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Sprintf("Ошибка чтения файла: %v", err)
	}
	
	return string(data)
}
