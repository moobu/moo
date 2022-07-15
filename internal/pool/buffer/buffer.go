package buffer

import (
	"bytes"
	"sync"
)

var pool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func Get() *bytes.Buffer {
	return pool.Get().(*bytes.Buffer)
}

func Put(buf *bytes.Buffer) {
	buf.Reset()
	pool.Put(buf)
}
