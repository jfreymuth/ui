package icons

import (
	"image"
	"image/draw"
	"strings"

	"golang.org/x/exp/shiny/iconvg"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Lookup struct {
	z iconvg.Rasterizer
}

func (l *Lookup) IconSize(s int) int {
	switch {
	case s < 24:
		return 18
	case s < 36:
		return 24
	case s < 48:
		return 36
	}
	return 48
}

func (l *Lookup) DrawIcon(dst *image.Alpha, name string) {
	l.z.SetDstImage(dst, dst.Rect, draw.Src)
	iconvg.Decode(&l.z, findIcon(name), nil)
}

func findIcon(name string) []byte {
	if b, ok := iconData[name]; ok {
		return b
	}
	if i := strings.LastIndexByte(name, '.'); i >= 0 {
		return findIcon(name[:i])
	}
	return icons.ActionHelp
}

var iconData = map[string][]byte{
	"down":        icons.NavigationExpandMore,
	"down.arrow":  icons.NavigationArrowDownward,
	"up":          icons.NavigationExpandMore,
	"up.arrow":    icons.NavigationArrowUpward,
	"left":        icons.NavigationChevronLeft,
	"left.arrow":  icons.NavigationArrowBack,
	"right":       icons.NavigationChevronRight,
	"right.arrow": icons.NavigationArrowForward,
	"zoomIn":      icons.ActionZoomIn,
	"zoomOut":     icons.ActionZoomOut,

	"add":              icons.ContentAdd,
	"add.file":         icons.ActionNoteAdd,
	"remove":           icons.ContentRemove,
	"delete":           icons.ActionDelete,
	"remove.backspace": icons.ContentBackspace,
	"edit":             icons.ImageEdit,

	"save":         icons.ContentSave,
	"open":         icons.FileFolderOpen,
	"folder":       icons.FileFolderOpen,
	"file":         icons.EditorInsertDriveFile,
	"file.text":    icons.ActionDescription,
	"file.image":   icons.ImagePhoto,
	"file.audio":   icons.ImageMusicNote,
	"file.video":   icons.MapsLocalMovies,
	"file.exec":    icons.ActionSettingsApplications,
	"file.archive": icons.FileFolder,

	"info":     icons.ActionInfo,
	"question": icons.ActionHelp,
	"error":    icons.AlertError,
	"warning":  icons.AlertWarning,

	"lock":   icons.ActionLock,
	"unlock": icons.ActionLockOpen,

	"play":         icons.AVPlayArrow,
	"pause":        icons.AVPause,
	"stop":         icons.AVStop,
	"rewind":       icons.AVFastRewind,
	"fastForward":  icons.AVFastForward,
	"skipNext":     icons.AVSkipNext,
	"skipPrevious": icons.AVSkipPrevious,

	"volumeNone": icons.AVVolumeMute,
	"volumeLow":  icons.AVVolumeDown,
	"volumeHigh": icons.AVVolumeUp,
	"volumeOff":  icons.AVVolumeOff,

	"search":  icons.ActionSearch,
	"refresh": icons.NavigationRefresh,

	"menu":     icons.NavigationMenu,
	"settings": icons.ActionSettings,
	"close":    icons.NavigationClose,

	"cut":   icons.ContentContentCut,
	"copy":  icons.ContentContentCopy,
	"paste": icons.ContentContentPaste,
	"undo":  icons.ContentUndo,
	"redo":  icons.ContentRedo,

	"checkbox":              icons.ToggleCheckBoxOutlineBlank,
	"checkboxChecked":       icons.ToggleCheckBox,
	"checkboxIndeterminate": icons.ToggleIndeterminateCheckBox,
	"radiobutton":           icons.ToggleRadioButtonUnchecked,
	"radiobuttonSelected":   icons.ToggleRadioButtonChecked,
}
