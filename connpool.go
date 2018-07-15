package gpool

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

type connPool struct {
	conns   chan net.Conn
	factory Factory
	mu      sync.RWMutex
}

// Factory generate a new connection
type Factory func() (net.Conn, error)

// NewConnPool create a connection pool
func NewConnPool(initCap, maxCap int, factory Factory) (Pool, error) {
	// test initCap and maxCap
	if initCap < 0 || maxCap < 0 || initCap > maxCap {
		return nil, errors.New("invalid capacity setting")
	}
	c := &connPool{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
	}

	// create initial connection, if wrong just close it
	for i := 0; i < initCap; i++ {
		conn, err := factory()
		if err != nil {
			c.Close()
			return nil, fmt.Errorf("factory is not able to fill the pool: %s", err)
		}
		c.conns <- conn
	}

	return c, nil
}

// Get implement Pool get interface
// if don't have any connection avaliable, it will try to new one
func (c *connPool) Get() (net.Conn, error) {
	return nil, nil
}

// Close implement Pool close interface
// it will close all the connection in the pool
func (c *connPool) Close() {

}

// Len implement Pool Len interface
// it will return current length of the pool
func (c *connPool) Len() {

}

// Idle implement Pool Idle interface
// it will return current idle length of the pool
func (c *connPool) Idle() {

}
