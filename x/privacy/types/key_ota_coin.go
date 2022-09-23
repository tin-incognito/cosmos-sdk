package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OTACoinKeyPrefix is the prefix to retrieve all OTACoin
	OTACoinKeyPrefix = "OTACoin/value/"
)

// OTACoinKey returns the store key to retrieve a OTACoin from the index fields
func OTACoinKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
