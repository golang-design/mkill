package mkill

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var (
	pid       = os.Getpid()
	maxThread = int32(runtime.NumCPU())
	interval  = time.Second
	debug     = false
)

func checkwork() {
	_, err := getThreads()
	if err != nil {
		panic(fmt.Sprintf("mkill: failed to use the library: %v", err))
	}
}

func init() {
	checkwork()

	if debug {
		fmt.Printf("mkill: pid %v, maxThread %v, interval %v\n", pid, maxThread, interval)
	}
	go func() {
		t := time.NewTicker(interval)
		for {
			select {
			case <-t.C:
				n, _ := getThreads()
				nkill := int32(n) - atomic.LoadInt32(&maxThread)
				if nkill <= 0 {
					if debug {
						fmt.Printf("mkill: checked #threads total %v / max %v\n", n, maxThread)
					}
					continue
				}
				for i := int32(0); i < nkill; i++ {
					go func() {
						runtime.LockOSThread()
					}()
				}
				if debug {
					fmt.Printf("mkill: killing #threads, remaining: %v\n", n)
				}
			}
		}
	}()
}

// GOMAXTHREADS change the limits of the maximum threads in runtime
// and returns the previous number of threads limit
func GOMAXTHREADS(n int) int {
	return int(atomic.SwapInt32(&maxThread, int32(n)))
}

// getThreads returns the number of running threads
// Linux:
func getThreads() (int, error) {
	out, err := exec.Command("bash", "-c", cmdThreads).Output()
	if err != nil {
		return 0, fmt.Errorf("mkill: failed to fetch #threads: %v", err)
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return 0, fmt.Errorf("mkill: failed to parse #threads: %v", err)
	}
	return n, nil
}
