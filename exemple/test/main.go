package main

import (
	"fmt"

	"github.com/goinggo/mapstructure"
)

func main() {
	// var wg sync.WaitGroup
	// wg.Add(1000)
	// for i := 0; i < 1000; i++ {
	// 	go func(i int) {
	// 		fmt.Println(i)
	// 		time.Sleep(30 * time.Millisecond)
	// 		wg.Done()
	// 	}(i)
	// }

	// wg.Wait()
	rsp := &Response{}
	m := make(map[string]int)
	m["result"] = 123
	mapstructure.Decode(m, rsp)
	fmt.Println(rsp)
}

type Response struct {
	Result int `mapstructure:"result"`
}
type HelloReply struct {
	Msg string
}
