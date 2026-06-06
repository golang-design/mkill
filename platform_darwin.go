// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

// +build darwin

package mkill

// psFlags lists the threads of a process. "M" prints one line per thread,
// preceded by a single column header line. See parseNumM.
var psFlags = []string{"M"}
