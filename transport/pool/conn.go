package pool

import (
	"log"
	"net"
	"sync"
	"time"
)

type Conn struct {
	conn        net.Conn
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
	// if p.unable {
	// 	if p.conn !=nil {
	// 		return
	// 	}
	// }
	err := p.conn.SetDeadline(time.Time{})
	if err != nil {
		log.Fatalf("SetDeadline error: %v\n", err)
	}
	return p.ch.Put(p)
}
