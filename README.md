# gpool[![GoDoc](https://godoc.org/github.com/Broadroad/gpool?status.svg)](https://godoc.org/github.com/Broadroad/gpool) [![Build Status](https://travis-ci.org/Broadroad/gpool.svg?branch=master)](https://travis-ci.org/Broadroad/gpool)

A go tcp connection pool

## Function
- gpool support net.Conn interface
- reuse connection in gpool
- get connection from gpool will error when there is no idle connection in gpool
- support block get from gpool when there is no idle connection in gpool

## Todo
- Add a timeout in BlockingGet
- Connection will be closed when idle for some time duration(keep idle connection alive for some time that users can config)

## Usage
### install with this command:
```shell
go get github.com/broadroad/gpool
```

### and then set the poolConfig:

```go
// create factory to create connection
factory    = func() (net.Conn, error) { return net.Dial(network, address) }

// create poolConfig
poolConfig = &PoolConfig{
	InitCap:     5,
	MaxCap:      30,
	Factory:     factory,
}
```

### Non blocking get, if no idle connection then return error.
```go
// create a new gpool
p, err := NewGPool(poolConfig)
if err != nil {
    fmt.Println(err)
}
// release all connection in gpool
p.Close()

// get a connection from gpool, if gpool has no idle connection, it will return error
conn, err := p.Get()
if err != nil {
	fmt.Println("Get error: ", err)
}

// return a connection to gpool
defer conn.Close()
```

### Blocking get, if no idle connection then block until a time out

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //3second timeout
defer cancel()
conn, err := p.BlockingGet(ctx)
if err != nil {
	fmt.Println("Get error:", err)
}
defer conn.Close()
```

## License
The Apache License 2.0 - see LICENSE for more details

## Issue
It will be very pleasure if you give some issue or pr. Feel free to contact tjbroadroad@163.com
