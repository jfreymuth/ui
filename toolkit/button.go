package toolkit

import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/text"
)

type Button struct {
	Action func(*ui.State)
	Theme  *Theme
	Text   string
	Icon   string
	text   text.Text
	anim   float32
}

func NewButton(text string, action func(*ui.State)) *Button {
	return &Button{Action: action, Theme: DefaultTheme, Text: text}
}

func NewButtonIcon(icon, text string, action func(*ui.State)) *Button {
	return &Button{Action: action, Theme: DefaultTheme, Text: text, Icon: icon}
}

func (b *Button) SetTheme(theme *Theme) { b.Theme = theme }

func (b *Button) PreferredSize(fonts draw.FontLookup) (int, int) {
	w, h := b.text.SizeIcon(b.Text, b.Theme.Font("buttonText"), b.Icon, 5, fonts)
	return w + 20, h + 20
}

func (b *Button) Update(g *draw.Buffer, state *ui.State) {
	w, h := g.Size()
	animate(state, &b.anim, 8, state.IsHovered())
	g.Fill(draw.WH(w, h), draw.Blend(b.Theme.Color("buttonBackground"), b.Theme.Color("buttonHovered"), b.anim))
	color := b.Theme.Color("buttonText")
	if state.HasKeyboardFocus() {
		color = b.Theme.Color("buttonFocused")
	}
	b.text.DrawCenteredIcon(g, draw.WH(w, h), b.Text, b.Theme.Font("buttonText"), color, b.Icon, 5)
	/*if b.Icon != "" {
		if b.text == "" {
			g.Icon(draw.XYWH((w-b.w-b.h)/2-10, 0, b.h+20, h).Inset(5), b.Icon, textColor)
		} else {
			g.Icon(draw.XYWH((w-b.w-b.h)/2-15, 0, b.h+20, h).Inset(5), b.Icon, textColor)
		}
		x += 10 + b.h
	}
	g.Text(draw.XYXY(x, 0, w, h), b.text, textColor, b.theme.Font("buttonText"), draw.TextCentered)*/
	action := state.MouseClick(ui.MouseLeft)
	for _, k := range state.KeyPresses() {
		if k == ui.KeySpace || k == ui.KeyEnter {
			action = true
		}
	}
	if b.Action != nil && action {
		b.Action(state)
		state.RequestUpdate()
	}
}
