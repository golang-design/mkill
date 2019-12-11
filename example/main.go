package main

import (
	"time"

	"github.com/changkun/mkill"
)

func main() {
	mkill.GOMAXTHREADS(10)
	for {
		time.Sleep(time.Second)
		go func() {
			time.Sleep(time.Second * 10)
		}()
	}
}
