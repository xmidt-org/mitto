// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type testListener[E any] struct {
	*assert.Assertions
	called   bool
	expected E
}

func (tl *testListener[E]) OnEvent(actual E) {
	tl.called = true
	tl.Equal(tl.expected, actual)
}

// DispatcherTestSuite is a generic suite for any Dispatcher implementation.
type DispatcherTestSuite[E any, D Dispatcher[E]] struct {
	suite.Suite

	// factory is the closure that creates a dispatcher of the
	// type under test.
	factory func() D

	// testEvent is used as an expected event
	testEvent E
}

func (suite *DispatcherTestSuite[E, D]) newTestListener() *testListener[E] {
	return &testListener[E]{
		Assertions: suite.Assert(),
		expected:   suite.testEvent,
	}
}

func (suite *DispatcherTestSuite[E, D]) newTestListeners(n int) (tests []*testListener[E], toAdd []Listener[E]) {
	tests = make([]*testListener[E], n)
	toAdd = make([]Listener[E], n)

	for i := 0; i < n; i++ {
		tl := suite.newTestListener()
		tests[i] = tl
		toAdd[i] = tl
	}

	return
}

func (suite *DispatcherTestSuite[E, D]) newTestListenerFuncs(n int) (tests []*testListener[E], toAdd []func(E)) {
	tests = make([]*testListener[E], n)
	toAdd = make([]func(E), n)

	for i := 0; i < n; i++ {
		tl := suite.newTestListener()
		tests[i] = tl
		toAdd[i] = tl.OnEvent
	}

	return
}

func (suite *DispatcherTestSuite[E, D]) newTestListenerChans(n int) (tests []chan E, toAdd []chan<- E) {
	tests = make([]chan E, n)
	toAdd = make([]chan<- E, n)

	for i := 0; i < n; i++ {
		tests[i] = make(chan E, 1)
		toAdd[i] = tests[i]
	}

	return
}

func (suite *DispatcherTestSuite[E, D]) resetTestListeners(tests []*testListener[E]) {
	for _, tl := range tests {
		tl.called = false
		tl.expected = suite.testEvent
	}
}

func (suite *DispatcherTestSuite[E, D]) assertTestListenersCalled(tests []*testListener[E]) {
	for i, tl := range tests {
		suite.True(tl.called, "testListener at index [%d] not called", i)
	}
}

func (suite *DispatcherTestSuite[E, D]) assertTestListenersNotCalled(tests []*testListener[E]) {
	for i, tl := range tests {
		suite.False(tl.called, "testListener at index [%d] called", i)
	}
}

func (suite *DispatcherTestSuite[E, D]) assertTestListenerChansCalled(tests []chan E) {
	for i, ch := range tests {
		select {
		case actual := <-ch:
			suite.Equal(suite.testEvent, actual)

		default:
			suite.Failf("failed to send event on a test channel", "channel at index [%d]", i)
		}
	}
}

func (suite *DispatcherTestSuite[E, D]) assertTestListenerChansNotCalled(tests []chan E) {
	for i, ch := range tests {
		select {
		case <-ch:
			suite.Failf("unexpected event on a test channel", "channel at index [%d]", i)

		default:
			// passing
		}
	}
}

func (suite *DispatcherTestSuite[E, D]) TestEmpty() {
	d := suite.factory()

	// all of the following should be nops
	d.Clear()
	d.AddListeners()
	d.AddListenerFuncs()
	d.AddListenerChans()
	d.RemoveListeners()
	d.Send(suite.testEvent)
}

func (suite *DispatcherTestSuite[E, D]) TestAddListeners() {
	for _, count := range []int{1, 2, 5} {
		suite.Run(fmt.Sprintf("count=%d", count), func() {
			tests, toAdd := suite.newTestListeners(count)
			d := suite.factory()

			d.AddListeners(toAdd...)
			suite.resetTestListeners(tests)
			d.Send(suite.testEvent)
			suite.assertTestListenersCalled(tests)

			d.RemoveListeners(toAdd...)
			suite.resetTestListeners(tests)
			d.Send(suite.testEvent)
			suite.assertTestListenersNotCalled(tests)

			d.AddListeners(toAdd...)
			d.Clear()
			suite.resetTestListeners(tests)
			d.Send(suite.testEvent)
			suite.assertTestListenersNotCalled(tests)
		})
	}
}

func (suite *DispatcherTestSuite[E, D]) TestAddListenerFuncs() {
	for _, count := range []int{1, 2, 5} {
		suite.Run(fmt.Sprintf("count=%d", count), func() {
			tests, toAdd := suite.newTestListenerFuncs(count)
			d := suite.factory()

			d.AddListenerFuncs(toAdd...)
			suite.resetTestListeners(tests)
			d.Send(suite.testEvent)
			suite.assertTestListenersCalled(tests)

			d.Clear()
			suite.resetTestListeners(tests)
			d.Send(suite.testEvent)
			suite.assertTestListenersNotCalled(tests)
		})
	}
}

func (suite *DispatcherTestSuite[E, D]) TestAddListenerChans() {
	for _, count := range []int{1, 2, 5} {
		suite.Run(fmt.Sprintf("count=%d", count), func() {
			tests, toAdd := suite.newTestListenerChans(count)
			d := suite.factory()

			d.AddListenerChans(toAdd...)
			d.Send(suite.testEvent)
			suite.assertTestListenerChansCalled(tests)

			d.Clear()
			d.Send(suite.testEvent)
			suite.assertTestListenerChansNotCalled(tests)
		})
	}
}
