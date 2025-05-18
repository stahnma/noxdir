package drive

import "sync"

// InoFilter filters files by inode value, preventing the double calculation of
// the same space on the disk. The filter must be manually reset before each new
// volume scan.
type InoFilter struct {
	inoMap map[uint64]struct{}
	mx     sync.Mutex
}

// Add adds a new inode value to the filter. It returns a bool value depending on
// whether the inode already exists - "false", or adds it to the filter - "true".
func (inf *InoFilter) Add(inode uint64) bool {
	inf.mx.Lock()
	defer inf.mx.Unlock()

	if _, ok := inf.inoMap[inode]; ok {
		return false
	}

	inf.inoMap[inode] = struct{}{}

	return true
}

func (inf *InoFilter) Reset() {
	inf.inoMap = make(map[uint64]struct{})
}

var InoFilterInstance = InoFilter{
	inoMap: make(map[uint64]struct{}),
}
