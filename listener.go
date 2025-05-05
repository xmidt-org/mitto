// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "reflect"

// Listener is a sink for events.
type Listener[E any] interface {
	// OnEvent allows this listener to receive the given event.
	// This method must not panic, and it must not call any methods
	// on the containing Dispatcher.
	//
	// This method may be called concurrently from different goroutines.
	OnEvent(E)
}

// Sink is a constraint for types that can be used as Listeners.
type Sink[E any] interface {
	~func(E) | ~chan E | ~chan<- E
}

type listenerFuncAdaptor[E any] struct {
	fn func(E)
}

func (a *listenerFuncAdaptor[E]) OnEvent(e E) { a.fn(e) }

type listenerChanAdaptor[E any] struct {
	ch chan<- E
}

func (a *listenerChanAdaptor[E]) OnEvent(e E) { a.ch <- e }

// AsListener converts a Sink into a Listener. If the supplied sink is nil,
// this function returns nil. The returned Listener is always comparable and can
// be passed to Dispatcher.RemoveListeners.
//
// For channels, the returned listener will block if the underlying channel blocks.
// Clients must create and manage channels to reduce or avoid blocking. Additionally,
// if a client wants to close the channel, care must be taken to remove it first
// to avoid panics.
func AsListener[E any, S Sink[E]](sink S) Listener[E] {
	if sink == nil {
		return nil
	}

	// fast conversions for simple types
	switch st := any(sink).(type) {
	case func(E):
		return &listenerFuncAdaptor[E]{fn: st}
	case chan E:
		return &listenerChanAdaptor[E]{ch: st}
	case chan<- E:
		return &listenerChanAdaptor[E]{ch: st}
	}

	// for custom types, use conversions through reflection
	sv := reflect.ValueOf(sink)

	ftype := reflect.TypeOf((func(E))(nil))
	if sv.CanConvert(ftype) {
		return &listenerFuncAdaptor[E]{
			fn: sv.Convert(ftype).Interface().(func(E)),
		}
	}

	chtype := reflect.TypeOf((chan E)(nil))
	if sv.CanConvert(chtype) {
		return &listenerChanAdaptor[E]{
			ch: sv.Convert(chtype).Interface().(chan E),
		}
	}

	sendtype := reflect.TypeOf((chan<- E)(nil))
	return &listenerChanAdaptor[E]{
		ch: sv.Convert(sendtype).Interface().(chan<- E),
	}
}
