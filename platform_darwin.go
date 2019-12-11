// +build darwin

package mkill

import "fmt"

var cmdThreads = fmt.Sprintf("ps M %d | wc -l", pid)
