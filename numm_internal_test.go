// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill

import (
	"errors"
	"testing"
)

func TestParseNumM(t *testing.T) {
	for _, tt := range []struct {
		name string
		out  string
		err  error
		want int
	}{
		// ps prints a header line followed by one line per thread, so the
		// thread count must exclude the header. This is the regression that
		// caused darwin (which kept the header) to overcount by one.
		{"three threads", "USER PID TT\n a\n b\n c", nil, 3},
		{"single thread", "USER PID TT\n a", nil, 1},
		{"trailing newline", "USER PID TT\n a\n b\n", nil, 2},
		{"header only", "USER PID TT", nil, 0},
		{"empty", "", nil, 0},
		{"error", "whatever", errors.New("boom"), 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNumM([]byte(tt.out), tt.err); got != tt.want {
				t.Fatalf("parseNumM(%q, %v) = %d, want %d", tt.out, tt.err, got, tt.want)
			}
		})
	}
}
