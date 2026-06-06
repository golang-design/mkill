// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// TestWaitDoesNotBusySpin guards against the regression where Wait looped in a
// tight default branch, hammering the (expensive) thread counter as fast as the
// CPU allowed. With the fix, Wait polls at most once per interval.
func TestWaitDoesNotBusySpin(t *testing.T) {
	var calls int32
	oldCount, oldInterval := threadCount, interval
	defer func() { threadCount, interval = oldCount, oldInterval }()

	interval = 10 * time.Millisecond
	threadCount = func() int {
		atomic.AddInt32(&calls, 1)
		return 1 << 30 // always above any allowed maximum
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	Wait(ctx)

	// Polling every 10ms for ~100ms is a handful of calls. A busy-spin produces
	// many thousands.
	if n := atomic.LoadInt32(&calls); n > 50 {
		t.Fatalf("Wait appears to busy-spin: %d thread-count calls in 100ms", n)
	}
}
