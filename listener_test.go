// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ListenerTestSuite struct {
	suite.Suite
}

func (suite *ListenerTestSuite) TestAsListener() {
	called := false
	l := AsListener(func(v int) {
		called = true
		suite.Equal(123, v)
	})

	suite.Require().NotNil(l)
	l.OnEvent(123)
	suite.True(called)
}

func (suite *ListenerTestSuite) TestListenerChan() {
	ch := make(chan int, 1)
	l := ListenerChan[int](ch)

	l.OnEvent(123)
	v, ok := <-ch
	suite.True(ok)
	suite.Equal(123, v)
}

func TestListener(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}
