package gen

import (
	"fmt"
	"strings"
)

type BinaryFormat int

const (
	BinaryFormatRaw BinaryFormat = iota
	BinaryFormatHex
)

func ReadBinFormat(name string) (BinaryFormat, error) {
	switch strings.ToLower(name) {
	case "", "raw":
		return BinaryFormatRaw, nil
	case "hex":
		return BinaryFormatHex, nil
	}
	return 0, fmt.Errorf("unknown binary format: %s (valid values: raw, hex)", name)
}
