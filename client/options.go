package client

import (
	"time"

	"github.com/hellozeikan/myrpc/interceptor"
)

type Options struct {
	serviceName string
	method      string
	target      string
	timeout     time.Duration
	// 服务发现
	// 拦截器设置
	interceptors []interceptor.ClientInterceptor
}

type Option func(*Options)

func WithServiceName(serviceName string) Option {
	return func(o *Options) {
		o.serviceName = serviceName
	}
}

func WithMethod(method string) Option {
	return func(o *Options) {
		o.method = method
	}
}

func WithTarget(target string) Option {
	return func(o *Options) {
		o.target = target
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}

func WithInterceptor(interceptors ...interceptor.ClientInterceptor) Option {
	return func(o *Options) {
		o.interceptors = append(o.interceptors, interceptors...)
	}
}

// func WithSelectorName(selectorName string) Option {
// 	return func(o *Options) {
// 		o.selectorName = selectorName
// 	}
// }
