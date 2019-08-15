package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	n := 400000
	m := 10
	wg := &sync.WaitGroup{}
	wg.Add(n * (m + 1))
	for i := 0; i < n; i++ {
		go func() {
			wg.Done()
			for j := 0; j < m; j++ {
				go func() {
					wg.Done()
					select {}
				}()
			}
			select {}
		}()
	}
	wg.Wait()
	a := runtime.NumGoroutine()
	fmt.Println(a)
	time.Sleep(time.Second * 3)
}
