package pool

import (
	"context"
	"errors"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type chanPool struct {
	mu          sync.Mutex
	conns       chan *Conn
	idleTimeout time.Duration
	dialTimeout time.Duration
	Dial        func(context.Context) (net.Conn, error)
	initCap     int
	maxCap      int
	maxIdle     int
	inflight    int32
}

func (ch *chanPool) Put(conn *Conn) error {
	if conn == nil {
		return errors.New("connection closed")
	}
	ch.mu.Lock()
	defer ch.mu.Unlock()
	if ch.conns == nil {
		conn.Makeunable()
		conn.Close() // 不存在管道直接关闭
	}
	select {
	case ch.conns <- conn:
		return nil
	default:
		return conn.Conn.Close()
	}
}

func (ch *chanPool) Get(ctx context.Context) (*Conn, error) {
	if ch.conns == nil {
		return nil, errors.New("connection closed")
	}
	select {
	case conn := <-ch.conns:
		if conn == nil {
			return nil, errors.New("connection closed")
		}
		if conn.unable {
			return nil, errors.New("connection closed")
		}
		return conn, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		if ch.inflight > int32(ch.maxCap) {
			select {
			case conn := <-ch.conns:
				if conn == nil {
					return nil, errors.New("connection closed")
				}

				if conn.unable {
					return nil, errors.New("connection closed") // 这里出错了没有自动重试，调用方根据错误类型来决定是否重试
				}

				return conn, nil
			case <-ctx.Done(): // context取消或超时，则退出
				return nil, ctx.Err()
			}
		}
		conn, err := ch.Dial(ctx)
		if err != nil {
			return nil, err
		}
		atomic.AddInt32(&ch.inflight, 1)
		return ch.wrapConn(conn), nil
	}
}

func (ch *chanPool) wrapConn(conn net.Conn) *Conn {
	p := &Conn{
		ch:          ch,
		t:           time.Now(),
		dialTimeout: ch.dialTimeout,
	}
	return p
}

func (ch *chanPool) Check(cn *Conn) bool {
	// check timeout
	if cn.t.Add(ch.idleTimeout).Before(time.Now()) {
		return false
	}
	if !isConnAlive(cn.Conn) {
		return false
	}
	return true
}

func isConnAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond))
	buff := make([]byte, 1)
	if n, err := conn.Read(buff); n > 0 || err == io.EOF {
		return false
	}
	_ = conn.SetReadDeadline(time.Time{})

	return true
}
