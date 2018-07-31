package connpool

import (
	"context"
	"net"
	"sync"
)

// PoolConfig used for config the connection pool
type PoolConfig struct {
	// InitCap of the connection pool
	InitCap int
	// Maxcap is max connection number of the pool
	MaxCap int
	// 127.0.0.1:8080
	key string
}

//GPool store connections and pool info
type GPool struct {
	conns      chan *GConn
	factory    *Factory
	mu         sync.RWMutex
	poolConfig *PoolConfig
	//will be used for blocking calls
	remainingSpace chan bool
}

func (p *GPool) addRemainingSpace() {
	p.remainingSpace <- true
}

func (p *GPool) removeRemainingSpace() {
	<-p.remainingSpace
}

// NewGPool create a connection pool
func NewGPool(pc *PoolConfig, fc *FactoryConfig) (*GPool, error) {
	// test initCap and maxCap
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, ParameterERROR
	}

	factory, err := NewFactory(fc)
	if err != nil {
		return nil, err
	}

	p := &GPool{
		conns:          make(chan *GConn, pc.MaxCap),
		poolConfig:     pc,
		factory:        factory,
		remainingSpace: make(chan bool, pc.MaxCap),
	}

	//fill the remainingSpace channel so we can use it for blocking calls
	for i := 0; i < pc.MaxCap; i++ {
		p.addRemainingSpace()
	}

	// create initial connection, if wrong just close it
	for i := 0; i < pc.InitCap; i++ {
		conn, err := p.factory.Create()
		p.removeRemainingSpace()
		if err != nil {
			p.Close()
			p.addRemainingSpace()
			return nil, FillERROR
		}
		p.conns <- conn
	}
	return p, nil
}

// Borrow borrows a connection from factory, and factory active the conn object
// it is lazy connect way
func (p *GPool) Borrow() (*GConn, error) {
	p.mu.RLock()
	conns := p.conns
	p.mu.RUnlock()
	select {
	case conn := <-conns:
		ActiveObject(conn)
		return conn, nil
	case _ = <-p.remainingSpace:
		conn, err := p.factory.Create()
		if err != nil {
			p.addRemainingSpace()
			return nil, err
		}

		return conn, nil
	default:
		return nil, NilERROR
	}
}

// Return return the connection back to the pool. If the pool is full or closed,
// conn is simply closed. A nil conn will be rejected.
func (p *GPool) Return(conn *GConn) error {
	if conn == nil {
		return NilERROR
	}

	p.mu.RLock()
	if p.conns == nil {
		// pool is closed, close passed connection
		return conn.Close()
	}
	p.mu.RUnlock()

	// put the resource back into the pool. If the pool is full, this will
	// block and the default case will be executed.
	select {
	case p.conns <- conn:
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// BlockingBorrow will block until it gets an idle connection from pool. Context timeout can be passed with context
// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
// available.
func (p *GPool) BlockingBorrow(ctx context.Context) (net.Conn, error) {
	p.mu.RLock()
	conns := p.conns
	p.mu.Unlock()

	if conns == nil {
		return nil, NilERROR
	}
	//if context is nil it means we have no timeout, we can wait indefinitely
	if ctx == nil {
		ctx = context.Background()
	}
	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-conns:
		if conn == nil {
			return nil, NilERROR
		}

		return conn, nil
	case _ = <-p.remainingSpace:
		p.mu.Lock()
		defer p.mu.Unlock()
		conn, err := p.factory.Create()
		if err != nil {
			p.addRemainingSpace()
			return nil, err
		}

		return conn, nil
	//if context deadline is reached, return timeout error
	case <-ctx.Done():
		return nil, ctx.Err()
	}

}

// Close implement Pool close interface
// it will close all the connection in the pool
func (p *GPool) Close() {
	p.mu.RLock()
	conns := p.conns
	p.conns = nil
	p.mu.RUnlock()

	if conns == nil {
		return
	}

	close(conns)

	for conn := range conns {
		conn.Close()
	}
}
