package connpool

import (
	"net"
	"time"
)

// factory is manage borrow or return connection from pool
type factory struct {
	factoryconfig FactoryConfig
}

// FactoryConfig manage factory config
type FactoryConfig struct {
	connectTimeout  time.Duration
	connectRetries  int
	connectMinRetry time.Duration
	idleTimeout     time.Duration
	poolConnections bool   // use connection pool or not
	protocol        string // here is tcp
	lazyCreate      bool   // create when pool create or when using it
}

// NewFactory return a new factory
func NewFactory(fc FactoryConfig) *factory {
	factory := &factory{factoryconfig: fc}
	return factory
}

// Create a new conn instance
func (f *factory) Create(key string) (net.Conn, error) {
	if f.factoryconfig.lazyCreate {
		return NewGConn(key), nil
	}

	return net.Dial(f.factoryconfig.protocol, key)
}

// DestoryObject destory the conn instance
func (f *factory) DestoryObject(key string, g GConn) error {
	g.Conn.Close()
	return nil
}

// ValidateObject validate whehter the connection is connected
func (f *factory) ValidateObject(key string, g GConn) bool {
	if g.Conn != nil {
		return true
	}
	return false
}

// ActiveObject really connect when called
func (f *factory) ActiveObject(key string, g GConn) error {
	if g.Conn == nil {
		return g.Connect()
	}
	return nil
}
