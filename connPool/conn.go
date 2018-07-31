package connpool

import (
	"errors"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
)

// GConn wrap net.Conn to borrow or return conn
type GConn struct {
	// wrap real connection
	net.Conn
	// gpool
	p *GPool
	//sync pool put or get
	mu sync.RWMutex
	// identify an GConn usable or can close
	unusable bool
	// key store the ip:port
	key string
	// connectMaxRetries
	connectMaxRetries int
	// connectMinRetry
	connectMinRetry time.Duration
	// uuid
	uuid string
}

// NewGConn return a new GConn
func NewGConn() *GConn {
	uuid := uuid.New()
	return &GConn{uuid: uuid.String()}
}

// Close puts the given connects back to the pool instead of closing it.
func (g *GConn) Close() error {
	if g.Conn != nil {
		g.p.addRemainingSpace()
		return g.Conn.Close()
	}
	return nil
}

// Connect real connect
func (g *GConn) Connect() error {
	connectAttempts := 0
	for connectAttempts < g.connectMaxRetries {
		conn, err := net.Dial(g.p.factory.factoryconfig.protocol, g.key)
		if err != nil {
			time.Sleep(time.Second * g.connectMinRetry)
			connectAttempts++
			continue
		}
		g.Conn = conn
		return nil
	}
	return errors.New("Connect fail after retry")
}
