// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill

import (
	"errors"
	"testing"
)

// TestEnabledDegradesWhenPSUnavailable guards the behavior that replaced the
// init-time panic: when threads cannot be counted, the package reports itself
// disabled instead of crashing the importing program.
func TestEnabledDegradesWhenPSUnavailable(t *testing.T) {
	old := runPS
	defer func() { runPS = old }()

	runPS = func() ([]byte, error) { return nil, errors.New("no ps in PATH") }
	if enabled() {
		t.Fatal("enabled() = true when ps is unavailable")
	}

	runPS = func() ([]byte, error) { return []byte("USER PID TT\n a"), nil }
	if !enabled() {
		t.Fatal("enabled() = false when ps works")
	}
}
