package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Broadroad/gpool"
)

var (
	network    = "tcp"
	address    = "127.0.0.1:8080"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
	poolConfig = &gpool.PoolConfig{
		InitCap:     5,
		MaxCap:      30,
		Factory:     factory,
		IdleTimeout: 15 * time.Second,
	}
)

func main() {

	// create a new conn pool
	p, err := gpool.NewGPool(poolConfig)
	if err != nil {
		fmt.Println("new pool error is", err)
	}

	// release pool
	defer p.Close()

	done := make(chan struct{})

	for i := 0; i < 20; i++ {
		go func() {

			defer func() {
				done <- struct{}{}
			}()

			// get a new connection from pool
			conn, err := p.Get()
			if err != nil {
				fmt.Println("Get error:", err)
			}

			_, ok := conn.(*gpool.GConn)
			if !ok {
				fmt.Println("Conn is not of type GConn")
			} else {
				time.Sleep(time.Second * 1)
				conn.Close()
			}
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
	}

	// return len of the pool
	fmt.Println("len = ", p.Len())

}

func tcpServer() {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Panicln(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Panicln(err)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]
		conn.Write(data)
	}
}
