package transport

import (
	"context"
	"encoding/binary"
	"io"
	"net"

	"github.com/hellozeikan/myrpc/code"
)

const DefaultPayloadLength = 1024
const MaxPayloadLength = 4 * 1024 * 1024

type ServerTransport interface {
	// 监听和处理请求
	ListenAndServe(context.Context, ...ServerTransportOption) error
}
type ClientTransport interface {
	// 发送请求
	Send(context.Context, []byte, ...ClientTransportOption) ([]byte, error)
}

type Framer interface {
	ReadFrame(net.Conn) ([]byte, error)
}

type framer struct {
	buffer  []byte
	counter int // to prevent the dead loop
}

// ReadFrame implements Framer
func (f *framer) ReadFrame(conn net.Conn) ([]byte, error) {

	frameHeader := make([]byte, code.FrameHeadLen)
	if num, err := io.ReadFull(conn, frameHeader); num != code.FrameHeadLen || err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(frameHeader) // 目前header里只有length

	if length > MaxPayloadLength {
		return nil, code.NewFrameworkError(code.ClientMsgErrorCode, "payload too large...")
	}

	f.buffer = make([]byte, tableSizeFor(int(length)))

	if num, err := io.ReadFull(conn, f.buffer[:length]); uint32(num) != length || err != nil {
		return nil, err
	}

	return append(frameHeader, f.buffer[:length]...), nil
}

func NewFramer() Framer {
	return &framer{
		buffer: make([]byte, DefaultPayloadLength),
	}
}
func (f *framer) Resize() {
	f.buffer = make([]byte, len(f.buffer)*2)
}

// // TCP/IP协议RFC1700规定使用“大端”字节序为网络字节序，开发的时候需要遵守这一规则
// func read(conn *net.TCPConn) ([]byte, error) {
// 	maxPayloadLength := uint32(1024 * 1024)
// 	frameHeader := make([]byte, 4)
// 	if num, err := io.ReadFull(conn, frameHeader); num != 4 || err != nil {
// 		return nil, err
// 	}
// 	length := binary.BigEndian.Uint32(frameHeader)
// 	if length > maxPayloadLength {
// 		return nil, nil //error("payload too large...")
// 	}
// 	buffer := make([]byte, tableSizeFor(int(length)))
// 	if num, err := io.ReadFull(conn, buffer[:length]); uint32(num) != length || err != nil {
// 		return nil, err
// 	}
// 	return append(frameHeader, buffer[:length]...), nil
// }

// 通过位移31次（1+2+4+8+16）以及或运算，把当前最高位下面的二进制都填满1，这样再加1以后就能得到比原先高一位的数字
func tableSizeFor(source int) int {
	maxCapacity := 1 << 30
	n := source - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	if n < 0 {
		return 1
	} else if n >= maxCapacity {
		return maxCapacity
	}
	return n + 1
}
