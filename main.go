package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
)

type App struct {
	gui           *gocui.Gui
	files         []string
	selectedFile  int
	inputText     string
	cursorX       int
	currentView   string // "files", "status", "help", "content", "input"
	fileContent   string
}

func NewApp() (*App, error) {
	app := &App{
		currentView: "files",
	}
	
	// Load files from ./files directory
	filesDir := "./files"
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create files directory: %v", err)
	}
	
	entries, err := ioutil.ReadDir(filesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read files directory: %v", err)
	}
	
	for _, entry := range entries {
		if !entry.IsDir() {
			app.files = append(app.files, entry.Name())
		}
	}
	
	if len(app.files) > 0 {
		app.loadFileContent(0)
	}
	
	return app, nil
}

func (app *App) loadFileContent(index int) {
	if index < 0 || index >= len(app.files) {
		return
	}
	
	content, err := ioutil.ReadFile(filepath.Join("./files", app.files[index]))
	if err != nil {
		app.fileContent = fmt.Sprintf("Error reading file: %v", err)
		return
	}
	app.fileContent = string(content)
}

func (app *App) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	
	// Service panel (top)
	if v, err := g.SetView("service", 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " LazyGit Clone Service "
		v.Frame = true
		fmt.Fprintln(v, "Welcome to LazyGit Clone - Press 'q' to quit")
	}
	
	// Left column panels
	leftWidth := maxX / 2
	
	// Files panel
	if v, err := g.SetView("files", 0, 3, leftWidth-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Files (Enter to select) "
		v.Frame = true
		v.Editable = false
		v.Wrap = false
		app.renderFiles(v)
	}
	
	// Status panel
	if v, err := g.SetView("status", 0, maxY/2, leftWidth-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Status "
		v.Frame = true
		fmt.Fprintf(v, "Selected: %s\nTotal files: %d\nNavigation: ↑/↓ arrows\nSelect: Enter\nSwitch panel: Tab", 
			func() string {
				if len(app.files) > 0 && app.selectedFile < len(app.files) {
					return app.files[app.selectedFile]
				}
				return "none"
			}(),
			len(app.files))
	}
	
	// Help panel
	if v, err := g.SetView("help", 0, maxY-2, leftWidth-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Help "
		v.Frame = true
		fmt.Fprintln(v, "↑/↓: Navigate files")
		fmt.Fprintln(v, "Enter: Select file")
		fmt.Fprintln(v, "Tab: Switch panels")
		fmt.Fprintln(v, "q: Quit")
	}
	
	// Right column panels
	rightStart := leftWidth
	
	// File Content panel
	if v, err := g.SetView("content", rightStart, 3, maxX-1, maxY/2+maxY/4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " File Content "
		v.Frame = true
		v.Wrap = true
		app.renderContent(v)
	}
	
	// Input panel
	if v, err := g.SetView("input", rightStart, maxY/2+maxY/4, maxX-1, maxY); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Input (type here) "
		v.Frame = true
		v.Editable = true
		v.Wrap = true
		if v.Buffer() == "" {
			v.Write([]byte(app.inputText))
		} else {
			app.inputText = v.Buffer()
		}
	}
	
	return nil
}

func (app *App) renderFiles(v *gocui.View) {
	v.Clear()
	for i, file := range app.files {
		if i == app.selectedFile {
			fmt.Fprintf(v, "> %s\n", file)
		} else {
			fmt.Fprintf(v, "  %s\n", file)
		}
	}
}

func (app *App) renderContent(v *gocui.View) {
	v.Clear()
	lines := strings.Split(app.fileContent, "\n")
	colors := []gocui.Attribute{gocui.ColorWhite, gocui.ColorCyan, gocui.ColorGreen, gocui.ColorYellow, gocui.ColorMagenta}
	
	for i, line := range lines {
		color := colors[i%len(colors)]
		fmt.Fprintf(v, "%s%s\n", color, line)
	}
}

func (app *App) nextFile() {
	if len(app.files) == 0 {
		return
	}
	app.selectedFile = (app.selectedFile + 1) % len(app.files)
	app.loadFileContent(app.selectedFile)
}

func (app *App) prevFile() {
	if len(app.files) == 0 {
		return
	}
	app.selectedFile = (app.selectedFile - 1 + len(app.files)) % len(app.files)
	app.loadFileContent(app.selectedFile)
}

func (app *App) selectFile() {
	// File is already selected, just update status
}

func (app *App) switchPanel() {
	panels := []string{"files", "status", "help", "content", "input"}
	currentIdx := 0
	for i, p := range panels {
		if p == app.currentView {
			currentIdx = i
			break
		}
	}
	nextIdx := (currentIdx + 1) % len(panels)
	app.currentView = panels[nextIdx]
}

func (app *App) keyBindings(g *gocui.Gui) error {
	// Global keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, app.handleTab); err != nil {
		return err
	}
	
	// Files panel keybindings
	if err := g.SetKeybinding("files", gocui.KeyArrowUp, gocui.ModNone, app.handleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", gocui.KeyArrowDown, gocui.ModNone, app.handleDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", gocui.KeyEnter, gocui.ModNone, app.handleEnter); err != nil {
		return err
	}
	
	// Content panel keybindings
	if err := g.SetKeybinding("content", gocui.KeyArrowUp, gocui.ModNone, app.handleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("content", gocui.KeyArrowDown, gocui.ModNone, app.handleDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("content", gocui.KeyEnter, gocui.ModNone, app.handleEnter); err != nil {
		return err
	}
	
	// Status panel keybindings
	if err := g.SetKeybinding("status", gocui.KeyArrowUp, gocui.ModNone, app.handleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("status", gocui.KeyArrowDown, gocui.ModNone, app.handleDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("status", gocui.KeyEnter, gocui.ModNone, app.handleEnter); err != nil {
		return err
	}
	
	// Help panel keybindings
	if err := g.SetKeybinding("help", gocui.KeyArrowUp, gocui.ModNone, app.handleUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyArrowDown, gocui.ModNone, app.handleDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("help", gocui.KeyEnter, gocui.ModNone, app.handleEnter); err != nil {
		return err
	}
	
	return nil
}

func (app *App) handleUp(g *gocui.Gui, v *gocui.View) error {
	if app.currentView == "files" || app.currentView == "content" || app.currentView == "status" || app.currentView == "help" {
		app.prevFile()
	}
	return nil
}

func (app *App) handleDown(g *gocui.Gui, v *gocui.View) error {
	if app.currentView == "files" || app.currentView == "content" || app.currentView == "status" || app.currentView == "help" {
		app.nextFile()
	}
	return nil
}

func (app *App) handleEnter(g *gocui.Gui, v *gocui.View) error {
	if app.currentView == "files" || app.currentView == "content" || app.currentView == "status" || app.currentView == "help" {
		app.selectFile()
	}
	return nil
}

func (app *App) handleTab(g *gocui.Gui, v *gocui.View) error {
	app.switchPanel()
	g.SetCurrentView(app.currentView)
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to create app: %v", err)
	}
	
	gui := gocui.NewGui()
	if gui == nil {
		log.Fatalf("Failed to create GUI")
	}
	defer gui.Close()
	
	app.gui = gui
	
	gui.SetLayout(app.layout)
	
	if err := app.keyBindings(gui); err != nil {
		log.Fatalf("Failed to set keybindings: %v", err)
	}
	
	// Set initial focus
	gui.SetCurrentView("files")
	
	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("GUI error: %v", err)
	}
}
