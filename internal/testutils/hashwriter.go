package testutils

import (
	"bytes"
	"crypto/sha512"
	"hash"
)

type CountingHashWriter struct {
	hash         hash.Hash
	bytesWritten int
}

func (h *CountingHashWriter) Write(p []byte) (n int, err error) {
	n, err = h.hash.Write(p)
	h.bytesWritten += n
	return n, err
}

func (h *CountingHashWriter) BytesWritten() int {
	return h.bytesWritten
}

func (a *CountingHashWriter) Equals(b *CountingHashWriter) bool {
	return bytes.Compare(a.hash.Sum(nil), b.hash.Sum(nil)) == 0
}

func NewCountingHashWriter() *CountingHashWriter {
	return &CountingHashWriter{
		hash:         sha512.New(),
		bytesWritten: 0,
	}
}
