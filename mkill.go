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

	// threadCount reports the current number of OS threads. It is a package
	// variable so tests can replace it with a deterministic stub.
	threadCount = NumM
)

// psArgs returns the arguments passed to ps to list this process' threads.
func psArgs() []string {
	args := make([]string, 0, len(psFlags)+1)
	args = append(args, psFlags...)
	return append(args, strconv.Itoa(pid))
}

// runPS lists the process' threads via ps. It is a package variable so tests
// can stub the subprocess out.
var runPS = func() ([]byte, error) {
	return exec.Command("ps", psArgs()...).Output()
}

// NumM returns the number of running threads. It returns 0 if the count cannot
// be determined.
func NumM() int {
	return parseNumM(runPS())
}

// parseNumM turns the output of ps into a thread count. ps prints one column
// header line followed by one line per thread, so the count is the number of
// lines after the header.
func parseNumM(out []byte, err error) int {
	if err != nil {
		if debug {
			fmt.Printf("mkill: failed to fetch #threads: %v\n", err)
		}
		return 0
	}
	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return 0
	}
	n := len(strings.Split(trimmed, "\n")) - 1 // drop the header line
	if n < 0 {
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
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		if threadCount() <= GOMAXTHREADS(0) {
			return true
		}
		select {
		case <-ctx.Done():
			return threadCount() <= GOMAXTHREADS(0)
		case <-t.C:
		}
	}
}

// enabled reports whether the environment supports counting threads. When it
// does not (for example ps is missing from PATH), the package degrades to a
// no-op rather than crashing the program that merely imported it.
func enabled() bool {
	_, err := runPS()
	return err == nil
}

func init() {
	if debug {
		fmt.Printf("mkill: pid %v, maxThread %v, interval %v\n", pid, maxThreads, interval)
	}
	if !enabled() {
		if debug {
			fmt.Println("mkill: disabled, cannot count threads")
		}
		return
	}
	go reap()
}

// reap periodically counts the threads and, whenever there are more than the
// configured maximum, retires the surplus by spawning goroutines that lock an
// OS thread and then exit: the runtime terminates the thread of a locked
// goroutine that returns without unlocking.
func reap() {
	var wg sync.WaitGroup
	t := time.NewTicker(interval)
	defer t.Stop()
	for range t.C {
		n := threadCount()
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
