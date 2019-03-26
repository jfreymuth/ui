package main

import (
	"fmt"
	"strconv"

	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/draw"
	"github.com/jfreymuth/ui/impl/gofont"
	"github.com/jfreymuth/ui/impl/icons"
	"github.com/jfreymuth/ui/impl/sdl"
	. "github.com/jfreymuth/ui/toolkit"
)

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {

	themeButton := NewButtonIcon("refresh", "Light Theme", nil)

	mb := NewMenuBar()
	{
		menu := mb.AddMenu("Test")
		menu.AddItem("Test", nil)
		sm := menu.AddMenu("Menu")
		sm.AddItem("Submenu1", nil)
		sm.AddItem("Submenu2", nil)
		sm = sm.AddMenu("Submenu3")
		sm.AddItem("Submenu1", nil)
		sm.AddItem("Submenu2", nil)
		sm = sm.AddMenu("Submenu3")
		sm.AddItem("Submenu1", nil)
		sm.AddItem("Submenu2", nil)
		menu.AddItem("Items", nil)
		menu.AddItem("Quit", (*ui.State).Quit)
	}
	{
		menu := mb.AddMenu("Dialog")
		menu.AddItem("Message", func(state *ui.State) {
			ShowMessageDialog(state, "Message", "This is a message.", "Ok")
		})
		menu.AddItem("Error", func(state *ui.State) {
			ShowErrorDialog(state, "Error", "This is an error.\nSeriously.", "Close")
		})
		menu.AddItem("Confirm", func(state *ui.State) {
			ShowConfirmDialog(state, "Confirm", "Quit?", "Quit", "Cancel", func(state *ui.State) {
				ShowMessageDialog(state, "Message", "Quit.", "Ok")
			})
		})
		menu.AddItem("Yes/No", func(state *ui.State) {
			ShowYesNoDialog(state, "Confirm", "Save before quitting?", "Save", "Don't Save", "Cancel", func(state *ui.State) {
				ShowMessageDialog(state, "Message", "Save and Quit", "Ok")
			}, func(state *ui.State) {
				ShowMessageDialog(state, "Message", "Quit without saving", "Ok")
			})
		})
		menu.AddItem("Input", func(state *ui.State) {
			ShowInputDialog(state, "Input", "Input something", "Ok", "Cancel", func(state *ui.State, text string) {
				ShowMessageDialog(state, "Message", fmt.Sprint("Your input: ", text), "Ok")
			})
		})
		menu.AddItem("Open", func(state *ui.State) {
			ShowOpenDialog(state, NewFileChooser(), "Open", "Open", "Cancel", func(state *ui.State, path string) {
				ShowMessageDialog(state, "Message", fmt.Sprint("Your input: ", path), "Ok")
			})
		})
		menu.AddItem("Save", func(state *ui.State) {
			ShowSaveDialog(state, NewFileChooser(), "Save As", "Save", "Cancel", func(state *ui.State, path string) {
				ShowMessageDialog(state, "Message", fmt.Sprint("Your input: ", path), "Ok")
			})
		})
	}

	ta := NewTextArea()
	ta.Font = draw.Font{Name: "mono", Size: 11}
	ta.SetText(ipsum)
	font := NewTextField()
	font.Text = "mono"
	size := NewTextField()
	size.Text = "11"
	size.MinWidth = 30
	update := func(state *ui.State, _ string) {
		sz, err := strconv.ParseFloat(size.Text, 32)
		if err != nil {
			ShowErrorDialog(state, "Error", fmt.Sprint("\"", size.Text, "\" is not a number"), "Close")
			return
		}
		ta.Font = draw.Font{Name: font.Text, Size: float32(sz)}
	}
	font.Action = update
	size.Action = update

	form := NewForm()
	form.AddField("TextField:", NewTextField())
	form.AddField("CheckBox:", NewCheckBox("Enabled"))
	cb := NewComboBox()
	cb.Items = []ListItem{
		{Text: "Text", Icon: "file.text"},
		{Text: "Image", Icon: "file.image"},
		{Text: "Audio", Icon: "file.audio"},
		{Text: "Video", Icon: "file.video"},
		{Text: "Archive", Icon: "file.archive"},
		{Text: "Executable", Icon: "file.exec"},
	}
	form.AddField("ComboBox:", cb)
	form.AddField("Theme:", themeButton)

	text := &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(100,
			NewButtonIcon("cut", "", ta.Cut),
			NewButtonIcon("copy", "", ta.Copy),
			NewButtonIcon("paste", "", ta.Paste),
			NewSeparator(3, 1),
			NewLabel(" Font: "), font,
			NewLabel(" Size: "), size,
		),
	}
	root := NewRoot(&Container{
		Top:    mb,
		Center: NewHorizontalDivider(NewScrollView(form), NewShadow(text)),
	})
	themeButton.Action = func(state *ui.State) {
		if DefaultTheme == LightTheme {
			root.SetTheme(DarkTheme)
			DefaultTheme = DarkTheme
			themeButton.Text = "Dark Theme"
		} else {
			root.SetTheme(LightTheme)
			DefaultTheme = LightTheme
			themeButton.Text = "Light Theme"
		}
	}
	sdl.Show(sdl.Options{
		Title:      "Test",
		Height:     500,
		FontLookup: gofont.Lookup,
		IconLookup: &icons.Lookup{},
		Root:       root,
		Update:     ui.HandleKeyboardShortcuts,
	})
}

const todo = `TODO:
word wrap for text area
improve file chooser
	file type detection?`

const ipsum = `Lorem ipsum dolor sit amet, consectetur
adipiscing elit. Nam ac massa et leo
vestibulum euismod. Sed neque nisi,
consectetur at magna in, convallis
fermentum metus. Sed vehicula metus urna,
vel dapibus sem vestibulum quis. Mauris
faucibus nunc erat. Curabitur laoreet
dictum turpis, ut ornare felis facilisis
at. Suspendisse volutpat aliquet erat,
vitae feugiat risus ultrices eu. Nulla a
interdum leo. Etiam vitae rutrum mauris.
Proin ex justo, tempor vel porta id,
interdum ut purus. Fusce eget libero eget
sapien ultricies facilisis non vel leo.
Duis dui urna, porttitor eu facilisis eget,
imperdiet a nulla. Nullam lorem ante,
ornare quis bibendum ac, maximus id elit.
Curabitur sollicitudin maximus felis, nec
faucibus dui scelerisque quis. Maecenas
eget odio venenatis, varius magna et,
tempor lorem. Curabitur et urna
condimentum, tempor arcu ac, vestibulum
tortor.

Sed mattis lectus ex, eu rhoncus leo
iaculis id. Aliquam erat volutpat. Donec
auctor sapien blandit turpis accumsan
auctor eu eu erat. Sed blandit, nibh
vehicula tristique aliquet, orci nisl
molestie neque, tincidunt maximus quam odio
luctus lacus. Etiam eu laoreet lectus.
Pellentesque congue mollis tristique. Ut
nunc massa, mattis nec justo ac, mattis
semper tellus.

Cras sed erat eu magna vehicula vulputate
eu ac justo. Donec magna est, fermentum non
ex ac, ultricies mattis enim. Integer
pharetra, enim eget fringilla aliquet, ex
eros tempor sem, in porta magna lorem et
ante. Donec eleifend, tortor a tristique
eleifend, velit risus imperdiet metus,
vitae lacinia enim odio a lorem. Duis
mattis bibendum porttitor. Praesent vitae
velit dui. Nam id interdum massa. Ut dictum
justo in nulla tristique pellentesque. Nunc
vel nisl fringilla, condimentum ligula ut,
auctor arcu. Quisque non risus ut felis
placerat venenatis eget ac enim.
Suspendisse at erat quis diam bibendum
rhoncus.

Morbi eu diam orci. Sed sollicitudin luctus
mollis. In hac habitasse platea dictumst.
In feugiat vulputate dui sit amet
sollicitudin. Phasellus id leo sapien.
Nulla dapibus vel est a rhoncus. Nulla
facilisi. Suspendisse sed urna ut ex
maximus rhoncus sed ut lacus. Duis suscipit
libero ut eleifend sollicitudin. Mauris
interdum rutrum fringilla. Phasellus in
orci a leo consectetur vehicula non
vehicula mauris. Sed pharetra nec neque
quis tempus. Donec massa diam, varius eget
metus eget, convallis mattis tellus. Sed
ultricies semper mauris sit amet euismod.

Donec dapibus tincidunt volutpat. Proin
finibus vitae metus a venenatis. Praesent
tincidunt urna ac quam accumsan, quis
ultricies tellus pharetra. Nam in orci
purus. Vestibulum non diam lectus. In hac
habitasse platea dictumst. Interdum et
malesuada fames ac ante ipsum primis in
faucibus. In tellus massa, consequat quis
fermentum vitae, laoreet in metus. Duis
auctor nec urna ac elementum. Integer
consequat, augue et molestie commodo, metus
tellus ultricies neque, ac venenatis nisi
magna a leo. Suspendisse et odio eleifend,
ornare libero at, iaculis risus. In nec
arcu a augue fermentum posuere nec eget
elit.
`
