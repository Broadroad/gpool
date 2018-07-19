package gpool

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

//connPool store connections and pool info
type connPool struct {
	conns   chan net.Conn
	factory Factory
	mu      sync.RWMutex
	idle    int64
}

// Factory generate a new connection
type Factory func() (net.Conn, error)

// NewConnPool create a connection pool
func NewConnPool(initCap, maxCap int, factory Factory) (Pool, error) {
	// test initCap and maxCap
	if initCap < 0 || maxCap < 0 || initCap > maxCap {
		return nil, errors.New("invalid capacity setting")
	}
	p := &connPool{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
		idle:    int64(initCap),
	}

	// create initial connection, if wrong just close it
	for i := 0; i < initCap; i++ {
		conn, err := factory()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		p.conns <- conn
	}

	return p, nil
}

// wrapConn wraps a standard net.Conn to a poolConn net.Conn.
func (p *connPool) wrapConn(conn net.Conn) net.Conn {
	gconn := &GConn{p: p}
	gconn.Conn = conn
	return gconn
}

// getConnsAndFactory get conn channel and factory by once
func (p *connPool) getConnsAndFactory() (chan net.Conn, Factory) {
	p.mu.RLock()
	conns := p.conns
	factory := p.factory
	p.mu.RUnlock()
	return conns, factory
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (p *connPool) Return(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	if p.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.conns <- conn:
		atomic.AddInt64(&p.idle, -1)
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// Get implement Pool get interface
// if don't have any connection avaliable, it will try to new one
func (p *connPool) Get() (net.Conn, error) {
	conns, factory := p.getConnsAndFactory()
	if conns == nil {
		return nil, ErrClosed
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		atomic.AddInt64(&p.idle, 1)
		return p.wrapConn(conn), nil
	default:
		conn, err := factory()
		if err != nil {
			return nil, err
		}

		return p.wrapConn(conn), nil
	}
}

// Close implement Pool close interface
// it will close all the connection in the pool
func (p *connPool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.factory = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
	}
}

// Len implement Pool Len interface
// it will return current length of the pool
func (p *connPool) Len() int {
	conns, _ := p.getConnsAndFactory()
	return len(conns)
}

// Idle implement Pool Idle interface
// it will return current idle length of the pool
func (p *connPool) Idle() int {
	return int(p.idle)
}
