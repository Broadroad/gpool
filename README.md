# gpool
A go tcp connection pool

## Function
- gpool support net.Conn interface
- reuse connection in gpool
- get connection from gpool will error when pool is full

## TODO
- get connection will block until timeout or a idle connection return
- connection will be closed when idle for some time duration

## License
The Apache License 2.0 - see LISENCE for more details

## Issue
It will be very pleasure if you give some issue or pr. Feel free to contact tjbroadroad@163.com