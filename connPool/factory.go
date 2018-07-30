package connpool

import (
	"net"
	"time"
)

// Factory is manage borrow or return connection from pool
type Factory struct {
	factoryconfig *FactoryConfig
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
	key             string // 127.0.0.1:8080
}

// NewFactory return a new factory
func NewFactory(fc *FactoryConfig) *Factory {
	factory := &Factory{factoryconfig: fc}
	return factory
}

// Create a new conn instance
func (f *Factory) Create() (net.Conn, error) {
	return create(f.factoryconfig.key, f.factoryconfig.lazyCreate, f.factoryconfig.protocol)
}

func create(key string, lazyCreate bool, protocol string) (net.Conn, error) {
	if lazyCreate {
		return NewGConn(), nil
	}

	return net.Dial(protocol, key)
}

// DestoryObject destory the conn instance
func (f *Factory) DestoryObject(key string, g *GConn) error {
	g.Conn.Close()
	return nil
}

// ValidateObject validate whehter the connection is connected
func (f *Factory) ValidateObject(key string, g *GConn) bool {
	if g.Conn != nil {
		return true
	}
	return false
}

// ActiveObject really connect when called
func (f *Factory) ActiveObject(key string, g *GConn) error {
	if g.Conn == nil {
		return g.Connect()
	}
	return nil
}
