package coin

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

func getMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parsePointForSetBytes(coinBytes *[]byte, offset *int) (*operation.Point, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return nil, fmt.Errorf("offset is larger than len(bytes), cannot parse point")
	}
	var point *operation.Point
	var err error
	lenField := b[*offset]
	*offset++
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return nil, fmt.Errorf("offset+curLen is larger than len(bytes), cannot parse point for set bytes")
		}
		data := b[*offset : *offset+int(lenField)]
		point, err = new(operation.Point).FromBytesS(data)
		if err != nil {
			return nil, err
		}
		*offset += int(lenField)
	}
	return point, nil
}

func parseInfoForSetBytes(coinBytes *[]byte, offset *int) ([]byte, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return []byte{}, fmt.Errorf("offset is larger than len(bytes), cannot parse info")
	}
	info := []byte{}
	lenField := b[*offset]
	*offset++
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return []byte{}, fmt.Errorf("offset+curLen is larger than len(bytes), cannot parse info for set bytes")
		}
		info = make([]byte, lenField)
		copy(info, b[*offset:*offset+int(lenField)])
		*offset += int(lenField)
	}
	return info, nil
}

func parseScalarForSetBytes(coinBytes *[]byte, offset *int) (*operation.Scalar, error) {
	b := *coinBytes
	if *offset >= len(b) {
		return nil, fmt.Errorf("offset is larger than len(bytes), cannot parse scalar")
	}
	var sc *operation.Scalar
	lenField := b[*offset]
	*offset++
	if lenField != 0 {
		if *offset+int(lenField) > len(b) {
			return nil, fmt.Errorf("offset+curLen is larger than len(bytes), cannot parse scalar for set bytes")
		}
		data := b[*offset : *offset+int(lenField)]
		sc = new(operation.Scalar).FromBytesS(data)
		*offset += int(lenField)
	}
	return sc, nil
}

func NewCoinFromBytes(b []byte) (*Coin, error) {
	c := NewCoin()
	err := c.SetBytes(b)
	return c, err
}

func NewCoinFromAmountAndTxRandomBytes(
	amount uint64, publicKey *operation.Point, txRandom *TxRandom, info []byte,
) *Coin {
	c := NewCoin()
	c.SetPublicKey(publicKey)
	c.SetAmount(new(operation.Scalar).FromUint64(amount))
	c.SetRandomness(operation.RandomScalar())
	c.SetTxRandom(txRandom)
	c.SetCommitment(operation.PedCom.CommitAtIndex(c.GetAmount(), c.GetRandomness(), operation.PedersenValueIndex))
	c.SetSharedRandom(nil)
	c.SetInfo(info)
	return c
}

func NewCoinFromPaymentInfo(paymentInfo *key.PaymentInfo) (*Coin, error) {
	c := NewCoin()
	// Amount, Randomness, SharedRandom are transparency until we call concealData
	c.SetAmount(new(operation.Scalar).FromUint64(paymentInfo.Amount))
	c.SetRandomness(operation.RandomScalar())
	c.SetSharedRandom(operation.RandomScalar())        // shared randomness for creating one-time-address
	c.SetSharedConcealRandom(operation.RandomScalar()) // shared randomness for concealing amount and blinding asset tag
	c.SetInfo(paymentInfo.Message)
	c.SetCommitment(operation.PedCom.CommitAtIndex(c.GetAmount(), c.GetRandomness(), operation.PedersenValueIndex))

	/*// If this is going to burning address then dont need to create ota*/
	/*if common.IsPublicKeyBurningAddress(p.PaymentAddress.Pk) {*/
	/*publicKey, err := new(operation.Point).FromBytesS(p.PaymentAddress.Pk)*/
	/*if err != nil {*/
	/*panic("Something is wrong with info.paymentAddress.pk, burning address should be a valid point")*/
	/*}*/
	/*c.SetPublicKey(publicKey)*/
	/*return c, nil*/
	/*}*/

	index := uint32(0)
	publicOTA := paymentInfo.PaymentAddress.GetOTAPublicKey()
	if publicOTA == nil {
		return nil, fmt.Errorf("public OTA from payment address is nil")
	}
	publicSpend := paymentInfo.PaymentAddress.GetPublicSpend()
	rK := new(operation.Point).ScalarMult(publicOTA, c.GetSharedRandom())

	// Get publickey
	hash := operation.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
	HrKG := new(operation.Point).ScalarMultBase(hash)
	publicKey := new(operation.Point).Add(HrKG, publicSpend)
	c.SetPublicKey(publicKey)
	otaRandomPoint := new(operation.Point).ScalarMultBase(c.GetSharedRandom())
	concealRandomPoint := new(operation.Point).ScalarMultBase(c.GetSharedConcealRandom())
	c.SetTxRandomDetail(concealRandomPoint, otaRandomPoint, index)
	return c, nil
}
