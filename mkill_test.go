// Copyright 2020 The golang.design Initiative authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.

package mkill_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"golang.design/x/mkill"
)

func TestMKill(t *testing.T) {
	mkill.GOMAXTHREADS(10)

	// create a lot of threads by sleep gs
	for i := 0; i < 100000; i++ {
		go func() {
			time.Sleep(time.Second * 10)
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*100)
	defer cancel()
	ok := mkill.Wait(ctx)
	if !ok {
		t.Fatal("mkill failed in 100s")
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
