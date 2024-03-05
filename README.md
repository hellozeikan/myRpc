# myRpc

## 通信作为互联网最重要且基础的组件，也是使用场景最为丰富的，在任何交互的地方都需要使用；在了解整个计算机网络体系后，想先从最上层的应用实现一个rpc组件

## rpc所使用的协议暂时使用msgback
```
使用说明：
go get -u  github.com/hellozeikan/myrpc@v1.1.0
go mod tidy

```
## server
```
package main

import (
	"context"
	"time"

	"github.com/hellozeikan/myrpc/service"
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

```
## client:
```
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hellozeikan/myrpc/client"
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

```

## 需要改善
- [ ] errors
- [ ] 服务注册与发现
- [ ] 性能优化
- [ ] 超时控制完善
