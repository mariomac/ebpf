package codec

import (
	"bytes"
	"encoding/hex"
	"io"
	"reflect"
)

var HexReaderPackage = reflect.TypeOf(HexReader{}).PkgPath()

type HexReader struct {
	buffer io.ReaderAt
}

func NewHexReader(input []byte) (io.ReaderAt, error) {
	var buffer bytes.Buffer

	hexString := []byte{0, 0}

	// Extract only hexadecimal characters, ignoring line breaks and whitespace
	// Read from source in chunks to avoid loading everything into memory
	hd := 0
	for _, b := range input {
		if (b >= '0' && b <= '9') || (b >= 'A' && b <= 'F') || (b >= 'a' && b <= 'f') {
			hexString[hd] = b
			hd = (hd + 1) % 2
			// always assumes that the file is formed by pairs of hex digits
			if hd == 0 {
				hexBytes, err := hex.DecodeString(string(hexString))
				if err != nil {
					return nil, err
				}
				buffer.Write(hexBytes)
			}
		}
	}
	return &HexReader{buffer: bytes.NewReader(buffer.Bytes())}, nil
}

func (h *HexReader) ReadAt(p []byte, off int64) (n int, err error) {
	return h.buffer.ReadAt(p, off)
}
