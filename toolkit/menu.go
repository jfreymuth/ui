package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type MenuItem struct {
	Theme  *Theme
	Text   string
	Icon   string
	text   text.Text
	Action func(*ui.State)
	parent menuParent
}

func (m *MenuItem) SetTheme(theme *Theme) { m.Theme = theme }

func (m *MenuItem) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := m.text.SizeIcon(m.Text, m.Theme.Font("text"), m.Icon, 4, fonts)
	return w + 10, h + 6
}

func (m *MenuItem) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	if state.IsHovered() {
		state.SetKeyboardFocus(m)
		drag, drop := state.DraggedContent()
		if drag == nil && state.MouseButtonDown(ui.MouseLeft) {
			state.InitiateDrag(ui.MenuDrag)
			drag = ui.MenuDrag
		}
		if drop && drag == ui.MenuDrag {
			state.ClosePopups()
			if m.Action != nil {
				m.Action(state)
			}
		}
	}
	if state.HasKeyboardFocus() {
		m.parent.setOpen(nil)
		g.Fill(draw.WH(w, h), m.Theme.Color("selection"))
		for _, k := range state.KeyPresses() {
			switch k {
			case ui.KeyDown:
				state.FocusNext()
			case ui.KeyUp:
				state.FocusPrevious()
			case ui.KeyLeft:
				state.SetKeyboardFocus(m.parent)
			case ui.KeySpace, ui.KeyEnter:
				state.ClosePopups()
				if m.Action != nil {
					m.Action(state)
				}
			}
		}
	}
	m.text.DrawLeftIcon(g, draw.XYXY(5, 0, w-5, h), m.Text, m.Theme.Font("text"), m.Theme.Color("text"), m.Icon, 4)
}

type Menu struct {
	Theme  *Theme
	Text   string
	text   text.Text
	parent menuParent
	items  []ui.Component
	open   *Menu
	popup  ui.Popup
}

func NewPopupMenu(text string) *Menu {
	return &Menu{Theme: DefaultTheme, Text: text}
}

func (m *Menu) SetTheme(theme *Theme) {
	m.Theme = theme
	for _, i := range m.items {
		SetTheme(i, theme)
	}
}

func (m *Menu) AddItem(text string, action func(*ui.State)) *MenuItem {
	i := &MenuItem{parent: m, Action: action, Theme: m.Theme, Text: text}
	m.items = append(m.items, i)
	return i
}

func (m *Menu) AddItemIcon(icon, text string, action func(*ui.State)) *MenuItem {
	i := &MenuItem{parent: m, Action: action, Theme: m.Theme, Text: text, Icon: icon}
	m.items = append(m.items, i)
	return i
}

func (m *Menu) AddMenu(text string) *Menu {
	menu := &Menu{parent: m, Theme: m.Theme, Text: text}
	m.items = append(m.items, menu)
	return menu
}

func (m *Menu) OpenPopupMenu(p image.Point, state *ui.State, fonts draw.FontLookup) {
	c := &menuBackground{&Stack{m.items}, m.Theme}
	w, h := c.PreferredSize(fonts)
	win := state.WindowBounds()
	if r := draw.XYWH(p.X, p.Y, w, h); r.In(win) {
		state.OpenPopup(r, c)
	} else if r = draw.XYWH(p.X, p.Y, w, -h); r.In(win) {
		state.OpenPopup(r, c)
	} else if r = draw.XYWH(p.X, p.Y, -w, h); r.In(win) {
		state.OpenPopup(r, c)
	} else {
		state.OpenPopup(draw.XYWH(p.X, p.Y, -w, -h), c)
	}
}

func (m *Menu) setOpen(p *Menu) {
	if m.open != nil {
		if m.open.popup != nil {
			m.open.popup.Close()
		}
		m.open.setOpen(nil)
	}
	m.open = p
}

func (m *Menu) isMenuBar() bool { return false }

func (m *Menu) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := m.text.Size(m.Text, m.Theme.Font("text"), fonts)
	if m.parent.isMenuBar() {
		return w + 10, h + 6
	} else {
		return w + h + 10, h + 6
	}
}

func (m *Menu) Update(g *draw.Buffer, state *ui.State) {
	submenu := !m.parent.isMenuBar()
	w, h := g.Size()
	if submenu {
		for _, k := range state.KeyPresses() {
			switch k {
			case ui.KeyDown:
				state.FocusNext()
			case ui.KeyUp:
				state.FocusPrevious()
			case ui.KeyRight:
				state.SetKeyboardFocus(m.items[0])
			case ui.KeyLeft:
				state.SetKeyboardFocus(m.parent)
			}
		}
	}
	open := state.HasPopups() && state.IsHovered()
	if state.IsHovered() {
		g.Fill(draw.WH(w, h), m.Theme.Color("buttonHovered"))
		if state.MouseButtonDown(ui.MouseLeft) {
			state.InitiateDrag(ui.MenuDrag)
			open = true
		}
	}
	if submenu && state.HasKeyboardFocus() {
		open = true
	}
	if open {
		if !isOpen(m.popup) {
			popup := &menuBackground{&Stack{m.items}, m.Theme}
			mw, mh := popup.PreferredSize(g.FontLookup)
			m.parent.setOpen(m)
			if submenu {
				m.popup = state.OpenPopup(draw.XYWH(w-3, -1, mw, mh), popup)
			} else {
				m.popup = state.OpenPopup(draw.XYWH(-3, h-1, mw, mh), popup)
			}
			state.SetKeyboardFocus(m)
		}
		if m.open != nil {
			m.setOpen(nil)
			state.SetKeyboardFocus(m)
		}
	}
	if submenu && isOpen(m.popup) {
		g.Fill(draw.WH(w, h), m.Theme.Color("selection"))
	}
	m.text.DrawLeft(g, draw.XYXY(5, 0, w-5, h), m.Text, m.Theme.Font("text"), m.Theme.Color("text"))
	if submenu {
		_, th := m.text.Size(m.Text, m.Theme.Font("text"), g.FontLookup)
		g.Icon(draw.XYXY(w-th-5, 0, w-5, h), "right", m.Theme.Color("text"))
	}
}

type MenuBar struct {
	Theme *Theme
	menus []*Menu
	open  *Menu
	popup ui.Popup
}

func NewMenuBar() *MenuBar {
	return &MenuBar{Theme: DefaultTheme}
}

func (m *MenuBar) SetTheme(theme *Theme) {
	m.Theme = theme
	for _, m := range m.menus {
		m.SetTheme(theme)
	}
}

func (m *MenuBar) AddMenu(name string) *Menu {
	menu := &Menu{parent: m, Theme: m.Theme, Text: name}
	m.menus = append(m.menus, menu)
	return menu
}

func (m *MenuBar) setOpen(p *Menu) {
	if m.open != nil {
		if m.open.popup != nil {
			m.open.popup.Close()
		}
		m.open.setOpen(nil)
	}
	m.open = p
}

func (m *MenuBar) isMenuBar() bool { return true }

func (m *MenuBar) PreferredSize(fonts draw.FontLookup) (int, int) {
	w := 0
	h := 0
	for _, m := range m.menus {
		mw, mh := m.PreferredSize(fonts)
		w += mw
		if mh > h {
			h = mh
		}

	}
	return w, h
}

func (m *MenuBar) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	g.Fill(draw.WH(w, h), m.Theme.Color("altBackground"))
	x := 0
	for _, menu := range m.menus {
		mw, _ := menu.PreferredSize(g.FontLookup)
		state.UpdateChild(g, draw.XYWH(x, 0, mw, h), menu)
		x += mw
	}
	if m.open != nil && isOpen(m.open.popup) && !isOpen(m.popup) {
		m.popup = state.OpenPopup(draw.WH(w, h), m)
	}
	if state.MouseClick(ui.MouseLeft) {
		state.ClosePopups()
	}
}

type menuParent interface {
	ui.Component
	setOpen(*Menu)
	isMenuBar() bool
}

type menuBackground struct {
	Content ui.Component
	theme   *Theme
}

func (m *menuBackground) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := m.Content.PreferredSize(fonts)
	return w + 6, h + 6
}

func (m *menuBackground) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	g.Shadow(draw.XYXY(3, 3, w-3, h-3), m.theme.Color("shadow"), 4)
	g.Fill(draw.XYXY(3, 1, w-3, h-5), m.theme.Color("background"))
	state.UpdateChild(g, draw.XYXY(3, 1, w-3, h-5), m.Content)
}

func isOpen(p ui.Popup) bool {
	return p != nil && !p.Closed()
}
