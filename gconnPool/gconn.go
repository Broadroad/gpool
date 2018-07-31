package connpool

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/google/uuid"
)

// GConn wrap net.Conn to borrow or return conn
type GConn struct {
	// wrap real connection
	net.Conn
	// key store the ip:port
	key string
	// connectMaxRetries
	connectMaxRetries int
	// connectMinRetry is time wait for next connect try
	connectMinRetry time.Duration
	// uuid
	uuid string
	// protocol
	protocol string
}

// NewGConn return a new GConn
func NewGConn(key string, connectMaxRetries int, connectMinRetry time.Duration, protocol string) *GConn {
	return &GConn{uuid: uuid.New().String(), connectMaxRetries: 10, protocol: "tcp", key: "127.0.0.1:8080"}
}

// Close puts the given connects back to the pool instead of closing it.
func (g *GConn) Close() error {
	if g.Conn != nil {
		err := g.Conn.Close()
		g.Conn = nil
		return err
	}
	return nil
}

// Connect real connect, it will try connectAttempts and wait connectMinRetry time
func (g *GConn) Connect() error {
	connectAttempts := 0
	for connectAttempts < g.connectMaxRetries {
		conn, err := net.Dial(g.protocol, g.key)
		if err != nil {
			time.Sleep(time.Second * g.connectMinRetry)
			fmt.Println(err)
			connectAttempts++
			continue
		}
		g.Conn = conn
		return nil
	}
	return errors.New("Connect fail after retry")
}
