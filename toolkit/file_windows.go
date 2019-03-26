package toolkit

import (
	"fmt"
	"path/filepath"
	"syscall"

	"github.com/jfreymuth/ui"
)

func (f *FileChooser) back(*ui.State) {
	dir := filepath.Dir(f.path)
	if dir != f.path {
		f.set(dir)
	} else {
		kernel32, _ := syscall.LoadLibrary("kernel32.dll")
		getLogicalDrivesHandle, _ := syscall.GetProcAddress(kernel32, "GetLogicalDrives")

		if ret, _, errno := syscall.Syscall(uintptr(getLogicalDrivesHandle), 0, 0, 0, 0); errno == 0 {
			f.files.Clear()
			for i := 0; i < 26; i++ {
				if ret&(1<<uint(i)) != 0 {
					f.files.AddItemIcon("folder", fmt.Sprintf("%c:\\", 'A'+i))
				}
			}
		}
	}
}
