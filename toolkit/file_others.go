// +build !windows

package toolkit

import (
	"path/filepath"

	"github.com/jfreymuth/ui"
)

func (f *FileChooser) back(*ui.State) {
	f.set(filepath.Dir(f.path))
}
