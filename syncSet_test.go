// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mitto

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSyncSet(t *testing.T) {
	suite.Run(t, &DispatcherTestSuite[int, *SyncSet[int]]{
		factory:   func() *SyncSet[int] { return new(SyncSet[int]) },
		testEvent: 123,
	})
}
