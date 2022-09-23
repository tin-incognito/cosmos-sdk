package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// CommitmentKeyPrefix is the prefix to retrieve all Commitment
	CommitmentKeyPrefix = "Commitment/value/"
)

// CommitmentKey returns the store key to retrieve a Commitment from the index fields
func CommitmentKey(
	index string,
) []byte {
	var key []byte

	indexBytes := []byte(index)
	key = append(key, indexBytes...)
	key = append(key, []byte("/")...)

	return key
}
