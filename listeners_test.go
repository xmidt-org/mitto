// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestListeners(t *testing.T) {
	suite.Run(t, &DispatcherTestSuite[int, *Listeners[int]]{
		factory:   func() *Listeners[int] { return new(Listeners[int]) },
		testEvent: 123,
	})
}
