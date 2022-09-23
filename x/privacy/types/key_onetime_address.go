package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OnetimeAddressKeyPrefix is the prefix to retrieve all OnetimeAddress
	OnetimeAddressKeyPrefix = "OnetimeAddress/value/"
)

// OnetimeAddressKey returns the store key to retrieve a OnetimeAddress from the index fields
func OnetimeAddressKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
