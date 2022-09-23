package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// OutputCoinKeyPrefix is the prefix to retrieve all OutputCoin
	OutputCoinKeyPrefix = "OutputCoin/value/"
)

// OutputCoinKey returns the store key to retrieve a OutputCoin from the index fields
func OutputCoinKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
