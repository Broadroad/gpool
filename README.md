# gpool[![GoDoc](http://godoc.org/github.com/silenceper/pool?status.svg)](http://godoc.org/github.com/silenceper/pool) [![Build Status](https://travis-ci.org/Broadroad/gpool.svg?branch=master)](https://travis-ci.org/Broadroad/gpool)

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
install with this command:
```shell
go get github.com/broadroad/gpool
```

and then use like this bellow:

```go
// create factory to create connection
factory    = func() (net.Conn, error) { return net.Dial(network, address) }

// create poolConfig
poolConfig = &PoolConfig{
	InitCap:     5,
	MaxCap:      30,
	Factory:     factory,
}

// create a new gpool
p, err := NewGPool(poolConfig)
if err != nil {
    fmt.Println(err)
}

// get a connection from gpool, if gpool has no idle connection, it will return error
conn, err := p.Get()
if err != nil {
	fmt.Println("Get error: ", err)
}

// return a connection to gpool
conn.Close()

// BlockingGet will block until it has a idle connection
conn, err := p.BlockingGet()
if err != nil {
	fmt.Println("Get error: ", err)
}

// return a connection to gpool
conn.Close()

// release all connection in gpool
p.Close()

```

## License
The Apache License 2.0 - see LICENCE for more details

## Issue
It will be very pleasure if you give some issue or pr. Feel free to contact tjbroadroad@163.com
