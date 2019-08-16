// http://www.sarathlakshman.com/2016/06/15/pitfall-of-golang-scheduler
package main

import (
	"runtime"
	"time"
)

func main() {
	processors := runtime.GOMAXPROCS(0)
	for i := 0; i < processors; i++ {
		go func() {
			for {
			}
		}()
	}
	time.Sleep(time.Second)
}
