// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

// +build darwin

package mkill

import "fmt"

var cmdThreads = fmt.Sprintf("ps M %d | wc -l", pid)
