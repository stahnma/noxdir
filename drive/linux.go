//go:build linux

package drive

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	path2 "path"
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

var excludedPaths = map[string]struct{}{
	"/proc":       {},
	"/mnt":        {},
	"/sys":        {},
	"/lost+found": {},
	"/boot":       {},
	"/":           {},
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
		if err != nil {
			return nil, err
		}

		if excluded {
			continue
		}

		list.pathInfoMap[mntList[i]] = info
		list.TotalCapacity += info.TotalBytes
		list.TotalFree += info.FreeBytes
		list.TotalUsed += info.UsedBytes
	}

	return list, err
}

func mntInfo(path string) (*Info, bool, error) {
	var stat unix.Statfs_t

	if err := unix.Statfs(path, &stat); err != nil {
		return nil, false, fmt.Errorf("failed to statfs: %v", err)
	}

	fsName, _ := fsTypesMap[int64(stat.Type)]

	// use an implicitly defined list of excluded FS types rather than names map
	if _, ok := excludedFSTypes[int64(stat.Type)]; ok || stat.Blocks == 0 {
		return nil, true, nil
	}

	usedBlocks := (stat.Blocks - stat.Bfree) * uint64(stat.Bsize)

	info := &Info{
		Path:        path,
		FSName:      fsName,
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

func ReadDir(path string) ([]os.FileInfo, error) {
	var infos []os.FileInfo

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
		if _, exclude := excludedPaths[path2.Join(path, child.Name())]; !exclude {
			infos = append(infos, child)
		}
	}

	return infos, nil
}

func Explore(path string) error {
	return nil
}
