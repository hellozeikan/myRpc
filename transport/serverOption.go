package transport

import (
	"context"
	"time"
)

type ServerTraOptions struct {
	Address         string        // address，e.g: ip://127.0.0.1：8080
	Network         string        // network type
	Timeout         time.Duration // transport layer request timeout ，default: 2 min
	Handler         Handler       // handler
	KeepAlivePeriod time.Duration // keepalive period
}
type Handler interface {
	Handle(context.Context, []byte) ([]byte, error)
}

type ServerTransportOption func(*ServerTraOptions)

func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTraOptions) {
		o.Address = address
	}
}

func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTraOptions) {
		o.Network = network
	}
}

// WithServerTimeout returns a ServerTransportOption which sets the value for timeout
func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTraOptions) {
		o.Timeout = timeout
	}
}

// WithHandler returns a ServerTransportOption which sets the value for handler
func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTraOptions) {
		o.Handler = handler
	}
}

// WithKeepAlivePeriod returns a ServerTransportOption which sets the value for keepAlivePeriod
func WithKeepAlivePeriod(keepAlivePeriod time.Duration) ServerTransportOption {
	return func(o *ServerTraOptions) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}
