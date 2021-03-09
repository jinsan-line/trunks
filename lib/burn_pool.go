package trunks

import (
	"fmt"
	"strings"
	"sync"

	"google.golang.org/grpc"
)

var (
	errPoolNoConn = fmt.Errorf("no connection in pool")
)

// simple client-side round-robin pool
type pool struct {
	conns []*grpc.ClientConn
	mu    sync.Mutex
	next  int
}

func newPool(hosts []string, numConnPerHost uint64, cos []grpc.CallOption) (*pool, error) {
	p := &pool{}
	err := p.connect(hosts, numConnPerHost, cos)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *pool) connect(hosts []string, numConnPerHost uint64, cos []grpc.CallOption) error {
	for _, h := range hosts {
		var x uint64
		for ; x < numConnPerHost; x++ {
			var c *grpc.ClientConn
			var err error

			if len(cos) > 0 {
				c, err = grpc.Dial(h, grpc.WithDefaultCallOptions(cos...), grpc.WithInsecure())
			} else {
				c, err = grpc.Dial(h, grpc.WithInsecure())
			}

			if err != nil {
				return fmt.Errorf("failed to dial %s: %v", h, err)
			}
			p.conns = append(p.conns, c)
		}
	}

	if len(p.conns) < 1 {
		return fmt.Errorf("no connection in pool")
	}

	return nil
}

func (p *pool) pick() (*grpc.ClientConn, error) {
	size := len(p.conns)
	if size == 0 {
		return nil, errPoolNoConn
	}

	p.mu.Lock()
	defer func() {
		p.mu.Unlock()
		p.next = (p.next + 1) % size
	}()
	return p.conns[p.next], nil
}

func (p *pool) close() error {
	var errs []string
	for _, c := range p.conns {
		if err := c.Close(); err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s", strings.Join(errs, "; "))
	}
	return nil
}
