// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

// Listener is a sink for events.
type Listener[E any] interface {
	// OnEvent allows this listener to receive the given event.
	// This method must not panic, and it must not call any methods
	// on the containing Dispatcher.
	//
	// This method may be called concurrently from different goroutines.
	OnEvent(E)
}

// ListenerFunc is a type constraint for closures which can act as Listeners.
type ListenerFunc[E any] interface {
	~func(E)
}

type listenerFuncAdaptor[E any, F ListenerFunc[E]] struct {
	f F
}

func (a *listenerFuncAdaptor[E, F]) OnEvent(e E) { a.f(e) }

// AsListener converts a closure into a Listener.
//
// Use of this function is optional, as closures can be added to a Dispatcher directly.
// However, functions are not comparable in golang. Thus, if a caller adds a closure
// to a Dispatcher directly, it cannot be removed because the == operator
// isn't defined for functions. Using AsListener gives a caller a comparable
// Listener that can be removed from a Dispatcher at a future time.
func AsListener[E any, F ListenerFunc[E]](f F) Listener[E] {
	return &listenerFuncAdaptor[E, F]{
		f: f,
	}
}

// ListenerChan is a channel type that can act as a Listener.
type ListenerChan[E any] chan<- E

// OnEvent puts the event onto the channel. This method will block
// will block if the channel's queue is full or was created with no queue.
func (lc ListenerChan[E]) OnEvent(e E) { lc <- e }
