// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

// Package mkill limits the number of threads in a Go program, without crash the whole program.
package mkill // import "golang.design/x/mkill"

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	pid        = os.Getpid()
	minThreads = int32(runtime.NumCPU()) + 2 // minimum number of threads required by the runtime
	maxThreads = int32(runtime.NumCPU()) + 2 // 2 meaning runtime sysmon thread + template thread
	interval   = time.Second
	debug      = false
)

// NumM returns the number of running threads.
func NumM() int {
	out, err := exec.Command("bash", "-c", cmdThreads).Output()
	if err != nil && debug {
		fmt.Printf("mkill: failed to fetch #threads: %v\n", err)
		return 0
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil && debug {
		fmt.Printf("mkill: failed to parse #threads: %v\n", err)
		return 0
	}
	return n
}

// GOMAXTHREADS sets the maximum number of system threads that allowed in a Go program
// and returns the previous setting. If n is lower than minimum required number of threads,
// it does not change the current setting.
// The minimum allowed number of threads of a program is runtime.NumCPU() + 2.
func GOMAXTHREADS(n int) int {
	if n < int(minThreads) {
		return int(atomic.LoadInt32(&maxThreads))
	}

	return int(atomic.SwapInt32(&maxThreads, int32(n)))
}

// Wait waits until the number of threads meet the GOMAXTHREADS settings.
// The function always returns true if the ctx is not canceled.
// Otherwise returns true only if the Wait is successed in the last check.
func Wait(ctx context.Context) (ok bool) {
	for {
		select {
		case <-ctx.Done():
			if NumM() <= GOMAXTHREADS(0) {
				ok = true
			}
			return
		default:
			if NumM() > GOMAXTHREADS(0) {
				continue
			}
			ok = true
			return
		}
	}
}

func checkwork() {
	_, err := exec.Command("bash", "-c", cmdThreads).Output()
	if err != nil {
		panic(fmt.Sprintf("mkill: failed to use the package: %v", err))
	}
}

func init() {
	checkwork()
	if debug {
		fmt.Printf("mkill: pid %v, maxThread %v, interval %v\n", pid, maxThreads, interval)
	}

	wg := sync.WaitGroup{}
	go func() {
		t := time.NewTicker(interval)
		for {
			select {
			case <-t.C:
				n := NumM()
				nkill := int32(n) - atomic.LoadInt32(&maxThreads)
				if nkill <= 0 {
					if debug {
						fmt.Printf("mkill: checked #threads total %v / max %v\n", n, maxThreads)
					}
					continue
				}
				wg.Add(int(nkill))
				for i := int32(0); i < nkill; i++ {
					go func() {
						runtime.LockOSThread()
						wg.Done()
					}()
				}
				wg.Wait()
				if debug {
					fmt.Printf("mkill: killing #threads, remaining: %v\n", n)
				}
			}
		}
	}()
}
