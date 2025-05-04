// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "fmt"

func ExampleListeners_AddListeners() {
	var l Listeners[int] // int is just an example, this could be a struct
	l.AddListeners(
		AsListener(func(event int) {
			fmt.Println(event)
		}),
	)

	l.Send(999)

	// Output:
	// 999
}

func ExampleSyncListeners_AddListeners() {
	var l SyncListeners[int] // int is just an example, this could be a struct
	l.AddListeners(
		AsListener(func(event int) {
			fmt.Println(event)
		}),
	)

	l.Send(999)

	// Output:
	// 999
}
