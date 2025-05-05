// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type customListenerFunc func(int)
type customChan chan int
type customSendOnlyChan chan<- int

type ListenerTestSuite struct {
	suite.Suite

	expected int
}

func (suite *ListenerTestSuite) SetupTest() {
	suite.expected = 999
}

func (suite *ListenerTestSuite) SetupSubTest() {
	suite.expected = 999
}

func (suite *ListenerTestSuite) testAsListenerNil() {
	suite.Nil(
		AsListener[int, func(int)](nil),
	)

	suite.Nil(
		AsListener[int, chan int](nil),
	)

	suite.Nil(
		AsListener[int, chan<- int](nil),
	)
}

func (suite *ListenerTestSuite) testAsListenerFunc() {
	called := false
	l := AsListener[int](func(v int) {
		called = true
		suite.Equal(suite.expected, v)
	})

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)
	suite.True(called)

}

func (suite *ListenerTestSuite) testAsListenerCustomFunc() {
	called := false
	l := AsListener[int](customListenerFunc(func(v int) {
		called = true
		suite.Equal(suite.expected, v)
	}))

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)
	suite.True(called)

}

func (suite *ListenerTestSuite) testAsListenerChan() {
	ch := make(chan int, 1)
	l := AsListener[int](ch)

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)

	select {
	case actual := <-ch:
		suite.Equal(suite.expected, actual)

	default:
		suite.Fail("failed to receive event")
	}
}

func (suite *ListenerTestSuite) testAsListenerCustomChan() {
	ch := make(customChan, 1)
	l := AsListener[int](ch)

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)

	select {
	case actual := <-ch:
		suite.Equal(suite.expected, actual)

	default:
		suite.Fail("failed to receive event")
	}
}

func (suite *ListenerTestSuite) testAsListenerSendOnlyChan() {
	ch := make(chan int, 1)
	l := AsListener[int]((chan<- int)(ch))

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)

	select {
	case actual := <-ch:
		suite.Equal(suite.expected, actual)

	default:
		suite.Fail("failed to receive event")
	}
}

func (suite *ListenerTestSuite) testAsListenerCustomSendOnlyChan() {
	ch := make(chan int, 1)
	l := AsListener[int](customSendOnlyChan(ch))

	suite.Require().NotNil(l)
	l.OnEvent(suite.expected)

	select {
	case actual := <-ch:
		suite.Equal(suite.expected, actual)

	default:
		suite.Fail("failed to receive event")
	}
}
func (suite *ListenerTestSuite) TestAsListener() {
	suite.Run("Nil", suite.testAsListenerNil)
	suite.Run("Func", suite.testAsListenerFunc)
	suite.Run("CustomFunc", suite.testAsListenerCustomFunc)
	suite.Run("Chan", suite.testAsListenerChan)
	suite.Run("CustomChan", suite.testAsListenerCustomChan)
	suite.Run("SendOnlyChan", suite.testAsListenerSendOnlyChan)
	suite.Run("CustomSendOnlyChan", suite.testAsListenerCustomSendOnlyChan)
}

func TestListener(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}
