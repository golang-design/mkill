// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

// +build linux

package mkill

// psFlags lists the threads of a process. "H" prints one line per thread and
// "p" selects the pid; the output is preceded by a single column header line.
// See parseNumM.
var psFlags = []string{"H", "p"}
