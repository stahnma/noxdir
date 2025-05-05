//go:build linux

package drive

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/unix"
)

const NTFS_SB_MAGIC = 0x5346544e

const mountInfoPath = "/proc/self/mounts"

var excludedFSTypes = map[int64]struct{}{
	unix.CGROUP_SUPER_MAGIC:    {},
	unix.CGROUP2_SUPER_MAGIC:   {},
	unix.SYSFS_MAGIC:           {},
	unix.OVERLAYFS_SUPER_MAGIC: {},
	unix.TMPFS_MAGIC:           {},
	unix.DEBUGFS_MAGIC:         {},
	unix.SQUASHFS_MAGIC:        {},
	unix.PROC_SUPER_MAGIC:      {},
	unix.SECURITYFS_MAGIC:      {},
}

var excludedPaths = map[string]map[string]struct{}{
	"/": {
		"mnt":        {},
		"sys":        {},
		"lost+found": {},
		"boot":       {},
		"proc":       {},
	},
}

var fsTypesMap = map[int64]string{
	unix.EXT4_SUPER_MAGIC:  "ext4",
	unix.XFS_SUPER_MAGIC:   "xfs",
	unix.BTRFS_SUPER_MAGIC: "btrfs",
	unix.NFS_SUPER_MAGIC:   "nfs",
	unix.MSDOS_SUPER_MAGIC: "msdos",
	unix.V9FS_MAGIC:        "v9",
	NTFS_SB_MAGIC:          "ntfs",
}

func NewList() (*List, error) {
	mntList, err := mounts()
	if err != nil {
		return nil, err
	}

	list := &List{pathInfoMap: make(map[string]*Info, len(mntList))}

	for i := range mntList {
		info, excluded, err := mntInfo(mntList[i])
		// suppress an error mostly related to the permissions, and requires
		// root access.
		if excluded || err != nil {
			continue
		}

		list.pathInfoMap[mntList[i]] = info
		list.TotalCapacity += info.TotalBytes
		list.TotalFree += info.FreeBytes
		list.TotalUsed += info.UsedBytes
	}

	return list, err
}

func NewFileInfo(name string, data *unix.Stat_t) FileInfo {
	return FileInfo{
		name:    name,
		isDir:   data.Mode&unix.S_IFMT == unix.S_IFDIR,
		size:    data.Size,
		modTime: time.Unix(int64(data.Mtim.Sec), int64(data.Mtim.Nsec)),
	}
}

func mntInfo(path string) (*Info, bool, error) {
	var stat unix.Statfs_t

	if err := unix.Statfs(path, &stat); err != nil {
		return nil, false, fmt.Errorf("failed to statfs: %v", err)
	}

	// use an implicitly defined list of excluded FS types rather than names map
	if _, ok := excludedFSTypes[int64(stat.Type)]; ok || stat.Blocks == 0 {
		return nil, true, nil
	}

	usedBlocks := stat.Blocks - stat.Bfree

	info := &Info{
		Path:        path,
		FSName:      fsTypesMap[int64(stat.Type)],
		TotalBytes:  stat.Blocks * uint64(stat.Bsize),
		FreeBytes:   stat.Bfree * uint64(stat.Bsize),
		UsedBytes:   usedBlocks * uint64(stat.Bsize),
		UsedPercent: (float64(usedBlocks) / float64(stat.Blocks)) * 100,
	}

	return info, false, nil
}

func mounts() ([]string, error) {
	mnt, err := os.Open(mountInfoPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", mountInfoPath, err)
	}

	defer func(mnt *os.File) {
		_ = mnt.Close()
	}(mnt)

	scanner := bufio.NewScanner(mnt)

	var mntList []string

	for scanner.Scan() {
		mountInfo := scanner.Bytes()
		start := bytes.IndexByte(mountInfo, ' ')
		if start == -1 {
			continue
		}

		end := bytes.IndexByte(mountInfo[start+1:], ' ')
		if end == -1 {
			continue
		}

		mntList = append(
			mntList, string(mountInfo[start+1:len(mountInfo[:start+1])+end]),
		)
	}

	return mntList, nil
}

var direntBufPool = sync.Pool{
	New: func() any {
		b := make([]byte, 1024*64)

		return &b
	},
}

func ReadDir(path string) ([]FileInfo, error) {
	fd, err := unix.Open(path, unix.O_RDONLY|unix.O_DIRECTORY, 0)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	defer func(fd int) {
		_ = unix.Close(fd)
	}(fd)

	buf := direntBufPool.Get().(*[]byte)
	defer direntBufPool.Put(buf)

	fis := make([]FileInfo, 0, 32)

	for {
		n, err := unix.Getdents(fd, *buf)
		if err != nil {
			return nil, fmt.Errorf("getdents error: %w", err)
		}

		if n == 0 {
			break
		}

		offset := 0

		for offset < n {
			dirent := (*unix.Dirent)(unsafe.Pointer(&(*buf)[offset]))

			nameBytes := (*[256]byte)(unsafe.Pointer(&dirent.Name[0]))
			name := bytePtrToString(nameBytes[:])

			if name == "." || name == ".." {
				offset += int(dirent.Reclen)
				continue
			}

			if excludedChild, excluded := excludedPaths[path]; excluded {
				if _, childExcluded := excludedChild[name]; childExcluded {
					offset += int(dirent.Reclen)
					continue
				}
			}

			var stat unix.Stat_t

			err = unix.Fstatat(fd, name, &stat, unix.AT_SYMLINK_NOFOLLOW)
			if err == nil {
				fis = append(fis, NewFileInfo(name, &stat))
			}

			offset += int(dirent.Reclen)
		}
	}

	return fis, nil
}

func Explore(path string) error {
	if len(path) == 0 {
		return nil
	}

	cmd := exec.Command("xdg-open", path)

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
