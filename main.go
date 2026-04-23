package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/jesseduffield/gocui"
)

var (
	filesList    []string
	selectedFile int
	fileContent  string
	inputText    string
	currentView  string // "files", "status", "help", "content", "input"
	filesDir     = "./files"
)

func main() {
	loadFiles()

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.BgColor = gocui.ColorDefault
	g.FgColor = gocui.ColorWhite
	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack
	g.Cursor = true
	g.Mouse = true

	g.SetLayout(layout)
	setupKeybindings(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func loadFiles() {
	filesList = []string{}
	entries, err := ioutil.ReadDir(filesDir)
	if err != nil {
		filesList = append(filesList, "Error: "+err.Error())
		return
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			filesList = append(filesList, entry.Name())
		}
	}
	if len(filesList) > 0 {
		loadFileContent(filesList[0])
	}
}

func loadFileContent(filename string) {
	path := filepath.Join(filesDir, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fileContent = "Error reading file: " + err.Error()
		return
	}
	fileContent = string(data)
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Clear all views
	g.DeleteView("files")
	g.DeleteView("content")
	g.DeleteView("input")
	g.DeleteView("status")
	g.DeleteView("help")

	// Layout similar to lazygit:
	// Left column (1/3 width): Files (top), Status (middle), Help (bottom)
	// Right column (2/3 width): Content (top), Input (bottom)

	leftWidth := maxX / 3

	// Files panel (left top)
	if v, err := g.SetView("files", 0, 0, leftWidth-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Files"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Frame = true
		renderFiles(v)
	}

	// Status panel (left middle)
	if v, err := g.SetView("status", 0, maxY/2, leftWidth-1, maxY*2/3-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Status"
		v.Frame = true
		fmt.Fprintf(v, "Ready\nSelected: %s", getCurrentFile())
	}

	// Help panel (left bottom)
	if v, err := g.SetView("help", 0, maxY*2/3, leftWidth-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Help"
		v.Frame = true
		fmt.Fprintln(v, "Up/Down: Navigate")
		fmt.Fprintln(v, "Enter: Select")
		fmt.Fprintln(v, "Tab: Switch panel")
		fmt.Fprintln(v, "q: Quit")
	}

	// Content panel (right top)
	if v, err := g.SetView("content", leftWidth, 0, maxX-1, maxY*2/3-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "File Content"
		v.Wrap = true
		v.Frame = true
		renderContent(v)
	}

	// Input panel (right bottom)
	if v, err := g.SetView("input", leftWidth, maxY*2/3, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Frame = true
		v.Write([]byte(inputText))
	}

	// Set focus
	if currentView == "" {
		currentView = "files"
	}
	g.SetCurrentView(currentView)

	return nil
}

func renderFiles(v *gocui.View) {
	v.Clear()
	for i, file := range filesList {
		if strings.HasPrefix(file, "Error:") {
			fmt.Fprintln(v, file)
			continue
		}
		if i == selectedFile {
			fmt.Fprintln(v, "> "+file)
		} else {
			fmt.Fprintln(v, "  "+file)
		}
	}
}

func renderContent(v *gocui.View) {
	v.Clear()
	lines := strings.Split(fileContent, "\n")
	for _, line := range lines {
		fmt.Fprintln(v, line)
	}
}

func getCurrentFile() string {
	if len(filesList) == 0 || strings.HasPrefix(filesList[selectedFile], "Error:") {
		return "None"
	}
	return filesList[selectedFile]
}

func setupKeybindings(g *gocui.Gui) error {
	// Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	// Navigation in files
	if err := g.SetKeybinding("files", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}
	if err := g.SetKeybinding("files", gocui.KeyEnter, gocui.ModNone, selectFile); err != nil {
		return err
	}

	// Tab to switch views
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextView); err != nil {
		return err
	}

	// Input handling
	if err := g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, inputEnter); err != nil {
		return err
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if selectedFile > 0 {
		selectedFile--
		loadFileContent(filesList[selectedFile])
		g.Execute(func(g *gocui.Gui) error {
			return layout(g)
		})
	}
	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if selectedFile < len(filesList)-1 && !strings.HasPrefix(filesList[selectedFile], "Error:") {
		selectedFile++
		loadFileContent(filesList[selectedFile])
		g.Execute(func(g *gocui.Gui) error {
			return layout(g)
		})
	}
	return nil
}

func selectFile(g *gocui.Gui, v *gocui.View) error {
	loadFileContent(filesList[selectedFile])
	return nil
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	views := []string{"files", "status", "help", "content", "input"}
	idx := -1
	for i, view := range views {
		if view == currentView {
			idx = i
			break
		}
	}
	idx = (idx + 1) % len(views)
	currentView = views[idx]
	g.SetCurrentView(currentView)
	return nil
}

func inputEnter(g *gocui.Gui, v *gocui.View) error {
	inputText = strings.TrimSpace(v.Buffer())
	v.Clear()
	return nil
}
