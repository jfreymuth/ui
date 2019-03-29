package toolkit

import (
	"image"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
)

type ComboBox struct {
	List
	sv   *ScrollView
	anim float32
}

func NewComboBox() *ComboBox {
	c := &ComboBox{List: List{Theme: DefaultTheme}}
	c.sv = NewScrollView(&c.List)
	return c
}

func (c *ComboBox) SetTheme(theme *Theme) { c.sv.SetTheme(theme) }

func (c *ComboBox) PreferredSize(fonts draw.FontLookup) (int, int) {
	if len(c.Items) == 0 {
		return 30, 16
	}
	w, h := c.List.PreferredSize(fonts)
	h /= len(c.Items)
	return w + 30 + h, h + 16
}

func (c *ComboBox) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	var item ListItem
	if len(c.Items) > 0 {
		if c.Selected < 0 {
			c.Selected = 0
		} else if c.Selected >= len(c.Items) {
			c.Selected = len(c.Items) - 1
		}
		item = c.Items[c.Selected]
	}
	animate(state, &c.anim, 8, state.IsHovered())
	g.Fill(draw.WH(w, h), draw.Blend(c.Theme.Color("buttonBackground"), c.Theme.Color("buttonHovered"), c.anim))
	_, th := item.text.Size(item.Text, c.Theme.Font("text"), g.FontLookup)
	item.text.DrawLeftIcon(g, draw.XYXY(15, 0, w-th-10, h), item.Text, c.Theme.Font("text"), c.Theme.Color("text"), item.Icon, 3)
	g.Icon(draw.XYXY(w-th-10, 0, w-10, h), "down", c.Theme.Color("text"))
	if state.MouseClick(ui.MouseLeft) {
		pw, ph := c.List.PreferredSize(g.FontLookup)
		pw, ph = pw+6, ph+6
		win := state.WindowBounds()
		spaceAbove := -win.Min.Y
		spaceBelow := win.Max.Y - h
		var r image.Rectangle
		if ph <= spaceBelow {
			r = draw.XYWH(0, h, pw, ph)
		} else if ph <= spaceAbove {
			r = draw.XYWH(0, -ph, pw, ph)
		} else if spaceAbove > spaceBelow {
			r = draw.XYWH(0, -spaceAbove, pw+15, spaceAbove)
		} else {
			r = draw.XYWH(0, h, pw+15, spaceBelow)
		}
		if r.Dx() < w {
			r.Max.X = r.Min.X + w
		}
		state.OpenPopup(r, &menuBackground{c.sv, c.Theme})
	}
}
