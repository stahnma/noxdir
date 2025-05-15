package main

import (
	"sync/atomic"

	"github.com/crumbyte/noxdir/drive"
	"github.com/crumbyte/noxdir/structure"
)

// State defines a custom type representing a current GUI state.
type State int

const (
	Drives State = iota
	Dirs
)

type ChangeLevelHandler func(e *structure.Entry, s State)

type stackItem struct {
	entry  *structure.Entry
	cursor int
}

type Navigation struct {
	entry        *structure.Entry
	drives       *drive.List
	currentDrive *drive.Info
	entryStack   []*stackItem
	state        State
	cursor       int
	locked       atomic.Bool
}

func NewNavigation(l *drive.List) *Navigation {
	return &Navigation{
		state:  Drives,
		drives: l,
	}
}

func (n *Navigation) State() State {
	return n.state
}

func (n *Navigation) DrivesList() *drive.List {
	return n.drives
}

func (n *Navigation) Entry() *structure.Entry {
	return n.entry
}

func (n *Navigation) Locked() bool {
	return n.locked.Load()
}

func (n *Navigation) Lock() bool {
	return !n.locked.Swap(true)
}

func (n *Navigation) Unlock() {
	n.locked.Swap(false)
}

func (n *Navigation) ParentSize() int64 {
	if n.entry == nil {
		//nolint:gosec // why not, let's overflow
		return int64(n.currentDrive.UsedBytes)
	}

	return n.entry.Size
}

func (n *Navigation) LevelUp() {
	if n.state == Drives {
		return
	}

	if len(n.entryStack) == 0 {
		n.state, n.cursor = Drives, 0

		return
	}

	lastItem := n.entryStack[len(n.entryStack)-1]

	n.entry, n.cursor = lastItem.entry, lastItem.cursor
	n.entryStack = n.entryStack[:len(n.entryStack)-1]
}

func (n *Navigation) LevelDown(path string, cursor int, clh ChangeLevelHandler) (chan struct{}, chan error) {
	if n.Lock() && len(path) == 0 {
		return nil, nil
	}

	if n.state == Drives {
		n.state = Dirs

		n.entry = structure.NewDirEntry(path, 0)
		n.currentDrive = n.drives.DriveInfo(path)

		return n.entry.TraverseAsync()
	}

	defer func() {
		clh(n.entry, n.state)
	}()

	entry := n.entry.GetChild(path)
	if entry == nil || !entry.IsDir {
		return nil, nil
	}

	n.entryStack = append(
		n.entryStack,
		&stackItem{entry: n.entry, cursor: cursor},
	)

	n.entry = entry
	n.cursor = 0

	return nil, nil
}

func (n *Navigation) Explore(path string) error {
	if len(path) == 0 {
		return nil
	}

	var fullPath string

	if n.state == Drives {
		d := n.drives.DriveInfo(path)
		if d == nil {
			return nil
		}

		fullPath = d.Path
	} else {
		entry := n.entry.GetChild(path)
		if entry == nil {
			return nil
		}

		fullPath = entry.Path
	}

	return drive.Explore(fullPath)
}
