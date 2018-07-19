package gpool

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	network    = "tcp"
	address    = "127.0.0.1:8080"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
	poolConfig = &PoolConfig{
		InitCap:     5,
		MaxCap:      30,
		Factory:     factory,
		IdleTimeout: 15 * time.Second,
	}
)

func init() {
	go simpleTCPServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled
}
func TestNew(t *testing.T) {
	_, err := NewGPool(poolConfig)
	if err != nil {
		t.Errorf("New error: %s", err)
	}
}

func TestGet(t *testing.T) {
	p, _ := NewGPool(poolConfig)
	defer p.Close()

	conn, err := p.Get()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}

	_, ok := conn.(*GConn)
	if !ok {
		t.Errorf("Conn is not of type GConn")
	}
}

func TestPressGet(t *testing.T) {
	p, _ := NewGPool(poolConfig)
	defer p.Close()
	done := make(chan struct{})

	for i := 0; i < 2000; i++ {
		go func() {
			defer func() {
				done <- struct{}{}
			}()
			conn, err := p.Get()
			if err != nil {
				t.Errorf("Get error: %s", err)
			}

			_, ok := conn.(*GConn)
			if !ok {
				t.Errorf("Conn is not of type GConn")
			} else {
				conn.Close()
			}
		}()
	}

	for i := 0; i < 2000; i++ {
		<-done
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
