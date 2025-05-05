// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

// Dispatcher is the common interface for anything which can manage a
// collection of Listeners and dispatch events to them.
//
// A Dispatcher does not guarantee any ordering for listeners. In particular,
// the order in which listeners were added is not necessarily the order in
// which they will be called.
//
// A Dispatcher implementation must be safe for concurrent access. Any
// of the methods on this interface may be called concurrently at any time.
type Dispatcher[E any] interface {
	// Clear removes all listeners from this dispatcher.
	Clear()

	// AddListeners adds listeners to this Dispatcher.
	//
	// A caller must take care to ensure that any added listener that
	// could be removed later should be comparable. In particular, functions
	// in golang are NOT comparable. Thus, a function that implements the
	// Listener[E] interface cannot be passed to RemoveListeners.
	AddListeners(...Listener[E])

	// RemoveListeners removes listeners from this dispatcher. Only listeners
	// that are comparable may be removed. In particular, closure types which
	// implement Listener[E] cannot be used with this method.
	RemoveListeners(...Listener[E])

	// Send dispatches the event to all listeners currently associated
	// with this dispatcher.
	Send(E)
}

// AddListeners adds custom listeners to a dispatcher. The custom listener type
// can be anything that implements Listener[E], rather than being exactly Listener[E].
func AddListeners[E any, L Listener[E]](d Dispatcher[E], ls ...L) {
	switch len(ls) {
	case 0:
		// do nothing

	case 1:
		// simple optimization
		d.AddListeners(ls[0])

	default:
		// we want to make adding a chunk of listeners atomic in the
		// case where AddListeners is synchronized
		more := make([]Listener[E], len(ls))
		for i, l := range ls {
			more[i] = l
		}

		d.AddListeners(more...)
	}
}

// RemoveListeners removes custom listeners from a dispatcher. The custom listener type
// can be anything that implements Listener[E], rather than being exactly Listener[E].
func RemoveListeners[E any, L Listener[E]](d Dispatcher[E], ls ...L) {
	switch len(ls) {
	case 0:
		// do nothing

	case 1:
		// simple optimization
		d.RemoveListeners(ls[0])

	default:
		// we want to make adding a chunk of listeners atomic in the
		// case where AddListeners is synchronized
		more := make([]Listener[E], len(ls))
		for i, l := range ls {
			more[i] = l
		}

		d.RemoveListeners(more...)
	}
}
