package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	var (
		num int
		mu  sync.Mutex

		qps = 10
		wg  sync.WaitGroup
		N   = 10000
	)

	wg.Add(N)

	limiter := rate.NewLimiter(rate.Every(time.Second), qps)

	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			for limiter.Wait(context.TODO()) == nil {
				mu.Lock()
				num++
				mu.Unlock()
			}
		}(i)
	}

	time.Sleep(time.Second)
	mu.Lock()
	fmt.Println("num:", num)
	mu.Unlock()

	fmt.Println("burst:", limiter.Burst())

	fmt.Println("blocking...")
	donec := make(chan struct{})
	go func() {
		wg.Wait()
		close(donec)
	}()
	select {
	case <-donec:
		fmt.Println("Done!")
	case <-time.After(time.Second):
		fmt.Println("Timed out!")
	}
}

/*
num: 11
burst: 10
blocking...
Timed out!
*/
