package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// SerialNumberKeyPrefix is the prefix to retrieve all SerialNumber
	SerialNumberKeyPrefix = "SerialNumber/value/"
)

// SerialNumberKey returns the store key to retrieve a SerialNumber from the index fields
func SerialNumberKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
