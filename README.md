# gpool
[![GoDoc](http://godoc.org/github.com/silenceper/pool?status.svg)](http://godoc.org/github.com/silenceper/pool)
A go tcp connection pool

## Function
- gpool support net.Conn interface
- reuse connection in gpool
- get connection from gpool will error when pool is full

## TODO
- get connection will block until timeout or a idle connection return
- connection will be closed when idle for some time duration

## Usage
```go
// create factory to create connection
factory    = func() (net.Conn, error) { return net.Dial(network, address) }

// create poolConfig
poolConfig = &PoolConfig{
	InitCap:     5,
	MaxCap:      30,
	Factory:     factory,
	IdleTimeout: 15 * time.Second,
}

// create a new gpool
p, err := NewGPool(poolConfig)
if err != nil {
    fmt.Println(err)
}

// get a connection from gpool
conn, err := p.Get()
if err != nil {
	t.Errorf("Get error: %s", err)
}

// return a connection to gpool
conn.Close()

// release all connection in gpool
p.Close()

```

## License
The Apache License 2.0 - see LISENCE for more details

## Issue
It will be very pleasure if you give some issue or pr. Feel free to contact tjbroadroad@163.com
