package server

import (
	"context"
	"time"
)

type ServerTraOption struct {
	Address         string        // address，e.g: ip://127.0.0.1：8080
	Network         string        // network type
	Timeout         time.Duration // transport layer request timeout ，default: 2 min
	Handler         Handler       // handler
	KeepAlivePeriod time.Duration // keepalive period
}
type Handler interface {
	Handle(context.Context, []byte) ([]byte, error)
}

type ServerTransportOption func(*ServerTraOption)

func WithServerAddress(address string) ServerTransportOption {
	return func(o *ServerTraOption) {
		o.Address = address
	}
}

func WithServerNetwork(network string) ServerTransportOption {
	return func(o *ServerTraOption) {
		o.Network = network
	}
}

// WithServerTimeout returns a ServerTransportOption which sets the value for timeout
func WithServerTimeout(timeout time.Duration) ServerTransportOption {
	return func(o *ServerTraOption) {
		o.Timeout = timeout
	}
}

// WithHandler returns a ServerTransportOption which sets the value for handler
func WithHandler(handler Handler) ServerTransportOption {
	return func(o *ServerTraOption) {
		o.Handler = handler
	}
}

// WithKeepAlivePeriod returns a ServerTransportOption which sets the value for keepAlivePeriod
func WithKeepAlivePeriod(keepAlivePeriod time.Duration) ServerTransportOption {
	return func(o *ServerTraOption) {
		o.KeepAlivePeriod = keepAlivePeriod
	}
}
