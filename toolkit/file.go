package toolkit

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type FileChooser struct {
	root      ui.Component
	files     *List
	ab, cb    *Button
	nameField TextField
	path      string
	Action    func(*ui.State, string)
}

func NewFileChooser() *FileChooser {
	f := &FileChooser{}
	f.ab = NewButton("", f.action)
	f.cb = NewButton("", (*ui.State).CloseDialog)
	f.files = NewList()
	f.files.Changed = func(state *ui.State, i ListItem) { f.nameField.Text = (i.Text) }
	f.files.Action = func(state *ui.State, i ListItem) { f.action(state) }
	f.root = &Container{
		Center: NewScrollView(f.files),
		Bottom: NewBar(1, NewButtonIcon("left.arrow", "", f.back), &f.nameField, f.ab, f.cb),
	}
	f.nameField = *NewTextField()
	f.set(".")
	return f
}

func (f *FileChooser) SetTheme(theme *Theme) {
	SetTheme(f.root, theme)
}

func (f *FileChooser) SetLabels(action, cancel string, existing bool) {
	f.ab.Text = action
	f.cb.Text = cancel
	f.nameField.Editable = !existing
}

func (f *FileChooser) SetPath(path string) {
	f.set(path)
}

func (f *FileChooser) PreferredSize(fonts draw.FontLookup) (int, int) {
	return 400, 320
}

func (f *FileChooser) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	state.UpdateChild(g, draw.WH(w, h), f.root)
}

func (f *FileChooser) action(state *ui.State) {
	s := f.files.Items[f.files.Selected]
	if s.Text == f.nameField.Text {
		if s.Icon == "folder" {
			f.enter()
		} else {
			state.CloseDialog()
			if f.Action != nil {
				f.Action(state, filepath.Join(f.path, s.Text))
			}
		}
	} else if f.nameField.Editable && f.nameField.Text != "" {
		state.CloseDialog()
		if f.Action != nil {
			f.Action(state, filepath.Join(f.path, f.nameField.Text))
		}
	}
}

func (f *FileChooser) set(p string) {
	p, _ = filepath.Abs(p)
	dir, err := os.Open(p)
	if err != nil {
		return
	}
	files, err := dir.Readdir(0)
	if err != nil {
		return
	}
	f.files.Items = nil
	for _, i := range files {
		name := i.Name()
		i.Mode()
		if !strings.HasPrefix(name, ".") {
			icon := "file"
			if i.IsDir() {
				icon = "folder"
			}
			f.files.Items = append(f.files.Items, ListItem{Text: name, Icon: icon})
		}
	}
	sort.Slice(f.files.Items, func(i, j int) bool {
		if (f.files.Items[i].Icon == "folder") != (f.files.Items[j].Icon == "folder") {
			return f.files.Items[i].Icon == "folder"
		}
		return f.files.Items[i].Text < f.files.Items[j].Text
	})
	if len(f.files.Items) > 0 {
		f.nameField.Text = (f.files.Items[0].Text)
	}
	f.path = p
}

func (f *FileChooser) enter() {
	sel := f.files.Items[f.files.Selected].Text
	if strings.HasSuffix(sel, string(filepath.Separator)) {
		f.set(sel)
	} else {
		f.set(filepath.Join(f.path, sel))
	}
}
