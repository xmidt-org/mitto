// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSet(t *testing.T) {
	suite.Run(t, &DispatcherTestSuite[int, *Set[int]]{
		factory:   func() *Set[int] { return new(Set[int]) },
		testEvent: 123,
	})
}
