// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "sync"

// SyncSet is a Dispatcher backed by a slice of Listener and is safe
// for concurrent access. A SyncSet must not be copied after creation.
//
// The zero value for this type is ready to use.
type SyncSet[E any] struct {
	lock sync.RWMutex
	set  Set[E]
}

// Clear removes all listeners. Send will block until this method completes.
func (ss *SyncSet[E]) Clear() {
	ss.lock.Lock()
	ss.set.Clear()
	ss.lock.Unlock()
}

// Add adds more listeners. This method ensures that no event can be sent
// until adding these listeners completes.
func (ss *SyncSet[E]) Add(toAdd ...Listener[E]) {
	ss.lock.Lock()
	ss.set.Add(toAdd...)
	ss.lock.Unlock()
}

// Remove removes the given listeners. This method ensures that no event
// can be sent until removing these listeners completes.
func (ss *SyncSet[E]) Remove(toRemove ...Listener[E]) {
	ss.lock.Lock()
	ss.set.Remove(toRemove...)
	ss.lock.Unlock()
}

// Send dispatches the given event to the contained listeners. Multiple goroutines
// can invoke this method concurrently. While an event is being sent, any method that
// would alter the set of Listeners is blocked.
//
// Listener implementations must be prepared to receive events concurrently. This method
// does not make any guarantees about the concurrent safety of the contained Listeners.
func (ss *SyncSet[E]) Send(e E) {
	defer ss.lock.RUnlock() // in case a listener panics
	ss.lock.RLock()
	ss.set.Send(e)
}
