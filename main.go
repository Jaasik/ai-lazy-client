package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
)

type App struct {
	files        []string
	selectedFile int
	fileContent  string
	inputText    string
	cursorX      int
	currentView  string // "files", "content", "input"
	gui          *gocui.Gui
}

func NewApp() *App {
	return &App{
		files:       loadFiles(),
		currentView: "files",
	}
}

func loadFiles() []string {
	var files []string
	filepath.Walk("./files", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	return files
}

func (app *App) loadFileContent(name string) {
	content, err := ioutil.ReadFile(filepath.Join("./files", name))
	if err != nil {
		app.fileContent = "Error reading file: " + err.Error()
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
		fmt.Fprintln(v, "Welcome to LazyGit Clone - File Manager")
	}

	// Files panel (left top)
	if v, err := g.SetView("files", 0, 3, maxX/2-1, maxY/2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Files (Enter to select) "
		v.Frame = true
		v.Editable = false
		v.Wrap = false
		app.renderFiles(v)
	}

	// Status panel (left middle)
	if v, err := g.SetView("status", 0, maxY/2, maxX/2-1, maxY*2/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Status "
		v.Frame = true
		fmt.Fprintln(v, "Ready")
		fmt.Fprintf(v, "Selected: %s\n", app.getSelectedFile())
	}

	// Help panel (left bottom)
	if v, err := g.SetView("help", 0, maxY*2/3, maxX/2-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Help "
		v.Frame = true
		fmt.Fprintln(v, "Up/Down: Navigate files")
		fmt.Fprintln(v, "Enter: Select file")
		fmt.Fprintln(v, "Tab: Switch panels")
		fmt.Fprintln(v, "q: Quit")
	}

	// Content panel (right top)
	if v, err := g.SetView("content", maxX/2, 3, maxX-1, maxY/2+maxY/3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " File Content "
		v.Frame = true
		v.Editable = false
		v.Wrap = true
		app.renderContent(v)
	}

	// Input panel (right bottom)
	if v, err := g.SetView("input", maxX/2, maxY/2+maxY/3, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = " Input "
		v.Frame = true
		v.Editable = true
		v.Wrap = true
		// Set cursor position based on input text length
		app.cursorX = len(app.inputText)
		v.Write([]byte(app.inputText))
	}

	// Set active view
	if err := app.setActiveView(g); err != nil {
		return err
	}

	return nil
}

func (app *App) renderFiles(v *gocui.View) {
	v.Clear()
	for i, file := range app.files {
		if i == app.selectedFile {
			fmt.Fprintln(v, "> "+file)
		} else {
			fmt.Fprintln(v, "  "+file)
		}
	}
}

func (app *App) renderContent(v *gocui.View) {
	v.Clear()
	lines := strings.Split(app.fileContent, "\n")
	for i, line := range lines {
		colorCode := i%5 + 31 // Cycle through colors 31-35
		fmt.Fprintf(v, "\x1b[%dm%s\x1b[0m\n", colorCode, line)
	}
}

func (app *App) setActiveView(g *gocui.Gui) error {
	// Set active view - current view will have focus
	g.SetCurrentView(app.currentView)
	return nil
}

func (app *App) getSelectedFile() string {
	if len(app.files) > 0 && app.selectedFile >= 0 && app.selectedFile < len(app.files) {
		return app.files[app.selectedFile]
	}
	return "none"
}

func nextView(current string) string {
	switch current {
	case "files":
		return "content"
	case "content":
		return "input"
	case "input":
		return "files"
	default:
		return "files"
	}
}

func (app *App) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (app *App) navigateFiles(g *gocui.Gui, v *gocui.View, dir int) error {
	newIdx := app.selectedFile + dir
	if newIdx >= 0 && newIdx < len(app.files) {
		app.selectedFile = newIdx
	}
	return nil
}

func (app *App) selectFile(g *gocui.Gui, v *gocui.View) error {
	if len(app.files) > 0 {
		app.loadFileContent(app.files[app.selectedFile])
		if contentV, err := g.View("content"); err == nil {
			contentV.Clear()
			app.renderContent(contentV)
		}
	}
	return nil
}

func (app *App) switchPanel(g *gocui.Gui, v *gocui.View) error {
	app.currentView = nextView(app.currentView)
	return nil
}

func (app *App) onMouseClick(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		app.currentView = v.Name()
	}
	return nil
}

func main() {
	app := NewApp()
	if len(app.files) > 0 {
		app.loadFileContent(app.files[0])
	}

	g := gocui.NewGui()
	defer g.Close()

	app.gui = g
	
	// Enable mouse support
	g.Mouse = true

	g.SetLayout(app.layout)

	// Key bindings
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, app.quit)
	g.SetKeybinding("", 'q', gocui.ModNone, app.quit)

	// Navigation in files
	g.SetKeybinding("files", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return app.navigateFiles(g, v, -1)
	})
	g.SetKeybinding("files", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return app.navigateFiles(g, v, 1)
	})
	g.SetKeybinding("files", gocui.KeyEnter, gocui.ModNone, app.selectFile)

	// Switch between panels
	g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, app.switchPanel)

	// Mouse support
	g.SetKeybinding("", gocui.MouseLeft, gocui.ModNone, app.onMouseClick)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}
