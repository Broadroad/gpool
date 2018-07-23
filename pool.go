// Package gpool implements a tcp connection pool to manage connection and use
package gpool

import (
	"errors"
	"net"
)

var (
	ErrClosed = errors.New("pool has been closed")
	ErrNil    = errors.New("pool is nil")
)

// Pool is interface which all type of pool need to implement
type Pool interface {
	// Get returns a new connection from pool.
	Get() (net.Conn, error)

	// Close close the pool and reclaim all the connections.
	Close()

	// Len get the length of the pool
	Len() int

	// Idle get the idle connection pool number
	Idle() int

	// BlockingGet will block until it get a idle connection from pool
	BlockingGet() (net.Conn, error)
}
