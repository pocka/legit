// Copyright 2026 Shota FUJI <pockawoooh@gmail.com>
// SPDX-License-Identifier: MIT

package git

import (
	"slices"
	"testing"
)

func TestRingBufferUnfilledOK(t *testing.T) {
	q := newRingBuffer[int](5)
	q.push(1)
	q.push(2)
	q.push(3)

	s := q.toSlice()

	if len(s) != 3 {
		t.Errorf("Expected slice length of 3, got %d", len(s))
	}

	if !slices.Equal(s, []int{1, 2, 3}) {
		t.Error("Incorrect slice returned")
	}
}

func TestRingBufferWrapping(t *testing.T) {
	q := newRingBuffer[int](5)
	for i := range 8 {
		q.push(i + 1)
	}

	s := q.toSlice()

	if len(s) != 5 {
		t.Errorf("Expected slice length of 5, got %d", len(s))
	}

	if !slices.Equal(s, []int{4, 5, 6, 7, 8}) {
		t.Error("Incorrect slice returned")
	}
}

func TestRingBufferNotCrashingOnEmptySlice(t *testing.T) {
	q := newRingBuffer[int](5)
	s := q.toSlice()

	if len(s) != 0 {
		t.Errorf("Expected slice length of 0, got %d", len(s))
	}

	if !slices.Equal(s, []int{}) {
		t.Error("Incorrect slice returned")
	}
}

func TestRingBufferManyRounds(t *testing.T) {
	q := newRingBuffer[int](3)
	for i := range 12 {
		q.push(i + 1)
	}

	s := q.toSlice()

	if len(s) != 3 {
		t.Errorf("Expected slice length of 3, got %d", len(s))
	}

	if !slices.Equal(s, []int{10, 11, 12}) {
		t.Error("Incorrect slice returned")
	}
}
