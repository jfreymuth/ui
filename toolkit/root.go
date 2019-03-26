package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type Root struct {
	Content ui.Component
	Dialog  ui.Component
	Theme   *Theme
	popups  []*popup
}

type popup struct {
	ui.Component
	bounds image.Rectangle
	closed bool
}

func NewRoot(content ui.Component) *Root {
	return &Root{Content: content, Theme: DefaultTheme}
}

func (r *Root) SetTheme(theme *Theme) {
	r.Theme = theme
	SetTheme(r.Content, theme)
	SetTheme(r.Dialog, theme)
}

func (r *Root) OpenDialog(dialog ui.Component) {
	r.Dialog = dialog
}

func (r *Root) CloseDialog() {
	r.Dialog = nil
}

func (r *Root) OpenPopup(bounds image.Rectangle, p ui.Component) ui.Popup {
	popup := &popup{p, bounds, false}
	r.popups = append(r.popups, popup)
	return popup
}

func (r *Root) ClosePopups() {
	for _, p := range r.popups {
		p.closed = true
	}
	r.popups = nil
}

func (r *Root) HasPopups() bool {
	return r.popups != nil
}

func (p *popup) Close()       { p.closed = true }
func (p *popup) Closed() bool { return p.closed }

func (r *Root) PreferredSize(fonts draw.FontLookup) (int, int) {
	return r.Content.PreferredSize(fonts)
}

func (r *Root) Update(g *draw.Buffer, state *ui.State) {
	state.SetRoot(r)
	w, h := g.Size()
	g.Fill(draw.WH(w, h), r.Theme.Color("background"))
	if r.Dialog == nil {
		if len(r.popups) == 0 {
			state.UpdateChild(g, draw.WH(w, h), r.Content)
		} else {
			state.DrawChild(g, draw.WH(w, h), r.Content)
		}
	} else {
		state.DrawChild(g, draw.WH(w, h), r.Content)
		g.Fill(draw.WH(w, h), r.Theme.Color("veil"))
		dw, dh := r.Dialog.PreferredSize(g.FontLookup)
		if dw > w*7/8 {
			dw = w * 7 / 8
		}
		if dh > h*7/8 {
			dh = h * 7 / 8
		}
		dx, dy := (w-dw)/2, (h-dh)/2
		if len(r.popups) == 0 {
			state.UpdateChild(g, draw.XYWH(dx, dy, dw, dh), r.Dialog)
		} else {
			state.DrawChild(g, draw.XYWH(dx, dy, dw, dh), r.Dialog)
		}
	}
	if r.popups != nil {
		allClosed := true
		for _, p := range r.popups {
			if !p.closed {
				state.UpdateChild(g, p.bounds, p.Component)
				allClosed = false
			}
		}
		if allClosed {
			r.popups = nil
		}
		if state.MouseButtonDown(ui.MouseLeft) || state.MouseButtonDown(ui.MouseRight) {
			state.ClosePopups()
			state.RequestRefocus()
		}
	}
	if drag, drop := state.DraggedContent(); drag != nil && !drop {
		mouse := state.MousePos()
		switch drag := drag.(type) {
		case string:
			h := g.FontLookup.Metrics(r.Theme.Font("text")).LineHeight()
			g.Text(draw.XYWH(mouse.X, mouse.Y, w, h), drag, r.Theme.Color("text"), r.Theme.Font("text"))
		}
	}
}
