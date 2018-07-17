package gpool

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	network = "tcp"
	address = "127.0.0.1:8080"
	factory = func() (net.Conn, error) { return net.Dial(network, address) }
)

func init() {
	go simpleTCPServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled
}
func TestNew(t *testing.T) {
	_, err := NewConnPool(1, 3, factory)
	if err != nil {
		t.Errorf("New error: %s", err)
	}
}

func TestPool_Get_Impl(t *testing.T) {
	p, _ := NewConnPool(1, 3, factory)
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	_, ok := conn.(*GConn)
	if !ok {
		t.Errorf("Conn is not of type poolConn")
	}
}

func simpleTCPServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func(conn net.Conn) {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}(conn)
	}
}
