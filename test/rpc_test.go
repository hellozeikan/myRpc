package main

import (
	"context"
	"myrpc/client"
	"testing"
	"time"
)

func TestRPC(t *testing.T) {
	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
	}
	c := client.DefaultClient

	numRequests := 1000000
	done := make(chan bool)
	for i := 0; i < numRequests; i++ {
		go func() {
			req := &Request{
				A: 1111,
				B: 2222,
			}
			rsp := &Response{}
			_ = c.Call(context.Background(), "helloworld.Greeter.Add", req, rsp, opts...)
			// fmt.Printf("%+v\n", rsp)
			done <- true
		}()
	}

	for i := 0; i < numRequests; i++ {
		<-done
	}
	// rsp1 := &HelloReply{}
	// err = c.Call(context.Background(), "helloworld.Greeter.SayHello", req, rsp1, opts...)
	// fmt.Printf("%+v\n", rsp1)
	// fmt.Println(err)
}

type Request struct {
	A int `msgpack:"a"`
	B int `msgpack:"b"`
}
type Response struct {
	Result int `mapstructure:"result" msgpack:"result"`
}
type HelloReply struct {
	Msg string
}
