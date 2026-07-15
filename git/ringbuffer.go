// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package git

// ringBuffer is a buffer focusing on pushing an item to end and retrieving
// the sorted slice at once, without popping or removing.
//
// As git (git-go) does not have a way to load N commits until commit X, we
// have to read from a tip till we find the commit X in order to retrieve
// commits for "?after=X" query. std has no suitable container for that.
//   - container/list ... Have to remember both ends then call PushBack and
//     Remove.
//   - container/ring ... Has no beggining or end, heavily relies on
//     reassignment / shadowing. From document, it's
//     unclear what'll happen to unused slots.
//
// I just want "a buffer holds last N pushed items." This is a naive implemention
// that solves the problem.
type ringBuffer[T any] struct {
	buffer []T
	size   uint32
	cursor uint32
}

func newRingBuffer[T any](size uint32) ringBuffer[T] {
	return ringBuffer[T]{
		buffer: make([]T, 0, size),
		size:   size,
		cursor: 0,
	}
}

func (r *ringBuffer[T]) push(item T) {
	if uint32(len(r.buffer)) < r.size {
		// Until buffer is filled up, just append an item as usual.
		//
		//   [1|2|3| | | ]
		//    ^
		//
		//   [1|2|3|4| | ]
		//    ^
		r.buffer = append(r.buffer, item)
	} else {
		// When buffer is full, we have to move the cursor and put an item behind
		// the cursor.
		//
		//   [1|2|3|4|5|6]
		//    ^
		//
		//   [7|2|3|4|5|6]
		//      ^
		r.buffer[r.cursor] = item
		r.cursor += 1
		if r.cursor >= r.size {
			r.cursor = 0
		}
	}
}

// toSlice returns sorted slice of last (size) items.
func (r ringBuffer[T]) toSlice() []T {
	s := make([]T, len(r.buffer))
	for i := range len(r.buffer) {
		idx := (uint32(i) + r.cursor) % r.size
		s[i] = r.buffer[idx]
	}
	return s
}
