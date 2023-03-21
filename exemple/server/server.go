package main

import (
	"context"
	"myrpc/service"
	"time"
)

func main() {
	opts := []service.ServerOption{
		service.WithAddress("127.0.0.1:8000"),
		service.WithTimeout(time.Millisecond * 20000),
	}
	s := service.NewServer(opts...)
	if err := s.RegisterService("helloworld.Greeter", &Service{}); err != nil {
		panic(err)
	}
	s.Serve()
}

type Service struct{}
type HelloRequest struct {
	Msg string
}

type HelloReply struct {
	Msg string
}
type Request struct {
	A int32 `msgpack:"a"`
	B int32 `msgpack:"b"`
}
type Reply struct {
	Result int32 `msgpack:"result"`
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloReply, error) {
	rsp := &HelloReply{
		Msg: "world",
	}

	return rsp, nil
}
func (s *Service) Add(ctx context.Context, req *Request) (*Reply, error) {
	rsp := &Reply{
		Result: req.A + req.B,
	}
	return rsp, nil
}
