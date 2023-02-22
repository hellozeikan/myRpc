package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func(i int) {
			fmt.Println(i)
			time.Sleep(30 * time.Millisecond)
			wg.Done()
		}(i)
	}

	wg.Wait()
}
