package models

import (
	"fmt"
	"math/big"
)

// SigPubKey defines the public key to sign ring signatures in version 2. It is an array of coin indexes.
type SigPubKey struct {
	Indexes [][]*big.Int
}

func (sigPub SigPubKey) Bytes() ([]byte, error) {
	n := len(sigPub.Indexes)
	if n == 0 {
		return nil, fmt.Errorf("txSigPublicKeyVer2.ToBytes: Indexes is empty")
	}
	if n > MaxSizeByte {
		return nil, fmt.Errorf("txSigPublicKeyVer2.ToBytes: Indexes is too large, too many rows")
	}
	m := len(sigPub.Indexes[0])
	if m > MaxSizeByte {
		return nil, fmt.Errorf("txSigPublicKeyVer2.ToBytes: Indexes is too large, too many columns")
	}
	for i := 1; i < n; i++ {
		if len(sigPub.Indexes[i]) != m {
			return nil, fmt.Errorf("txSigPublicKeyVer2.ToBytes: Indexes is not a rectangle array")
		}
	}

	b := make([]byte, 0)
	b = append(b, byte(n))
	b = append(b, byte(m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			currentByte := sigPub.Indexes[i][j].Bytes()
			lengthByte := len(currentByte)
			if lengthByte > MaxSizeByte {
				return nil, fmt.Errorf("txSigPublicKeyVer2.ToBytes: IndexesByte is too large")
			}
			b = append(b, byte(lengthByte))
			b = append(b, currentByte...)
		}
	}
	return b, nil
}
