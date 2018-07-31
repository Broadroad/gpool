package connpool

import (
	"log"
	"net"
	"testing"
	"time"
)

var (
	network = "tcp"
	address = "127.0.0.1:8080"
	//	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
	poolConfig = &PoolConfig{
		InitCap: 5,
		MaxCap:  30,
	}
	factoryConfig = &FactoryConfig{
		connectTimeout:    10,
		connectMaxRetries: 10,
		lazyCreate:        false,
		protocol:          "tcp",
		key:               "127.0.0.1:8080",
	}
)

func init() {
	go tcpServer()
	time.Sleep(time.Millisecond * 300) // wait until tcp server has been settled
}

func TestGet(t *testing.T) {
	p, _ := NewGPool(poolConfig, factoryConfig)
	defer p.Close()

	gconn, err := p.Borrow(address)
	p.Return(gconn)
	gconn.Close()
	if err != nil {
		t.Errorf("Get error: %s", err)
	}
}

/*func TestPressGet(t *testing.T) {
	p, _ := NewGPool(poolConfig, factoryConfig)
	defer p.Close()
	done := make(chan struct{})

	for i := 0; i < 20; i++ {
		go func() {
			defer func() {
				done <- struct{}{}
			}()
			conn, err := p.Borrow(address)
			if err != nil {
				t.Errorf("Get error: %s", err)
			}

			time.Sleep(time.Second * 1)
			conn.Close()
		}()
	}

	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestBlockingGetWithTimeout(t *testing.T) {
	p, _ := NewGPool(poolConfig, factoryConfig)
	defer p.Close()
	done := make(chan struct{})

	//context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := 0; i < 30; i++ {
		go func(i int) {
			defer func() {
				done <- struct{}{}
			}()

			conn, err := p.BlockingBorrow(ctx)
			if err != nil {
				t.Errorf("Get error: %s", err)
			}

			_, ok := conn.(*GConn)
			if !ok {
				t.Errorf("Conn is not of type GConn")
			} else {
				conn.Write([]byte("Message"))
				r := make([]byte, 1024)
				time.Sleep(time.Second * 1)
				_, err := conn.Read(r)
				if err != nil {
					t.Error("error reading from conn", err)
				}
				conn.Close()
			}
		}(i)
	}

	for i := 0; i < 30; i++ {
		<-done
	}
}

func TestBlockingGetWithoutTimeout(t *testing.T) {
	p, _ := NewGPool(poolConfig, factoryConfig)
	defer p.Close()
	done := make(chan struct{})

	//context without timeout - will block indefinitely
	ctx := context.Background()

	for i := 0; i < 30; i++ {
		go func(i int) {
			defer func() {
				done <- struct{}{}
			}()

			conn, err := p.BlockingBorrow(ctx)
			if err != nil {
				t.Errorf("Get error: %s", err)
			}

			_, ok := conn.(*GConn)
			if !ok {
				t.Errorf("Conn is not of type GConn")
			} else {
				conn.Write([]byte("Message"))
				r := make([]byte, 1024)
				time.Sleep(time.Second * 1)
				_, err := conn.Read(r)
				if err != nil {
					t.Error("error reading from conn", err)
				}
				conn.Close()
			}
		}(i)
	}

	for i := 0; i < 30; i++ {
		<-done
	}
}

func TestBlockingGetWithNil(t *testing.T) {
	p, _ := NewGPool(poolConfig, factoryConfig)
	defer p.Close()
	done := make(chan struct{})

	for i := 0; i < 30; i++ {
		go func(i int) {
			defer func() {
				done <- struct{}{}
			}()
			//nil as ctx - will block indefinitely
			conn, err := p.BlockingBorrow(nil)
			if err != nil {
				t.Errorf("Get error: %s", err)
			}

			_, ok := conn.(*GConn)
			if !ok {
				t.Errorf("Conn is not of type GConn")
			} else {
				conn.Write([]byte("Message"))
				r := make([]byte, 1024)
				time.Sleep(time.Second * 1)
				_, err := conn.Read(r)
				if err != nil {
					t.Error("error reading from conn", err)
				}
				conn.Close()
			}
		}(i)
	}

	for i := 0; i < 30; i++ {
		<-done
	}
}
*/

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
	log.Println("Accepted new connection.")
	defer conn.Close()
	defer log.Println("Closed connection.")

	for {
		buf := make([]byte, 1024)
		size, err := conn.Read(buf)
		if err != nil {
			return
		}
		data := buf[:size]
		//log.Println("Read new data from connection", data)
		conn.Write(data)
		//log.Println("Wrote data to connection", data)
	}
}
