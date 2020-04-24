package testutils

import "io"

// Custom NopCloser returning an io.WriteCloser as io.Util only supports io.ReadCloser
type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

func NopCloser(r io.Writer) io.WriteCloser {
	return nopCloser{r}
}
