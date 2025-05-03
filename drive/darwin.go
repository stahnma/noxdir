//go:build darwin

package drive

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/sys/unix"
)

var excludedVolumes = map[string]struct{}{
	"/":                          {},
	"/dev":                       {},
	"/System/Volumes/VM":         {},
	"/System/Volumes/Preboot":    {},
	"/System/Volumes/Update":     {},
	"/System/Volumes/xarts":      {},
	"/System/Volumes/iSCPreboot": {},
	"/System/Volumes/Hardware":   {},
	"/System/Volumes/Data/home":  {},
}

var excludedPaths = map[string]struct{}{
	"/System/Volumes/Data/Volumes": {},
}

func NewList() (*List, error) {
	mounts, err := mntList()
	if err != nil {
		return nil, err
	}

	list := &List{pathInfoMap: make(map[string]*Info, len(mounts))}

	for _, mount := range mounts {
		info := statFSToInfo(&mount)

		if _, excluded := excludedVolumes[info.Path]; excluded {
			continue
		}

		list.pathInfoMap[info.Path] = info
		list.TotalCapacity += info.TotalBytes
		list.TotalFree += info.FreeBytes
		list.TotalUsed += info.UsedBytes
	}

	return list, nil
}

func mntList() ([]unix.Statfs_t, error) {
	count, err := unix.Getfsstat(nil, unix.MNT_NOWAIT)
	if err != nil {
		return nil, fmt.Errorf("error getting fsstat: %v", err)
	}

	fs := make([]unix.Statfs_t, count)

	if _, err = unix.Getfsstat(fs, unix.MNT_NOWAIT); err != nil {
		return nil, fmt.Errorf("error getting fsstat: %v", err)
	}

	return fs, nil
}

func statFSToInfo(stat *unix.Statfs_t) *Info {
	usedBlocks := stat.Blocks - stat.Bfree

	return &Info{
		Path:        bytePtrToString(stat.Mntonname[:]),
		FSName:      bytePtrToString(stat.Fstypename[:]),
		TotalBytes:  stat.Blocks * uint64(stat.Bsize),
		FreeBytes:   stat.Bfree * uint64(stat.Bsize),
		UsedBytes:   usedBlocks * uint64(stat.Bsize),
		UsedPercent: (float64(usedBlocks) / float64(stat.Blocks)) * 100,
	}
}

func ReadDir(path string) ([]FileInfo, error) {
	fis := make([]FileInfo, 0, 32)

	entry, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = entry.Close()
	}()

	nodeEntries, err := entry.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, child := range nodeEntries {
		if child.IsDir() {
			fullPath := path

			if fullPath[len(fullPath)-1] != filepath.Separator {
				fullPath += string(filepath.Separator)
			}

			fullPath += child.Name()

			if _, excluded := excludedPaths[fullPath]; excluded {
				continue
			}
		}

		fis = append(
			fis,
			FileInfo{
				name:    child.Name(),
				isDir:   child.IsDir(),
				size:    child.Size(),
				modTime: child.ModTime(),
			},
		)
	}

	return fis, nil
}

func Explore(path string) error {
	if len(path) == 0 {
		return nil
	}

	cmd := exec.Command("open", path)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("error starting open: %v", err)
	}

	go func() {
		_ = cmd.Wait()
	}()

	return nil
}

func bytePtrToString(b []byte) string {
	for n := 0; n < len(b); n++ {
		if b[n] == 0 {
			return string(b[:n])
		}
	}

	return ""
}
