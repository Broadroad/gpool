package connpool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	. "github.com/gpool"
)

// PoolConfig used for config the connection pool
type PoolConfig struct {
	// InitCap of the connection pool
	InitCap int
	// Maxcap is max connection number of the pool
	MaxCap int
	// WaitTimeout is the timeout for waiting to borrow a connection
	WaitTimeout time.Duration
	// 127.0.0.1:8080
	key string
}

//gPool store connections and pool info
//gPool key -> conns, each ip:port mapping a connection pool
type gPool struct {
	conns      chan *GConn
	factory    *Factory
	mu         sync.RWMutex
	poolConfig *PoolConfig
	idleConns  int
	createNum  int
	//will be used for blocking calls
	remainingSpace chan bool
}

func (p *gPool) addRemainingSpace() {
	p.remainingSpace <- true
}

func (p *gPool) removeRemainingSpace() {
	<-p.remainingSpace
}

// NewGPool create a connection pool
func NewGPool(pc *PoolConfig, fc *FactoryConfig) (Pool, error) {
	// test initCap and maxCap
	if pc.InitCap < 0 || pc.MaxCap < 0 || pc.InitCap > pc.MaxCap {
		return nil, errors.New("invalid capacity setting")
	}
	p := &gPool{
		conns:          make(chan *GConn, pc.MaxCap),
		poolConfig:     pc,
		idleConns:      pc.InitCap,
		factory:        NewFactory(fc),
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
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		p.createNum = pc.InitCap
		p.conns <- &GConn{Conn: conn}
	}
	return p, nil
}

func (p *gPool) Borrow() (*GConn, error) {
	p.mu.Lock()
	if p.conns == nil {
		return nil, ErrNil
	}
	p.mu.Unlock()

	select {
	case conn := <-p.conns:
		if conn == nil {
			return nil, ErrClosed
		}

		p.mu.Lock()
		p.idleConns--
		p.mu.Unlock()

		p.factory.ActiveObject(p.factory.factoryconfig.key, conn)
		return conn, nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.createNum++
		if p.createNum > p.poolConfig.MaxCap {
			return nil, errors.New("More than MaxCap")
		}
		conn, err := p.factory.Create()
		p.removeRemainingSpace()

		if err != nil {
			p.addRemainingSpace()
			return nil, err
		}

		return &GConn{Conn: conn}, nil
	}
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
	case p.conns <- &GConn{Conn: conn}:
		p.idleConns++
		return nil
	default:
		// pool is full, close passed connection
		return conn.Close()
	}
}

// Get implement Pool get interface
// if don't have any connection available, it will try to new one
func (p *gPool) Get() (net.Conn, error) {
	p.mu.Lock()
	if p.conns == nil {
		return nil, ErrNil
	}
	p.mu.Unlock()

	// wrap our connections with out custom net.Conn implementation (wrapConn
	// method) that puts the connection back to the pool if it's closed.
	select {
	case conn := <-p.conns:
		if &conn == nil {
			return nil, ErrClosed
		}

		p.mu.Lock()
		p.idleConns--
		p.mu.Unlock()
		return conn.Conn, nil
	default:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.createNum++
		if p.createNum > p.poolConfig.MaxCap {
			return nil, errors.New("More than MaxCap")
		}
		conn, err := p.factory.Create()
		p.removeRemainingSpace()

		if err != nil {
			p.addRemainingSpace()
			return nil, err
		}

		return conn, nil
	}
}

// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
// available.
func (p *gPool) BlockingGet(ctx context.Context) (net.Conn, error) {
	p.mu.RLock()
	conns := p.conns
	p.mu.Unlock()

	if conns == nil {
		return nil, ErrNil
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
			return nil, ErrClosed
		}

		p.mu.Lock()
		p.idleConns--
		p.mu.Unlock()
		return conn, nil
	case _ = <-p.remainingSpace:
		p.mu.Lock()
		defer p.mu.Unlock()
		p.createNum++
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
func (p *gPool) Close() {
	p.mu.Lock()
	conns := p.conns
	p.conns = nil
	p.mu.Unlock()

	if conns == nil {
		return
	}

	close(conns)
	for conn := range conns {
		conn.Close()
		p.addRemainingSpace()
	}
}

// Len implement Pool Len interface
// it will return current length of the pool
func (p *gPool) Len() int {
	p.mu.RLock()
	length := len(p.conns)
	p.mu.Unlock()
	return length
}

// Idle implement Pool Idle interface
// it will return current idle length of the pool
func (p *gPool) Idle() int {
	p.mu.RLock()
	idle := p.idleConns
	p.mu.Unlock()
	return int(idle)
}
