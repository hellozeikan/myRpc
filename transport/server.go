package transport

import (
	"context"
	"fmt"
	"io"
	"log"
	"myrpc/code"
	"net"
	"strconv"
	"time"
)

type serverTransport struct {
	opts *ServerTraOptions
}

var serverTransportMap = make(map[string]ServerTransport)

func init() {
	serverTransportMap["default"] = DefaultServerTransport
}

// RegisterServerTransport supports business custom registered ServerTransport
func RegisterServerTransport(name string, serverTransport ServerTransport) {
	if serverTransportMap == nil {
		serverTransportMap = make(map[string]ServerTransport)
	}
	serverTransportMap[name] = serverTransport
}

// Get the ServerTransport
func GetServerTransport(transport string) ServerTransport {
	if v, ok := serverTransportMap[transport]; ok {
		return v
	}

	return DefaultServerTransport
}

var DefaultServerTransport = NewServerTransport()

var NewServerTransport = func() ServerTransport {
	return &serverTransport{
		opts: &ServerTraOptions{},
	}
}

func (s *serverTransport) ListenAndServe(ctx context.Context, opts ...ServerTransportOption) error {
	for _, o := range opts {
		o(s.opts)
	}

	lis, err := net.Listen("tcp", s.opts.Address)
	if err != nil {
		return err
	}

	go func() {
		if err = s.serve(ctx, lis); err != nil {
			log.Fatalf("transport serve error, %v", err)
		}
	}()

	addr, err := Extract(s.opts.Address, lis)
	if err != nil {
		return err
	}
	log.Fatalf("server listening on %s\n", addr)

	return nil
}

func (s *serverTransport) serve(ctx context.Context, lis net.Listener) error {
	var tempDelay time.Duration

	tl, ok := lis.(*net.TCPListener)
	if !ok {
		return code.NetworkNotSupportedError
	}

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := tl.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		if err = conn.SetKeepAlive(true); err != nil {
			return err
		}

		if s.opts.KeepAlivePeriod != 0 {
			err := conn.SetKeepAlivePeriod(s.opts.KeepAlivePeriod)
			if err != nil {
				log.Fatalf("SetKeepAlivePeriod error: %v", err)
			}
		}

		go func() {
			defer func() {
				if err := recover(); err != nil {
					log.Fatalf("panic: %v", err)
				}
			}()

			if err := s.handleConn(ctx, wrapConn(conn)); err != nil {
				log.Fatalf("mrpc handle tcp conn error, %v", err)
			}
		}()
	}
}

func (s *serverTransport) handleConn(ctx context.Context, conn *connWrapper) error {
	defer conn.Close()

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		frame, err := s.read(ctx, conn)
		if err == io.EOF {
			// read compeleted
			return nil
		}

		if err != nil {
			return err
		}

		rsp, err := s.handle(ctx, frame)
		if err != nil {
			log.Fatalf("s.handle err is not nil, %v", err)
		}

		if err = s.write(ctx, conn, rsp); err != nil {
			return err
		}
	}

}

func (s *serverTransport) read(ctx context.Context, conn *connWrapper) ([]byte, error) {
	frame, err := conn.framer.ReadFrame(conn)

	if err != nil {
		return nil, err
	}

	return frame, nil
}

func (s *serverTransport) handle(ctx context.Context, frame []byte) ([]byte, error) {
	// parse reqbuf into req interface {}
	serverCodec := code.DefaultCodec

	reqbuf, err := serverCodec.Decode(frame)
	if err != nil {
		log.Fatalf("server Decode error: %v", err)
		return nil, err
	}

	rspbuf, err := s.opts.Handler.Handle(ctx, reqbuf)
	if err != nil {
		// todo: handle error
		log.Fatalf("server Handle error: %v", err)
	}

	rspbody, err := serverCodec.Encode(rspbuf)
	if err != nil {
		log.Fatalf("server Encode error, response: %v, err: %v", rspbuf, err)
		return nil, err
	}

	return rspbody, nil
}

func (s *serverTransport) write(ctx context.Context, conn net.Conn, rsp []byte) error {
	if _, err := conn.Write(rsp); err != nil {
		log.Fatalf("conn Write err: %v", err)
	}

	return nil
}

type connWrapper struct {
	net.Conn
	framer Framer
}

func wrapConn(rawConn net.Conn) *connWrapper {
	return &connWrapper{
		Conn:   rawConn,
		framer: NewFramer(),
	}
}

func isValidIP(addr string) bool {
	ip := net.ParseIP(addr)
	return ip.IsGlobalUnicast() && !ip.IsInterfaceLocalMulticast()
}

func Port(lis net.Listener) (int, bool) {
	if addr, ok := lis.Addr().(*net.TCPAddr); ok {
		return addr.Port, true
	}
	return 0, false
}

func Extract(hostPort string, lis net.Listener) (string, error) {
	addr, port, err := net.SplitHostPort(hostPort)
	if err != nil && lis == nil {
		return "", err
	}
	if lis != nil {
		p, ok := Port(lis)
		if !ok {
			return "", fmt.Errorf("failed to extract port: %v", lis.Addr())
		}
		port = strconv.Itoa(p)
	}
	if len(addr) > 0 && (addr != "0.0.0.0" && addr != "[::]" && addr != "::") {
		return net.JoinHostPort(addr, port), nil
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	minIndex := int(^uint(0) >> 1)
	ips := make([]net.IP, 0)
	for _, iface := range ifaces {
		if (iface.Flags & net.FlagUp) == 0 {
			continue
		}
		if iface.Index >= minIndex && len(ips) != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for i, rawAddr := range addrs {
			var ip net.IP
			switch addr := rawAddr.(type) {
			case *net.IPAddr:
				ip = addr.IP
			case *net.IPNet:
				ip = addr.IP
			default:
				continue
			}
			if isValidIP(ip.String()) {
				minIndex = iface.Index
				if i == 0 {
					ips = make([]net.IP, 0, 1)
				}
				ips = append(ips, ip)
				if ip.To4() != nil {
					break
				}
			}
		}
	}
	if len(ips) != 0 {
		return net.JoinHostPort(ips[len(ips)-1].String(), port), nil
	}
	return "", nil
}
