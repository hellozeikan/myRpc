# myRpc

## 通信作为互联网最重要且基础的组件，也是使用场景最为丰富的，在任何交互的地方都需要使用；在了解整个计算机网络体系后，想先从最上层的应用实现一个rpc组件

## rpc所使用的协议暂时使用msgback
## 以下是各各序列化协议的对比
```
JSON:可读性好、简单易用。
Protobuf:二进制消息，效率高，性能好。
Thrift:省流量，体积较小 。
MessagePack:序列化反序列化效率高，文件体积小，比json小一倍。
第二种和第三种已有成熟的框架，不必再造轮子
```

## 需要完成
> 1.server端
>> 序列化组件,异常处理

> 2.client端
>> 连接池