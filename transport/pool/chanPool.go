package pool

import (
	"errors"
	"net"
	"sync"
	"time"
)

type chanPool struct {
	mu    sync.Mutex
	conns chan *Conn
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
		return conn.conn.Close()
	}
}

func isConnAlive(conn net.Conn) bool {
	conn.SetReadDeadline(time.Now().Add(time.Millisecond))

	return true
}
