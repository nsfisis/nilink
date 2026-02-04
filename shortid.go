package main

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
)

var b32 = base32.StdEncoding.WithPadding(base32.NoPadding)

func encodeID(id int64) (string, error) {
	if id < 0 {
		return "", fmt.Errorf("id out of range: %d", id)
	}
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	// Trim leading zero bytes, keep at least 2 bytes (4 chars).
	i := 0
	for i < 6 && buf[i] == 0 {
		i++
	}
	return b32.EncodeToString(buf[i:]), nil
}

func decodeID(s string) (int64, error) {
	s = strings.ToUpper(s)
	buf, err := b32.DecodeString(s)
	if err != nil {
		return 0, fmt.Errorf("invalid short id: %w", err)
	}
	if len(buf) == 0 || len(buf) > 8 {
		return 0, fmt.Errorf("invalid short id")
	}
	padded := make([]byte, 8)
	copy(padded[8-len(buf):], buf)
	return int64(binary.BigEndian.Uint64(padded)), nil
}
