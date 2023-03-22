package main

import (
	"context"
	"fmt"
	"myrpc/client"
	"time"
)

func main() {
	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithTimeout(20000 * time.Millisecond),
	}
	c := client.DefaultClient

	req := &Request{
		A: 1111,
		B: 2222,
	}
	var err error
	rsp := &Response{}
	err = c.Call(context.Background(), "helloworld.Greeter.Add", req, rsp, opts...)
	fmt.Printf("%+v\n", rsp)
	fmt.Println(err)
	rsp1 := &HelloReply{}
	err = c.Call(context.Background(), "helloworld.Greeter.SayHello", req, rsp1, opts...)
	fmt.Printf("%+v\n", rsp1)
	fmt.Println(err)
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
