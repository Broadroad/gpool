# gpool[![GoDoc](https://godoc.org/github.com/Broadroad/gpool?status.svg)](https://godoc.org/github.com/Broadroad/gpool) [![Build Status](https://travis-ci.org/Broadroad/gpool.svg?branch=master)](https://travis-ci.org/Broadroad/gpool)

A golang pool which will support connection pool, buffer pool, goroutine pool. Help developers to use pool easily. Now gpool only support tcp connection pool. It will support other pools soon. And Thanks to https://github.com/fatih/pool, ideas comes from fatih. 
 
## Function
- gpool support tcp connection now
- reuse tcp connection in gpool
- get connection from gpool will error when there is no idle connection in gpool
- support block get from gpool when there is no idle connection in gpool with timeout

## Todo
- support lazy connect when real borrow conn from gpool
- support buffer pool
- support goroutine pool

## Usage
### 1. install with this command:
```shell
go get github.com/broadroad/gpool
```

### 2. set the poolConfig:

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

### 3. create new gpool
```go
// create a new gpool
p, err := NewGPool(poolConfig)
if err != nil {
    fmt.Println(err)
}
// release all connection in gpool
defer p.Close()
```

### 4. non blocking get, if no idle connection then return error.
```go
// get a connection from gpool, if gpool has no idle connection, it will return error
conn, err := p.Get()
if err != nil {
	fmt.Println("Get error: ", err)
}

// return a connection to gpool
defer conn.Close()
```

### 5. blocking get, if no idle connection then block until a timeout
#### 5.1 Block until specified timeout
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) //3second timeout
defer cancel()
conn, err := p.BlockingGet(ctx)
if err != nil {
	fmt.Println("Get error:", err)
}

// return a connection to gpool
defer conn.Close()
```
#### 5.2 Block indefinitely
```go
conn, err := p.BlockingGet(nil)
if err != nil {
	fmt.Println("Get error:", err)
}

// return a connection to gpool
defer conn.Close()
```

## License
The Apache License 2.0 - see LICENSE for more details

## Issue
It will be very pleasure if you give some issue or pr. Feel free to contact tjbroadroad@163.com

## Contributor
* [Farhad Farahi](https://github.com/FarhadF)
