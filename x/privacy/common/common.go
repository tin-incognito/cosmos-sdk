package common

import (
	"crypto/rand"

	"encoding/binary"
	"fmt"
	"math/big"
)

// AddPaddingBigInt adds padding to big int to it is fixed size
// and returns bytes array
func AddPaddingBigInt(numInt *big.Int, fixedSize int) []byte {
	numBytes := numInt.Bytes()
	lenNumBytes := len(numBytes)
	zeroBytes := make([]byte, fixedSize-lenNumBytes)
	numBytes = append(zeroBytes, numBytes...)
	return numBytes
}

// BytesToUint32 converts big endian 4-byte array to uint32 number
func BytesToUint32(b []byte) (uint32, error) {
	if len(b) != Uint32Size {
		return 0, fmt.Errorf("invalid length of input BytesToUint32")
	}
	return binary.BigEndian.Uint32(b), nil
}

// Uint32ToBytes converts uint32 number to big endian 4-byte array
func Uint32ToBytes(value uint32) []byte {
	b := make([]byte, Uint32Size)
	binary.BigEndian.PutUint32(b, value)
	return b
}

// IntToBytes converts an integer number to 2-byte array in big endian
func IntToBytes(n int) []byte {
	if n == 0 {
		return []byte{0, 0}
	}

	a := big.NewInt(int64(n))

	if len(a.Bytes()) > 2 {
		return []byte{}
	}

	if len(a.Bytes()) == 1 {
		return []byte{0, a.Bytes()[0]}
	}

	return a.Bytes()
}

// BytesToInt reverts an integer number from 2-byte array
func BytesToInt(bytesArr []byte) int {
	if len(bytesArr) != 2 {
		return 0
	}

	numInt := new(big.Int).SetBytes(bytesArr)
	return int(numInt.Int64())
}

// BoolToByte receives a value in bool
// and returns a value in byte
func BoolToByte(value bool) byte {
	var bitSetVar byte
	if value {
		bitSetVar = 1
	}
	return bitSetVar
}

// RandBytes generates random bytes with length
func RandBytes(length int) []byte {
	rbytes := make([]byte, length)
	rand.Read(rbytes)
	return rbytes
}

// RandBigIntMaxRange generates a big int with maximum value
func RandBigIntMaxRange(max *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, max)
}
