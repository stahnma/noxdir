package drive

import (
	"cmp"
	"maps"
	"slices"
)

type SortKey string

const (
	TotalCap   SortKey = "t"
	TotalUsed  SortKey = "u"
	TotalFree  SortKey = "f"
	TotalUsedP SortKey = "g"
)

// Info contains information about a single drive/volume/mount. The data can be
// fetched using a corresponding operating system syscall against a target entry.
type Info struct {
	Path        string  `json:"path"`
	Volume      string  `json:"volume"`
	FSName      string  `json:"fsName"`
	TotalBytes  uint64  `json:"total"`
	FreeBytes   uint64  `json:"free"`
	UsedBytes   uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type List struct {
	pathInfoMap   map[string]*Info
	TotalCapacity uint64
	TotalUsed     uint64
	TotalFree     uint64
}

func (l *List) All() map[string]*Info {
	return l.pathInfoMap
}

func (l *List) DriveInfo(path string) *Info {
	return l.pathInfoMap[path]
}

func (l *List) Find(diskPath string) *Info {
	for disk := range maps.Values(l.pathInfoMap) {
		if disk.Path == diskPath {
			return disk
		}
	}

	return nil
}

func (l *List) Sort(sk SortKey, desc bool) []*Info {
	drives := make([]*Info, 0, len(l.pathInfoMap))

	for disk := range maps.Values(l.pathInfoMap) {
		drives = append(drives, disk)
	}

	slices.SortFunc(
		drives,
		func(a, b *Info) int {
			var compared int

			switch sk {
			case TotalCap:
				compared = cmp.Compare(a.TotalBytes, b.TotalBytes)
			case TotalUsed:
				compared = cmp.Compare(a.UsedBytes, b.UsedBytes)
			case TotalFree:
				compared = cmp.Compare(a.FreeBytes, b.FreeBytes)
			case TotalUsedP:
				compared = cmp.Compare(a.UsedPercent, b.UsedPercent)
			}

			if desc {
				compared *= -1
			}

			return compared
		},
	)

	return drives
}
