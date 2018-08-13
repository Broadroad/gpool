package bufferpool

import "bytes"

type BufferPool struct {
	c     chan *bytes.Buffer
	alloc int
}
