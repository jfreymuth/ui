package main

import (
	"os"
	"path/filepath"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/impl/gofont"
	"github.com/jfreymuth/ui/impl/icons"
	"github.com/jfreymuth/ui/impl/sdl"
	. "github.com/jfreymuth/ui/toolkit"
)

type Editor struct {
	editor *TextArea
	files  *FileChooser

	unsavedChanges bool
	filePath       string
}

func main() {
	var e Editor
	e.editor = NewTextArea()
	e.editor.Font = draw.Font{Name: "gomono", Size: 11}
	e.files = NewFileChooser()
	e.files.SetPath(".")

	menuBar := NewMenuBar()
	fileMenu := menuBar.AddMenu("File")
	fileMenu.AddItemIcon("add", "New", e.New)
	fileMenu.AddItemIcon("open", "Open...", e.ShowOpenDialog)
	fileMenu.AddItemIcon("save", "Save", e.Save)
	fileMenu.AddItemIcon("save", "Save As...", e.ShowSaveDialog)
	fileMenu.AddItemIcon("close", "Exit", func(state *ui.State) { e.DoDestructive(state, (*ui.State).Quit) })
	editMenu := menuBar.AddMenu("Edit")
	editMenu.AddItemIcon("", "Select All", e.editor.SelectAll)
	editMenu.AddItemIcon("cut", "Cut", e.editor.Cut)
	editMenu.AddItemIcon("copy", "Copy", e.editor.Copy)
	editMenu.AddItemIcon("paste", "Paste", e.editor.Paste)

	root := NewRoot(&Container{
		Top:    menuBar,
		Center: NewScrollView(e.editor),
	})

	sdl.Show(sdl.Options{
		Title: "Editor",
		Width: 600, Height: 450,
		Root:       root,
		FontLookup: gofont.Lookup,
		IconLookup: &icons.Lookup{},
		Init: func(state *ui.State) {
			// Handle command line arguments in the Init function, so we can e.g. show an error dialog
			if len(os.Args) == 2 {
				e.files.SetPath(filepath.Dir(os.Args[1]))
				e.Open(state, os.Args[1])
				state.RequestUpdate()
			}
		},
		Update: func(state *ui.State) {
			// Handle custom keyboard shortcuts (Ctrl+O, Ctrl+S)
			if state.HasModifiers(ui.Control) {
				for _, k := range state.PeekKeyPresses() {
					switch k {
					case ui.KeyS:
						if state.HasModifiers(ui.Shift) {
							e.ShowSaveDialog(state)
						} else {
							e.Save(state)
						}
					case ui.KeyO:
						e.ShowOpenDialog(state)
					}
				}
			}
			// Handle default keyboard shortcuts (copy, paste etc.)
			ui.HandleKeyboardShortcuts(state)
		},
		Close: func(state *ui.State) {
			// If the Close function is set, it must call state.Quit() somewhere.
			// In this case, show a confirmation dialog.
			e.DoDestructive(state, (*ui.State).Quit)
		},
	})
}

func (e *Editor) ShowOpenDialog(state *ui.State) {
	e.DoDestructive(state, func(state *ui.State) {
		ShowOpenDialog(state, e.files, "Open", "Open", "Cancel", e.Open)
	})
}

func (e *Editor) ShowSaveDialog(state *ui.State) {
	ShowSaveDialog(state, e.files, "Save As", "Save", "Cancel", e.SaveAs)
}

func (e *Editor) New(state *ui.State) {
	e.DoDestructive(state, func(state *ui.State) {
		e.filePath = ""
		e.editor.SetText("")
		e.editor.Changed()
		e.unsavedChanges = false
	})
}

func (e *Editor) Open(state *ui.State, path string) {
	file, err := os.Open(path)
	if err != nil {
		e.Error(state, err.Error())
		return
	}
	err = e.editor.SetTextFromReader(file)
	file.Close()
	if err != nil {
		e.Error(state, err.Error())
		return
	}
	e.filePath = path
	e.editor.Changed()
	e.unsavedChanges = false
}

func (e *Editor) Save(state *ui.State) {
	if e.filePath == "" {
		e.ShowSaveDialog(state)
		return
	}
	e.SaveAs(state, e.filePath)
}

func (e *Editor) SaveAs(state *ui.State, path string) {
	file, err := os.Create(path)
	if err != nil {
		e.Error(state, err.Error())
		return
	}
	err = e.editor.WriteTextTo(file)
	file.Close()
	if err != nil {
		e.Error(state, err.Error())
		return
	}
	e.editor.Changed()
	e.unsavedChanges = false
	e.filePath = path
}

// If there are unsaved changes, asks the user for confirmation.
// Calls action only after a successful save or with the user's approval.
func (e *Editor) DoDestructive(state *ui.State, action func(*ui.State)) {
	if e.editor.Changed() {
		e.unsavedChanges = true
	}
	if e.unsavedChanges {
		ShowYesNoDialog(state, "Save Changes", "Save Changes?", "Save", "Discard", "Cancel", e.SaveAndThen(action), action)
	} else {
		action(state)
	}
}

// Returns a callback that will show a save dialog and only call action after a successful save.
func (e *Editor) SaveAndThen(action func(*ui.State)) func(*ui.State) {
	return func(state *ui.State) {
		if e.filePath == "" {
			ShowSaveDialog(state, e.files, "Save As", "Save", "Cancel", func(state *ui.State, path string) {
				e.SaveAs(state, path)
				if !e.unsavedChanges {
					action(state)
				}
			})
			return
		} else {
			e.SaveAs(state, e.filePath)
			if !e.unsavedChanges {
				action(state)
			}
		}
	}
}

func (e *Editor) Error(state *ui.State, msg string) {
	ShowErrorDialog(state, "Error", msg, "Close")
}
