package server

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"time"
)

type ServerOption struct {
	Addr    string // e.g:127.0.0.1:8000
	NetWork string // tcp/udp
	TimeOut time.Duration
}

func ListenAndServe(option ServerOption) error {
	ls, err := net.Listen(option.NetWork, option.Addr)
	if err != nil {
		return err
	}
	go func() {
		if err = serve(ls); err != nil {
			log.Fatalf("transport serve error, %v", err)
		}
	}()
	log.Fatalf("server listening on %s\n", option.Addr)
	return nil
}

func serve(ls net.Listener) error {
	var timeOut time.Duration
	tl, _ := ls.(*net.TCPListener)
	// tl, ok := lis.(*net.TCPListener)
	// if !ok {
	// 	return codes.NetworkNotSupportedError
	// }
	for {
		conn, err := tl.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if timeOut == 0 {
					timeOut = 5 * time.Millisecond
				} else {
					timeOut *= 2
				}
				if max := 1 * time.Second; timeOut > max {
					timeOut = max
				}
				time.Sleep(timeOut)
				continue
			}
			return err
		}

		if err = conn.SetKeepAlive(true); err != nil {
			return err
		}
		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Fatalf("panic: %v", err)
				}
			}()

			if err := handleConn(conn); err != nil {
				log.Fatalf("mrpc handle tcp conn error, %v", err)
			}
		}()
	}
}

func handleConn(conn *net.TCPConn) error {
	defer conn.Close()
	for {
		byte, err := read(conn)
		if err != nil {
			return err
		}
		rsp, err := handle(byte)
		if err != nil {
			log.Fatalf("s.handle err is not nil, %v", err)
		}
		if err = write(conn, rsp); err != nil {
			return err
		}
	}
}

// 处理请求
func handle(req []byte) ([]byte, error) {
	return req, nil
}

func write(conn *net.TCPConn, rsp []byte) error {
	if _, err := conn.Write(rsp); err != nil {
		log.Fatalf("conn Write err: %v", err)
	}
	return nil
}

// TCP/IP协议RFC1700规定使用“大端”字节序为网络字节序，开发的时候需要遵守这一规则
func read(conn *net.TCPConn) ([]byte, error) {
	maxPayloadLength := uint32(1024 * 1024)
	frameHeader := make([]byte, 4)
	if num, err := io.ReadFull(conn, frameHeader); num != 4 || err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(frameHeader)
	if length > maxPayloadLength {
		return nil, nil //error("payload too large...")
	}
	buffer := make([]byte, tableSizeFor(int(length)))
	if num, err := io.ReadFull(conn, buffer[:length]); uint32(num) != length || err != nil {
		return nil, err
	}
	return append(frameHeader, buffer[:length]...), nil
}

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
