// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import "fmt"

func ExampleSet_Add() {
	var s Set[int] // int is just an example, this could be a struct
	s.Add(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	s.Send(999)

	// Output:
	// 999
}

func ExampleSyncSet_Add() {
	var ss SyncSet[int] // int is just an example, this could be a struct
	ss.Add(
		AsListener[int](func(event int) {
			fmt.Println(event)
		}),
	)

	ss.Send(999)

	// Output:
	// 999
}
