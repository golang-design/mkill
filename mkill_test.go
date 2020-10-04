// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill_test

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"golang.design/x/mkill"
)

func TestMKill(t *testing.T) {
	mkill.GOMAXTHREADS(10)

	// create a lot of threads by sleep gs
	wg := sync.WaitGroup{}
	wg.Add(100000)
	for i := 0; i < 100000; i++ {
		go func() {
			time.Sleep(time.Second * 1)
			wg.Done()
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	ok := mkill.Wait(ctx)
	if !ok {
		t.Fatal("mkill failed in 10s")
	}
	wg.Wait()
}

func TestMinThreads(t *testing.T) {
	old := mkill.GOMAXTHREADS(0)
	n := runtime.NumCPU()
	if mkill.GOMAXTHREADS(n-1) != old {
		t.Fatalf("number of threads is less than required in the runtime")
	}
}

func ExampleGOMAXTHREADS() {
	mkill.GOMAXTHREADS(10)
	// Output:
}

func ExampleNumM() {
	mkill.GOMAXTHREADS(10)
	mkill.Wait(context.Background())
	fmt.Println(mkill.NumM() <= 10)
	// Output:
	// true
}

func ExampleWait() {
	mkill.GOMAXTHREADS(10)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	fmt.Println(mkill.Wait(ctx))
	// Output:
	// true
}
