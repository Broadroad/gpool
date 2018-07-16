package gpool

import (
	"net"
	"sync"
)

// GConn wrap net.Conn to borrow or return conn
type GConn struct {
	net.Conn
	p        *connPool
	mu       sync.RWMutex
	unusable bool
}

// Close puts the given connects back to the pool instead of closing it.
func (g *GConn) Close() error {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if g.unusable {
		if g.Conn != nil {
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
