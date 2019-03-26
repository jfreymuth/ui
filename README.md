# package ui

Lightweight go desktop ui

**this package is highly experimental, drastic changes may occur**

### What this package tries to do

- be intuitive for users; applications created with this packge should superficially behave like other ui toolkits, even in subtle ways
- be intuitive for programmers; writing applications should be straightforward and not require much boilerplate or magic incantations
- make it easy to create custom components
- be the basis of a potential ecosystem; many parts of this package can easily be replaced, and the replacements should work well together (e.g. someone could write a different renderer, or an alternative library of components)
- not be dependent on cgo in principle; while the current backend uses cgo, one could write e.g. an X11 backend in pure go.

### What this package does *not* try to do

- be the "official" go ui package
- be useful in all situations (a big limitation right now is that applications can only use a single window)
- be suited for mobile or web applications

## Overview

The design of this package is inspired by immediate mode ui and uses similar techniques to minimize the complexity of components. Many ui libraries make heavy use of inheritance, which can be imitated in go through struct embedding. However this approach has disadvantages and is generally considered unidiomatic. This package follows a different approach by hiding most of the complexity in the `ui.State` struct instead of an abstract base class. This leads to the `Component` interface being very small, and components can be extremly lightweight.

There are three parts to this package:

### `ui` and `ui/draw`

`ui` defines important interfaces and constants. It also provides the type `State`, which is central to the functionality of the library.

`ui/draw` contains data structures for drawing. Components do not draw directly, but rather create a list of commands that will be passed on to the renderer.

### `ui/toolkit`

`ui/toolkit` is a collection of widgets. It provides common widgets (buttons, text fields etc.), simple layouts and basic themeing support.

### `ui/impl/...`

The packages in `ui/impl` make up the backend. They are designed to be replaced by better packages later. These are the only packages with external dependencies.

It consists of the following packages (so far):
- `ui/impl/gofont` supports truetype fonts via [github.com/golang/freetype](https://github.com/golang/freetype) and includes the [go fonts](https://blog.golang.org/go-fonts) as default fonts
- `ui/impl/icons` provides some of google's material design icons through [golang.org/x/exp/shiny/iconvg](https://godoc.org/golang.org/x/exp/shiny/iconvg)
- `ui/impl/gldraw` renders using OpenGL 3.3 ([github.com/go-gl/gl](https://github.com/go-gl/gl))
- `ui/impl/sdl` is the main backend, it handles window creation and user input. It uses [github.com/veandco/go-sdl2](https://github.com/veandco/go-sdl2)

## Usage

```go
import (
	"github.com/jfreymuth/ui"
	"github.com/jfreymuth/ui/impl/gofont"
	"github.com/jfreymuth/ui/impl/icons"
	"github.com/jfreymuth/ui/impl/sdl"
	"github.com/jfreymuth/ui/toolkit"
)

func main() {
	button := toolkit.NewButton("Button", func(state *ui.State) {
		toolkit.ShowMessageDialog(state, "Message", "The button was pressed!", "Ok")
	})

	root := toolkit.NewRoot(button)

	sdl.Show(sdl.Options{
		Title: "Example",
		Width: 400, Height: 300,
		FontLookup: gofont.Lookup,
		IconLookup: &icons.Lookup{},
		Root:       root,
	})
}
```

### Examples

- [`examples/editor`](examples/editor/main.go): a basic text editor, showcases the text area and demonstrates how to implement "save before quitting" comfirmations
- [`examples/synth`](examples/synth): a very simple synthesizer, demonstrates custom widgets and also some concurrency
- [`examples/test`](examples/test/main.go): rather messy, but uses almost every single widget
