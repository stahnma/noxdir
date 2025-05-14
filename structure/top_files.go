package structure

import (
	"container/heap"
	"sync"
)

// TopFiles contains a list of the biggest files found on a specific drive/volume.
// It implements the heap.Interface, hence it persists the files sorted by their
// sizes. The TopFiles could contain up to n files, where n is defined by the
// TopFiles.size. Other files will be discarded.
type TopFiles struct {
	files []*Entry
	mx    sync.RWMutex
	size  int
}

// PushSafe provides a thread-safe method for adding elements to the queue.
// On each call, it will check the current number of items and pop the oldest.
func (tf *TopFiles) PushSafe(e *Entry) {
	tf.mx.Lock()
	defer tf.mx.Unlock()

	heap.Push(tf, e)

	if tf.Len() > tf.size {
		heap.Pop(tf)
	}
}

func (tf *TopFiles) Reset() {
	tf.files = make([]*Entry, 0)
}

func (tf *TopFiles) Less(i, j int) bool {
	return tf.files[i].Size < tf.files[j].Size
}

func (tf *TopFiles) Swap(i, j int) {
	tf.files[i], tf.files[j] = tf.files[j], tf.files[i]
}

func (tf *TopFiles) Len() int {
	return len(tf.files)
}

func (tf *TopFiles) Pop() (v any) {
	v, tf.files = tf.files[tf.Len()-1], tf.files[:tf.Len()-1]

	return
}

func (tf *TopFiles) Push(v any) {
	entry, ok := v.(*Entry)
	if !ok {
		return
	}

	tf.files = append(tf.files, entry)
}
