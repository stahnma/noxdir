package drive

import (
	"os"
	"time"
)

// FileInfo defines a custom fs.FileInfo implementation for wrapping the results
// from the file info system calls.
type FileInfo struct {
	modTime time.Time
	name    string
	size    uint64
	isDir   bool
}

func (fi FileInfo) Name() string {
	return fi.name
}

func (fi FileInfo) Size() uint64 {
	return fi.size
}

func (fi FileInfo) Mode() os.FileMode {
	// since we are not using the os.FileMode values we can skip the mapping
	// from the Windows API file attributes.
	panic("os.FileMode not supported")
}

func (fi FileInfo) ModTime() time.Time {
	return fi.modTime
}

func (fi FileInfo) IsDir() bool {
	return fi.isDir
}

func (fi FileInfo) Sys() any {
	return nil
}
