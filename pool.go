// Package gpool implements a tcp connection pool to manage connection and use
package gpool

import (
	"context"
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

	// BlockingGet will block until it gets an idle connection from pool. Context timeout can be passed with context
	// to wait for specific amount of time. If nil is passed, this will wait indefinitely until a connection is
	// available.
	BlockingGet(context.Context) (net.Conn, error)
}
