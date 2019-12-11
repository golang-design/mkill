// +build linux

package mkill

import "fmt"

var cmdThreads = fmt.Sprintf("ps hH p %d | wc -l", pid)
