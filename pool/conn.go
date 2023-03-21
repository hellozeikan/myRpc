package pool

import (
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

var ErrConnClosed = errors.New("connection closed")

type Conn struct {
	net.Conn
	ch          *chanPool
	unable      bool
	mu          sync.Mutex
	t           time.Time
	dialTimeout time.Duration // connection timeout duration
}

func (p *Conn) Makeunable() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.unable = true
}

func (p *Conn) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.unable {
		if p.Conn != nil {
			return p.Conn.Close()
		}
	}
	err := p.Conn.SetDeadline(time.Time{})
	if err != nil {
		log.Printf("SetDeadline error: %v\n", err)
	}
	return p.ch.Put(p)
}
func (p *Conn) Read(b []byte) (int, error) {
	// 判断该连接状态
	if p.unable {
		return 0, errors.New("connection closed")
	}
	n, err := p.Conn.Read(b)
	if err != nil {
		p.Makeunable()
		// 关闭连接
		p.Conn.Close()
	}
	return n, err
}

func (p *Conn) Write(b []byte) (int, error) {
	if p.unable {
		return 0, errors.New("connection closed")
	}
	n, err := p.Conn.Write(b)
	if err != nil {
		p.Makeunable()
		p.Conn.Close()
	}
	return n, err
}
