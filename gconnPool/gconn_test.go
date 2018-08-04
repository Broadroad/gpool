package connpool

import (
	"fmt"
	"testing"
	"time"
)

var (
	key               = "127.0.0.1:8080"
	connectMaxRetries = 10
	connectMinRetry   = 1 * time.Second
	protocol          = "tcp"
)

// TestNewGConn test connpool
func TestNewGConn(t *testing.T) {
	gconn := NewGConn(key, connectMaxRetries, connectMinRetry, protocol)
	fmt.Println(gconn.Uuid)
}

// TestClose test Close
func TestClose(t *testing.T) {
	gconn := NewGConn(key, connectMaxRetries, connectMinRetry, protocol)
	err := gconn.Close()
	if err != nil {
		fmt.Println("err: ", err)
	}
}

// TestConnect test Connect
func TestConnect(t *testing.T) {
	gconn := NewGConn(key, connectMaxRetries, connectMinRetry, protocol)
	err := gconn.Connect()
	if err != nil {
		t.Error(err)
	}
	if gconn.Conn == nil {
		t.Error("gconn.Conn is nil")
	}
	err = gconn.Close()
	if err != nil {
		fmt.Println("err: ", err)
	}
}
