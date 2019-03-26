package toolkit

import "github.com/jfreymuth/ui"

func ShowMessageDialog(state *ui.State, title, message, button string) {
	ta := NewTextArea()
	ta.Editable = false
	ta.SetText(message)
	state.OpenDialog(NewFrame("info", title, &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(-1, NewButton(button, (*ui.State).CloseDialog)),
	}, DefaultTheme.Color("titleBackground")))
}

func ShowErrorDialog(state *ui.State, title, message, button string) {
	ta := NewTextArea()
	ta.Editable = false
	ta.SetText(message)
	state.OpenDialog(NewFrame("warning", title, &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(-1, NewButton(button, (*ui.State).CloseDialog)),
	}, DefaultTheme.Color("titleBackgroundError")))
}

func ShowConfirmDialog(state *ui.State, title, message, ok, cancel string, action func(*ui.State)) {
	ta := NewTextArea()
	ta.Editable = false
	ta.SetText(message)
	state.OpenDialog(NewFrame("question", title, &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(-1, NewButton(ok, func(state *ui.State) {
			state.CloseDialog()
			if action != nil {
				action(state)
			}
		}), NewButton(cancel, (*ui.State).CloseDialog)),
	}, DefaultTheme.Color("titleBackground")))
}

func ShowYesNoDialog(state *ui.State, title, message, yes, no, cancel string, yesAction, noAction func(*ui.State)) {
	ta := NewTextArea()
	ta.Editable = false
	ta.SetText(message)
	state.OpenDialog(NewFrame("question", title, &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(-1, NewButton(yes, func(state *ui.State) {
			state.CloseDialog()
			if yesAction != nil {
				yesAction(state)
			}
		}), NewButton(no, func(state *ui.State) {
			state.CloseDialog()
			if noAction != nil {
				noAction(state)
			}
		}), NewButton(cancel, (*ui.State).CloseDialog)),
	}, DefaultTheme.Color("titleBackground")))
}

func ShowInputDialog(state *ui.State, title, message, button, cancel string, action func(*ui.State, string)) {
	ta := NewTextArea()
	ta.Editable = false
	ta.SetText(message)
	tf := NewTextField()
	tf.Action = func(state *ui.State, text string) {
		state.CloseDialog()
		if action != nil {
			action(state, text)
		}
	}
	state.OpenDialog(NewFrame("question", title, &Container{
		Center: NewScrollView(ta),
		Bottom: NewBar(-1, tf, NewButton(button, tf.TriggerAction), NewButton(cancel, (*ui.State).CloseDialog)),
	}, DefaultTheme.Color("titleBackground")))
}

func ShowOpenDialog(state *ui.State, fc *FileChooser, title, open, cancel string, action func(*ui.State, string)) {
	fc.SetLabels(open, cancel, true)
	fc.Action = action
	state.OpenDialog(NewFrame("open", title, fc, DefaultTheme.Color("titleBackground")))
}

func ShowSaveDialog(state *ui.State, fc *FileChooser, title, save, cancel string, action func(*ui.State, string)) {
	fc.SetLabels(save, cancel, false)
	fc.Action = action
	state.OpenDialog(NewFrame("save", title, fc, DefaultTheme.Color("titleBackground")))
}
