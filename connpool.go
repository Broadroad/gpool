package gpool

import "net"

type connPool struct {
	conns   chan net.Conn
	factory Factory
}

// Factory generate a new connection
type Factory func() (net.Conn, error)

// NewConnPool create a connection pool
func NewConnPool(initCap, maxCap int64, factory Factory) (Pool, error) {
	c := &connPool{
		conns:   make(chan net.Conn, maxCap),
		factory: factory,
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
