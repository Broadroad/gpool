package gpool

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// PoolConfig used for config the connection pool
type PoolConfig struct {
	// InitCap of the connection pool
	InitCap int
	// Maxcap is max connection number of the pool
	MaxCap int
	// WaitTimeout is the timeout for waiting to borrow a connection
	WaitTimeout time.Duration
	// IdleTimeout is the timeout for a connection to be alive
	IdleTimeout time.Duration
	Factory     func() (net.Conn, error)
}

//gPool store connections and pool info
type gPool struct {
	conns         chan net.Conn
	factory       Factory
	mu            sync.RWMutex
	poolConfig    *PoolConfig
	idleConns     int
	borrowedConns int
}

// Factory generate a new connection
type Factory func() (net.Conn, error)

// NewGPool create a connection pool
func NewGPool(pc *PoolConfig) (Pool, error) {
	// test initCap and maxCap
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, errors.New("invalid capacity setting")
	}
	p := &gPool{
		conns:      make(chan net.Conn, pc.MaxCap),
		factory:    pc.Factory,
		poolConfig: pc,
		idleConns:  pc.InitCap,
	}

	// create initial connection, if wrong just close it
	for i := 0; i < pc.InitCap; i++ {
		conn, err := pc.Factory()
		if err != nil {
			p.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		p.conns <- conn
	}

	return p, nil
}

// wrapConn wraps a standard net.Conn to a poolConn net.Conn.
func (p *gPool) wrapConn(conn net.Conn) net.Conn {
	gconn := &GConn{p: p}
	gconn.Conn = conn
	return gconn
}

// getConnsAndFactory get conn channel and factory by once
func (p *gPool) getConnsAndFactory() (chan net.Conn, Factory) {
	p.mu.RLock()
	conns := p.conns
	factory := p.factory
	p.mu.RUnlock()
	return conns, factory
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (p *gPool) Return(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil. rejecting")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.conns <- conn:
		p.idleConns++
		return nil
	default:
		// pool is full, close passed connection
		p.borrowedConns--
		return conn.Close()
	}
}

// Get implement Pool get interface
// if don't have any connection avaliable, it will try to new one
func (p *gPool) Get() (net.Conn, error) {
	conns, factory := p.getConnsAndFactory()
	if conns == nil {
		return nil, ErrNil
	}

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, ErrClosed
		}

		p.idleConns--
		return p.wrapConn(conn), nil
	default:

		conn, err := factory()
		p.mu.Lock()
		if p.borrowedConns > p.poolConfig.MaxCap {
			return nil, err
		}
		p.borrowedConns++
		defer p.mu.Unlock()
		if err != nil {
			return nil, err
		}

		return p.wrapConn(conn), nil
	}
}

// Close implement Pool close interface
// it will close all the connection in the pool
func (p *gPool) Close() {
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
func (p *gPool) Len() int {
	conns, _ := p.getConnsAndFactory()
	return len(conns)
}

// Idle implement Pool Idle interface
// it will return current idle length of the pool
func (p *gPool) Idle() int {
	return int(p.idleConns)
}
