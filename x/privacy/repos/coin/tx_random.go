package coin

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

type TxRandom [TxRandomGroupSize]byte

func NewTxRandom() *TxRandom {
	txRandom := new(operation.Point).Identity()
	index := uint32(0)

	res := new(TxRandom)
	res.SetTxConcealRandomPoint(txRandom)
	res.SetIndex(index)
	return res
}

func (t TxRandom) GetTxConcealRandomPoint() (*operation.Point, error) {
	return new(operation.Point).FromBytesS(t[operation.Ed25519KeySize+4:])
}

func (t TxRandom) GetTxOTARandomPoint() (*operation.Point, error) {
	return new(operation.Point).FromBytesS(t[:operation.Ed25519KeySize])
}

func (t TxRandom) GetIndex() (uint32, error) {
	return common.BytesToUint32(t[operation.Ed25519KeySize : operation.Ed25519KeySize+4])
}

func (t *TxRandom) SetTxConcealRandomPoint(txConcealRandom *operation.Point) {
	txRandomBytes := txConcealRandom.ToBytesS()
	copy(t[operation.Ed25519KeySize+4:], txRandomBytes)
}

func (t *TxRandom) SetTxOTARandomPoint(txRandom *operation.Point) {
	txRandomBytes := txRandom.ToBytesS()
	copy(t[:operation.Ed25519KeySize], txRandomBytes)
}

func (t *TxRandom) SetIndex(index uint32) {
	indexBytes := common.Uint32ToBytes(index)
	copy(t[operation.Ed25519KeySize:], indexBytes)
}

func (t TxRandom) Bytes() []byte {
	return t[:]
}

func (t *TxRandom) SetBytes(b []byte) error {
	if b == nil || len(b) != TxRandomGroupSize {
		return fmt.Errorf("cannot SetByte to TxRandom. Input is invalid")
	}
	_, err := new(operation.Point).FromBytesS(b[:operation.Ed25519KeySize])
	if err != nil {
		return fmt.Errorf("cannot TxRandomGroupSize.SetBytes: bytes is not operation.Point err: %v", err)
	}
	_, err = new(operation.Point).FromBytesS(b[operation.Ed25519KeySize+4:])
	if err != nil {
		return fmt.Errorf("cannot TxRandomGroupSize.SetBytes: bytes is not operation.Point err: %v", err)
	}
	copy(t[:], b)
	return nil
}

// MarshalText returns the hex string representation of `p` but as bytes
func (t TxRandom) MarshalText() []byte {
	return []byte(hex.EncodeToString(t.Bytes()))
}

// UnmarshalText decodes a Point from its hex string form and sets `p` to it, then returns p
func (t *TxRandom) UnmarshalText(data []byte) (*TxRandom, error) {
	byteSlice, _ := hex.DecodeString(string(data))
	if len(byteSlice) != TxRandomGroupSize {
		return nil, fmt.Errorf("invalid point byte size")
	}
	err := t.SetBytes(byteSlice)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func (t *TxRandom) Marshal() ([]byte, error) {
	return []byte(hex.EncodeToString(t.Bytes())), nil
}

func (t *TxRandom) Unmarshal(data []byte) error {
	byteSlice, _ := hex.DecodeString(string(data))
	if len(byteSlice) != TxRandomGroupSize {
		return fmt.Errorf("invalid point byte size")
	}
	err := t.SetBytes(byteSlice)
	if err != nil {
		return err
	}
	return nil
}
