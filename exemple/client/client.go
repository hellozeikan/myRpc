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
	rsp1 := &Response{}
	err = c.Call(context.Background(), "helloworld.Greeter.Add", req, rsp1, opts...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", rsp1)
	rsp2 := &HelloReply{}
	err = c.Call(context.Background(), "helloworld.Greeter.SayHello", req, rsp2, opts...)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", rsp2)
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
