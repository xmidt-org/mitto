// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "fmt"

func ExampleListeners_AddListeners() {
	var ls Listeners[int] // int is just an example, this could be a struct
	ls.AddListeners(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	ls.Send(999)

	// Output:
	// 999
}

func ExampleSyncListeners_AddListeners() {
	var ls SyncListeners[int] // int is just an example, this could be a struct
	ls.AddListeners(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	ls.Send(999)

	// Output:
	// 999
}
