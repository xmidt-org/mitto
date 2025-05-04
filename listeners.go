// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "slices"

// Listeners is a Dispatcher backed by a simple slice of Listeners.
//
// A Listeners is not safe for concurrent use. This type can be used
// in situations where concurrent safety is guaranteed by containing code
// or where concurrency isn't an issue.
type Listeners[E any] struct {
	all []Listener[E]
}

func (ls *Listeners[E]) Clear() {
	for i := range len(ls.all) {
		ls.all[i] = nil
	}

	ls.all = ls.all[:0]
}

func (ls *Listeners[E]) AddListeners(toAdd ...Listener[E]) {
	ls.all = append(ls.all, toAdd...)
}

func (ls *Listeners[E]) AddListenerFuncs(toAdd ...func(E)) {
	ls.all = slices.Grow(ls.all, len(toAdd))
	for _, f := range toAdd {
		ls.all = append(ls.all,
			AsListener(f),
		)
	}
}

func (ls *Listeners[E]) AddListenerChans(toAdd ...chan<- E) {
	ls.all = slices.Grow(ls.all, len(toAdd))
	for _, c := range toAdd {
		ls.all = append(ls.all,
			ListenerChan[E](c),
		)
	}
}

func (ls *Listeners[E]) RemoveListeners(toRemove ...Listener[E]) {
	for _, r := range toRemove {
		if p := slices.Index(ls.all, r); p >= 0 {
			last := len(ls.all) - 1
			ls.all[p], ls.all[last] = ls.all[last], nil
			ls.all = ls.all[:last]
		}
	}
}

// Send dispatches an event to all contained listeners. Listener implementations
// should be prepared to receive events concurrently.
func (ls *Listeners[E]) Send(e E) {
	for _, l := range ls.all {
		l.OnEvent(e)
	}
}
