package main

import (
	"container/heap"
	"sync"
)

// topFiles contains a list of the biggest files found on a specific drive/volume.
// It implements the heap.Interface, hence it persists the files sorted by their
// sizes. The topFiles could contain up to n files, where n is defined by the
// topFiles.size. Other files will be discarded.
type topFiles struct {
	size  int
	files []*Entry
	mx    sync.RWMutex
}

// PushSafe provides a thread-safe method for adding elements to the queue.
// On each call, it will check the current number of items and pop the oldest.
func (tf *topFiles) PushSafe(e *Entry) {
	tf.mx.Lock()
	defer tf.mx.Unlock()

	heap.Push(tf, e)

	if tf.Len() > tf.size {
		heap.Pop(tf)
	}
}

func (tf *topFiles) Reset() {
	tf.files = make([]*Entry, 0)
}

func (tf *topFiles) Less(i, j int) bool {
	return tf.files[i].Size < tf.files[j].Size
}

func (tf *topFiles) Swap(i, j int) {
	tf.files[i], tf.files[j] = tf.files[j], tf.files[i]
}

func (tf *topFiles) Len() int {
	return len(tf.files)
}

func (tf *topFiles) Pop() (v any) {
	v, tf.files = tf.files[tf.Len()-1], tf.files[:tf.Len()-1]

	return
}

func (tf *topFiles) Push(v any) {
	tf.files = append(tf.files, v.(*Entry))
}
