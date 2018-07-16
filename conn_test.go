package gpool

import (
	"net"
	"testing"
)

// TestConn_Impl test connpool
func TestConn_Impl(t *testing.T) {
	var _ net.Conn = new(GConn)
}
