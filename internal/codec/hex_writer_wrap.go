package codec

import (
	"fmt"
	"io"
)

type HexWriterWrapper struct {
	writes int
	inner  io.Writer
}

func NewHexWriterWrapper(inner io.Writer) *HexWriterWrapper {
	return &HexWriterWrapper{inner: inner}
}

func (h *HexWriterWrapper) Write(p []byte) (n int, err error) {
	for _, b := range p {
		if _, err := fmt.Fprintf(h.inner, "%02x", b); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}
