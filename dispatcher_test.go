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

func (suite *DispatcherTestSuite[E, D]) newTestListeners(n int) (tests []*testListener[E]) {
	tests = make([]*testListener[E], n)

	for i := 0; i < n; i++ {
		tests[i] = suite.newTestListener()
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

func (suite *DispatcherTestSuite[E, D]) TestEmpty() {
	d := suite.factory()

	// all of the following should be nops
	d.Clear()
	d.AddListeners()
	d.RemoveListeners()
	d.Send(suite.testEvent)
}

func (suite *DispatcherTestSuite[E, D]) testAddListenersEmpty() {
	var tests []*testListener[E]
	d := suite.factory()

	AddListeners(d, tests...) // should add nothing
	d.Send(suite.testEvent)

	RemoveListeners(d, tests...)
	d.Send(suite.testEvent)
}

func (suite *DispatcherTestSuite[E, D]) testAddListenersLifecycle(count int) {
	tests := suite.newTestListeners(count)
	d := suite.factory()

	AddListeners(d, tests...)
	suite.resetTestListeners(tests)
	d.Send(suite.testEvent)
	suite.assertTestListenersCalled(tests)

	RemoveListeners(d, tests...)
	suite.resetTestListeners(tests)
	d.Send(suite.testEvent)
	suite.assertTestListenersNotCalled(tests)

	AddListeners(d, tests...)
	d.Clear()
	suite.resetTestListeners(tests)
	d.Send(suite.testEvent)
	suite.assertTestListenersNotCalled(tests)

	// check that nils are skipped
	AddListeners[E, Listener[E]](d, nil, nil)
	AddListeners(d, tests...)
	AddListeners[E, Listener[E]](d, nil, nil)
	suite.resetTestListeners(tests)
	d.Send(suite.testEvent)
	suite.assertTestListenersCalled(tests)
}

func (suite *DispatcherTestSuite[E, D]) tesetAddListenersRemoveSinks() {
	var (
		f = func(E) {
			suite.Fail("closure should not have received an event")
		}

		ch1 = make(chan E, 1)
		ch2 = make(chan E, 1)

		toAdd = []Listener[E]{
			AsListener[E](f),
			AsListener[E](ch1),
			AsListener[E]((chan<- E)(ch2)),
		}

		d = suite.factory()
	)

	d.AddListeners(toAdd...)
	d.RemoveListeners(toAdd...)
	d.Send(suite.testEvent)

	select {
	case <-ch1:
		suite.Fail("should not have received an event on a channel")

	default:
		// passing
	}

	select {
	case <-ch2:
		suite.Fail("should not have received an event on a send-only channel")

	default:
		// passing
	}
}

func (suite *DispatcherTestSuite[E, D]) TestAddListeners() {
	suite.Run("Empty", suite.testAddListenersEmpty)

	for _, count := range []int{1, 2, 5} {
		suite.Run(fmt.Sprintf("count=%d", count), func() {
			suite.testAddListenersLifecycle(count)
		})
	}

	suite.Run("RemoveSinks", suite.tesetAddListenersRemoveSinks)
}
