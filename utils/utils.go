// Copyright © 2021 Jacob Hansen. All rights reserved.
// Licensed under the MIT License. See LICENSE file in the project root for full license information.

// Package utils is a collection of utilities for Jobs.
package utils

import "sync"

// AtomicBool is a thread-safe wrapper for a bool.
type AtomicBool struct {
	value bool
	mu    *sync.Mutex
}

// NewAtomicBool creates a new thread-safe bool with the specified value.
func NewAtomicBool(initialVal bool) AtomicBool {
	safeVar := &AtomicBool{
		value: initialVal,
		mu:    new(sync.Mutex),
	}
	return *safeVar
}

// SafeSet is a thread-safe function to set the value of an AtomicBool.
func (s *AtomicBool) SafeSet(val bool) {
	defer s.mu.Unlock()
	s.mu.Lock()
	s.value = val
}

// SafeGet is a thread-safe function to get the value of an AtomicBool.
func (s *AtomicBool) SafeGet() bool {
	defer s.mu.Unlock()
	s.mu.Lock()
	return s.value
}
