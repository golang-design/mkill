// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill_test

import (
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"golang.design/x/mkill"
)

// TestNumMMatchesOracle checks NumM against an independent, kernel-backed
// thread count rather than against synthetic ps output. It guards the darwin
// header off-by-one fix: any constant skew in NumM (e.g. counting the ps
// header) would make it disagree with the oracle.
//
// The only portable oracle without cgo is Linux's /proc/self/status, so the
// test runs there (which is also where CI runs) and skips elsewhere. The
// manual darwin verification lives in the commit history.
func TestNumMMatchesOracle(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skipf("no /proc thread oracle on %s", runtime.GOOS)
	}

	// Thread count can drift as the runtime parks/spawns Ms. Read the oracle
	// on both sides of NumM and only assert when the window was stable, so the
	// comparison is exact without being flaky.
	for i := 0; i < 50; i++ {
		before := procThreads(t)
		n := mkill.NumM()
		after := procThreads(t)
		if before != after {
			continue
		}
		if n != before {
			t.Fatalf("NumM() = %d, /proc oracle = %d", n, before)
		}
		return
	}
	t.Skip("thread count never stabilized; skipping to avoid flakiness")
}

// procThreads reads the live thread count of this process from /proc.
func procThreads(t *testing.T) int {
	t.Helper()
	data, err := ioutil.ReadFile("/proc/self/status")
	if err != nil {
		t.Fatalf("read /proc/self/status: %v", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "Threads:") {
			n, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Threads:")))
			if err != nil {
				t.Fatalf("parse Threads line %q: %v", line, err)
			}
			return n
		}
	}
	t.Fatal("no Threads field in /proc/self/status")
	return 0
}
