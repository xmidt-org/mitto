// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "fmt"

func ExampleListeners_Add() {
	var ls Listeners[int] // int is just an example, this could be a struct
	ls.Add(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	ls.Send(999)

	// Output:
	// 999
}

func ExampleSyncListeners_Add() {
	var ls SyncListeners[int] // int is just an example, this could be a struct
	ls.Add(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	ls.Send(999)

	// Output:
	// 999
}
