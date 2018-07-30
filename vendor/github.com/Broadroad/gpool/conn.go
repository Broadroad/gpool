package gpool

import (
	"net"
	"sync"
)

// GConn wrap net.Conn to borrow or return conn
type GConn struct {
	// wrap real connection
	net.Conn
	// gpool
	p *gPool
	//sync pool put or get
	mu sync.RWMutex
	// identify an GConn usable or can close
	unusable bool
}

// Close puts the given connects back to the pool instead of closing it.
func (g *GConn) Close() error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.unusable {
		if g.Conn != nil {
			g.p.addRemainingSpace()
			return g.Conn.Close()
		}
		return nil
	}
	return g.p.Return(g.Conn)
}

// MarkUnusable marks the connection not usable any more, to let the pool close it instead of returning it to pool.
func (g *GConn) MarkUnusable() {
	g.mu.Lock()
	g.unusable = true
	g.mu.Unlock()
}
