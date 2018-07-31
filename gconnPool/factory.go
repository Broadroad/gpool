package connpool

import (
	"time"
)

// Factory is manage borrow or return connection from pool
type Factory struct {
	factoryconfig *FactoryConfig
}

// FactoryConfig manage factory config
type FactoryConfig struct {
	connectTimeout    time.Duration
	connectMaxRetries int
	connectMinRetry   time.Duration
	protocol          string // here is tcp
	lazyCreate        bool   // create when pool create or when using it
	key               string // 127.0.0.1:8080
}

// NewFactory return a new factory
func NewFactory(fc *FactoryConfig) (*Factory, error) {
	if fc.connectTimeout < 0 || fc.connectMaxRetries < 0 || fc.connectMinRetry < 0 {
		return nil, ParameterERROR
	}
	factory := &Factory{factoryconfig: fc}
	return factory, nil
}

// Create a new gconn instance
func (f *Factory) Create() (*GConn, error) {
	return create(f.factoryconfig.key, f.factoryconfig.connectMaxRetries, f.factoryconfig.connectMinRetry, f.factoryconfig.lazyCreate, f.factoryconfig.protocol)
}

func create(key string, connectMaxRetries int, connectMinRetry time.Duration, lazyCreate bool, protocol string) (*GConn, error) {
	if lazyCreate {
		return NewGConn(key, connectMaxRetries, connectMinRetry, protocol), nil
	}

	gconn := NewGConn(key, connectMaxRetries, connectMinRetry, protocol)
	err := ActiveObject(gconn)
	if err != nil {
		return nil, err
	}
	return gconn, nil
}

// DestoryObject destory the conn instance
func DestoryObject(g *GConn) error {
	g.Conn.Close()
	return nil
}

// ValidateObject validate whehter the connection is connected
func ValidateObject(g *GConn) bool {
	if g.Conn != nil {
		return true
	}
	return false
}

// ActiveObject really connect when called
func ActiveObject(g *GConn) error {
	if g.Conn == nil {
		return g.Connect()
	}
	return nil
}
