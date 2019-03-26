package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type CheckBox struct {
	Checked bool
	Changed func(*ui.State, bool)
	Theme   *Theme
	Text    string
	text    text.Text
	anim    float32
}

func NewCheckBox(text string) *CheckBox {
	//return &CheckBox{Label: Label{Color: DefaultTheme.Color("buttonText"), font: DefaultTheme.Font("buttonText"), text: text, h: -1}}
	return &CheckBox{Theme: DefaultTheme, Text: text}
}

func (c *CheckBox) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := c.text.Size(c.Text, DefaultTheme.Font("buttonText"), fonts)
	return w + h + 30, h + 20
}

func (c *CheckBox) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	_, s := c.text.Size(c.Text, DefaultTheme.Font("buttonText"), g.FontLookup)

	action := state.MouseClick(ui.MouseLeft)
	for _, k := range state.KeyPresses() {
		if k == ui.KeySpace || k == ui.KeyEnter {
			action = true
		}
	}
	if action {
		c.Checked = !c.Checked
		if c.Changed != nil {
			c.Changed(state, c.Checked)
			state.RequestUpdate()
		}
	}

	animate(state, &c.anim, 8, c.Checked)
	x := int(c.anim * float32(s+10))
	g.Add(draw.Command{Bounds: draw.XYXY(5, 5, s+15, h-5), Clip: draw.XYXY(5, 5, 5+x, h-5), Color: DefaultTheme.Color("buttonText"), Text: "checkboxChecked", Style: draw.Icon})
	g.Add(draw.Command{Bounds: draw.XYXY(5, 5, s+15, h-5), Clip: draw.XYXY(5+x, 5, s+15, h-5), Color: DefaultTheme.Color("buttonText"), Text: "checkbox", Style: draw.Icon})
	color := DefaultTheme.Color("buttonText")
	if state.HasKeyboardFocus() {
		color = DefaultTheme.Color("buttonFocused")
	}
	c.text.DrawLeft(g, draw.XYXY(s+20, 0, w, h), c.Text, DefaultTheme.Font("buttonText"), color)
}
