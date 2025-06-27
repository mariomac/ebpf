package codec

import (
	"bytes"
	"encoding/hex"
	"errors"
	"io"
	"strings"
)

type HexReader struct {
	buffer io.ReaderAt
}

func NewHexReader(source io.Reader) (*HexReader, error) {
	var buffer bytes.Buffer

	// Extract only hexadecimal characters, ignoring line breaks and whitespace
	var hexString strings.Builder

	// Read from source in chunks to avoid loading everything into memory
	chunk := make([]byte, 4096) // 4KB chunks
	for {
		n, err := source.Read(chunk)
		if n > 0 {
			// Process the chunk and extract hex characters
			for i := 0; i < n; i++ {
				b := chunk[i]
				if (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f') {
					hexString.WriteByte(b)
					if hexString.Len() == 2 {
						hexBytes, err := hex.DecodeString(hexString.String())
						if err != nil {
							return nil, err
						}
						buffer.Write(hexBytes)
						hexString.Reset()
					}
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
	}

	return &HexReader{buffer: bytes.NewReader(buffer.Bytes())}, nil
}

func (h *HexReader) ReadAt(p []byte, off int64) (n int, err error) {
	return h.buffer.ReadAt(p, off)
}
