// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "slices"

// Set is a Dispatcher backed by a simple slice of Listener. The zero
// value for this type is ready to use.
//
// A Set is not safe for concurrent use. This type can be used
// in situations where concurrent safety is guaranteed by containing code
// or where concurrency isn't an issue.
type Set[E any] struct {
	all []Listener[E]
}

// Clear removes all listeners from this set.
func (s *Set[E]) Clear() {
	for i := range len(s.all) {
		s.all[i] = nil
	}

	s.all = s.all[:0]
}

// Add appends listeners to this set. Any nil listeners
// are skipped.
//
// AsListener can be used to convert closures and channels into
// listeners to pass to this method.
func (s *Set[E]) Add(toAdd ...Listener[E]) {
	s.all = slices.Grow(s.all, len(toAdd))
	for _, nl := range toAdd {
		if nl != nil {
			s.all = append(s.all, nl)
		}
	}
}

// Remove deletes the given listeners. Nil listeners and listeners
// that are not part of this set are ignored.
func (s *Set[E]) Remove(toRemove ...Listener[E]) {
	for _, r := range toRemove {
		if p := slices.Index(s.all, r); p >= 0 {
			last := len(s.all) - 1
			s.all[p], s.all[last] = s.all[last], nil
			s.all = s.all[:last]
		}
	}
}

// Send dispatches an event to all contained listeners. Listener implementations
// should be prepared to receive events concurrently.
func (s *Set[E]) Send(e E) {
	for _, l := range s.all {
		l.OnEvent(e)
	}
}
