// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSyncListeners(t *testing.T) {
	suite.Run(t, &DispatcherTestSuite[int, *SyncListeners[int]]{
		factory:   func() *SyncListeners[int] { return new(SyncListeners[int]) },
		testEvent: 123,
	})
}
