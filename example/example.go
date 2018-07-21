package main

import (
	"fmt"
	"net"

	"github.com/broadroad/gpool"
)

func main() {
	// factory is the function that create connection
	factory := func() (net.Conn, error) { return net.Dial("tcp", "127.0.0.1") }

	// poolConfig is the config of gpool
	poolConfig := &gpool.PoolConfig{
		InitCap: 5,
		MaxCap:  30,
		Factory: factory,
	}

	// create a new conn pool
	p, err := gpool.NewGPool(poolConfig)
	if err != nil {
		fmt.Println("new pool error is", err)
	}

	// release pool
	defer p.Close()

	// get a new connection from pool
	conn, err := p.Get()
	if err != nil {
		fmt.Println("Get error: %s", err)
	}

	_, ok := conn.(*GConn)
	if !ok {
		fmt.Println("Conn is not of type GConn")
	}

	// return connection to pool
	conn.Close()

	// return len of the pool
	fmt.Println("len=", p.Len())

}
