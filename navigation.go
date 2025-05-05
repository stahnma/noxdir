package main

import (
	"sync/atomic"
	"time"

	"github.com/crumbyte/noxdir/drive"
)

// State defines a custom type representing a current GUI state.
type State int

const (
	Drives State = iota
	Dirs
)

type ChangeLevelHandler func(e *Entry, s State)

type Navigation struct {
	entry        *Entry
	drives       *drive.List
	currentDrive *drive.Info
	entryStack   []*Entry
	state        State
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

func (n *Navigation) Entry() *Entry {
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
		n.state = Drives

		return
	}

	n.entry = n.entryStack[len(n.entryStack)-1]
	n.entryStack = n.entryStack[:len(n.entryStack)-1]
}

func (n *Navigation) LevelDown(path string, clh ChangeLevelHandler) (chan struct{}, chan error) {
	if n.Lock() && len(path) == 0 {
		return nil, nil
	}

	if n.state == Drives {
		n.state = Dirs

		n.entry = NewDirEntry(path, time.Now())
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

	n.entryStack, n.entry = append(n.entryStack, n.entry), entry

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
