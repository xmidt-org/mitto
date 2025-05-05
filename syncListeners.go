// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "sync"

// SyncListeners is a Dispatcher backed by a slice of Listeners that is safe
// for concurrent access. A SyncListeners must not be copied after creation.
type SyncListeners[E any] struct {
	lock      sync.RWMutex
	listeners Listeners[E]
}

// Clear atomically removes all listeners.
func (sl *SyncListeners[E]) Clear() {
	sl.lock.Lock()
	sl.listeners.Clear()
	sl.lock.Unlock()
}

// AddListeners adds more listeners. This method ensures that no event can be sent
// until adding these listeners completes.
func (sl *SyncListeners[E]) AddListeners(toAdd ...Listener[E]) {
	sl.lock.Lock()
	sl.listeners.AddListeners(toAdd...)
	sl.lock.Unlock()
}

// RemoveListeners removes the given listeners. This method ensures that no event
// can be sent until removing these listeners completes.
func (sl *SyncListeners[E]) RemoveListeners(toRemove ...Listener[E]) {
	sl.lock.Lock()
	sl.listeners.RemoveListeners(toRemove...)
	sl.lock.Unlock()
}

// Send dispatches the given event to the contained listeners. Multiple goroutines
// can invoke this method concurrently. While an event is being sent, any method that
// would alter the set of Listeners is blocked.
//
// Listener implementations must be prepared to receive events concurrently. This method
// does not make any guarantees about the concurrent safety of the contained Listeners.
func (sl *SyncListeners[E]) Send(e E) {
	defer sl.lock.RUnlock() // in case a listener panics
	sl.lock.RLock()
	sl.listeners.Send(e)
}
