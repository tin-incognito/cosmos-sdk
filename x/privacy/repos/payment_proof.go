package repos

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/bulletproofs"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
)

type PaymentProof struct {
	aggregatedRangeProof *bulletproofs.AggregatedRangeProof
	inputCoins           []*coin.Coin
	outputCoins          []*coin.Coin
}

func NewPaymentProof() *PaymentProof {
	proof := &PaymentProof{}
	proof.aggregatedRangeProof = bulletproofs.NewAggregatedRangeProof()
	proof.inputCoins = []*coin.Coin{}
	proof.outputCoins = []*coin.Coin{}
	return proof
}

// Bytes does byte serialization for this payment proof
func (proof PaymentProof) Bytes() []byte {
	var bytes []byte

	comOutputMultiRangeProof := proof.aggregatedRangeProof.Bytes()
	var rangeProofLength uint32 = uint32(len(comOutputMultiRangeProof))
	bytes = append(bytes, common.Uint32ToBytes(rangeProofLength)...)
	bytes = append(bytes, comOutputMultiRangeProof...)

	// InputCoins
	bytes = append(bytes, byte(len(proof.inputCoins)))
	for i := 0; i < len(proof.inputCoins); i++ {
		inputCoins := proof.inputCoins[i].Bytes()
		lenInputCoins := len(inputCoins)
		var lenInputCoinsBytes []byte
		if lenInputCoins < 256 {
			lenInputCoinsBytes = []byte{byte(lenInputCoins)}
		} else {
			lenInputCoinsBytes = common.IntToBytes(lenInputCoins)
		}

		bytes = append(bytes, lenInputCoinsBytes...)
		bytes = append(bytes, inputCoins...)
	}

	// OutputCoins
	bytes = append(bytes, byte(len(proof.outputCoins)))
	for i := 0; i < len(proof.outputCoins); i++ {
		outputCoins := proof.outputCoins[i].Bytes()
		lenOutputCoins := len(outputCoins)
		var lenOutputCoinsBytes []byte
		if lenOutputCoins < 256 {
			lenOutputCoinsBytes = []byte{byte(lenOutputCoins)}
		} else {
			lenOutputCoinsBytes = common.IntToBytes(lenOutputCoins)
		}

		bytes = append(bytes, lenOutputCoinsBytes...)
		bytes = append(bytes, outputCoins...)
	}

	return bytes
}

// SetBytes does byte deserialization for this payment proof
func (proof *PaymentProof) SetBytes(proofbytes []byte) error {
	if len(proofbytes) == 0 {
		return fmt.Errorf("Proof bytes is zero")
	}
	offset := 0

	// ComOutputMultiRangeProofSize *aggregatedRangeProof
	if offset+common.Uint32Size >= len(proofbytes) {
		return fmt.Errorf("Out of range aggregated range proof")
	}
	lenComOutputMultiRangeUint32, _ := common.BytesToUint32(proofbytes[offset : offset+common.Uint32Size])
	lenComOutputMultiRangeProof := int(lenComOutputMultiRangeUint32)
	offset += common.Uint32Size

	if offset+lenComOutputMultiRangeProof > len(proofbytes) {
		return fmt.Errorf("Out of range aggregated range proof")
	}
	if lenComOutputMultiRangeProof > 0 {
		bulletproof := bulletproofs.NewAggregatedRangeProof()
		proof.aggregatedRangeProof = bulletproof
		err := proof.aggregatedRangeProof.SetBytes(proofbytes[offset : offset+lenComOutputMultiRangeProof])
		if err != nil {
			return err
		}
		offset += lenComOutputMultiRangeProof
	}

	if offset >= len(proofbytes) {
		return fmt.Errorf("Out of range input coins")
	}
	lenInputCoinsArray := int(proofbytes[offset])
	offset++
	proof.inputCoins = make([]*coin.Coin, lenInputCoinsArray)
	var err error
	for i := 0; i < lenInputCoinsArray; i++ {
		// try get 1-byte for len
		if offset >= len(proofbytes) {
			return fmt.Errorf("Out of range input coins")
		}
		lenInputCoin := int(proofbytes[offset])
		offset++

		if offset+lenInputCoin > len(proofbytes) {
			return fmt.Errorf("Out of range input coins")
		}
		proof.inputCoins[i], err = coin.NewCoinFromBytes(proofbytes[offset : offset+lenInputCoin])
		if err != nil {
			// 1-byte is wrong
			// try get 2-byte for len
			if offset+1 > len(proofbytes) {
				return fmt.Errorf("Out of range input coins")
			}
			lenInputCoin = common.BytesToInt(proofbytes[offset-1 : offset+1])
			offset++

			if offset+lenInputCoin > len(proofbytes) {
				return fmt.Errorf("Out of range input coins")
			}
			proof.inputCoins[i], err = coin.NewCoinFromBytes(proofbytes[offset : offset+lenInputCoin])
			if err != nil {
				return err
			}
		}
		offset += lenInputCoin
	}

	if offset >= len(proofbytes) {
		return fmt.Errorf("Out of range output coins")
	}
	lenOutputCoinsArray := int(proofbytes[offset])
	offset++
	proof.outputCoins = make([]*coin.Coin, lenOutputCoinsArray)
	for i := 0; i < lenOutputCoinsArray; i++ {
		proof.outputCoins[i] = new(coin.Coin)
		// try get 1-byte for len
		if offset >= len(proofbytes) {
			return fmt.Errorf("Out of range output coins")
		}
		lenOutputCoin := int(proofbytes[offset])
		offset++

		if offset+lenOutputCoin > len(proofbytes) {
			return fmt.Errorf("Out of range output coins")
		}
		err := proof.outputCoins[i].SetBytes(proofbytes[offset : offset+lenOutputCoin])
		if err != nil {
			// 1-byte is wrong
			// try get 2-byte for len
			if offset+1 > len(proofbytes) {
				return fmt.Errorf("Out of range output coins")
			}
			lenOutputCoin = common.BytesToInt(proofbytes[offset-1 : offset+1])
			offset++

			if offset+lenOutputCoin > len(proofbytes) {
				return fmt.Errorf("Out of range output coins")
			}
			e := proof.outputCoins[i].SetBytes(proofbytes[offset : offset+lenOutputCoin])
			if e != nil {
				return e
			}
		}
		offset += lenOutputCoin
	}

	return nil
}

func (p *PaymentProof) SetInputCoins(inputCoins []*coin.Coin) error {
	var err error
	p.inputCoins = make([]*coin.Coin, len(inputCoins))
	for i := 0; i < len(inputCoins); i++ {
		b := inputCoins[i].Bytes()
		if p.inputCoins[i], err = coin.NewCoinFromBytes(b); err != nil {
			return err
		}
	}
	return err
}

func (p *PaymentProof) SetOutputCoins(outputCoins []*coin.Coin) error {
	var err error
	p.outputCoins = make([]*coin.Coin, len(outputCoins))
	for i := 0; i < len(outputCoins); i++ {
		b := outputCoins[i].Bytes()
		if p.outputCoins[i], err = coin.NewCoinFromBytes(b); err != nil {
			return err
		}
	}
	return err
}

func (p *PaymentProof) InputCoins() []coin.Coin {
	res := make([]coin.Coin, len(p.inputCoins))
	for i, v := range p.inputCoins {
		res[i] = *v
	}
	return res
}

func (p *PaymentProof) OutputCoins() []coin.Coin {
	res := make([]coin.Coin, len(p.outputCoins))
	for i, v := range p.outputCoins {
		res[i] = *v
	}
	return res
}

func (p *PaymentProof) SetInputCoinAtIndex(index int, c *coin.Coin) error {
	if index >= len(p.inputCoins) {
		return fmt.Errorf("Index out of range")
	}
	p.inputCoins[index] = new(coin.Coin)
	*p.inputCoins[index] = *c
	return nil
}

func (p *PaymentProof) SetOutputCoinAtIndex(index int, c *coin.Coin) error {
	if index >= len(p.outputCoins) {
		return fmt.Errorf("Index out of range")
	}
	p.outputCoins[index] = new(coin.Coin)
	*p.outputCoins[index] = *c
	return nil
}

func (p *PaymentProof) SetAggregatedRangeProof(proof *bulletproofs.AggregatedRangeProof) {
	p.aggregatedRangeProof = proof
}

func (p *PaymentProof) ValidateSanity() (bool, error) {
	if len(p.inputCoins) > 255 {
		return false, fmt.Errorf("Input coins in tx are very large:" + strconv.Itoa(len(p.inputCoins)))
	}

	if len(p.outputCoins) > 255 {
		return false, fmt.Errorf("Output coins in tx are very large:" + strconv.Itoa(len(p.outputCoins)))
	}

	if !p.aggregatedRangeProof.ValidateSanity() {
		return false, fmt.Errorf("validate sanity Aggregated range proof failed")
	}

	// check output coins with privacy
	duplicatePublicKeys := make(map[string]bool)
	outputCoins := p.outputCoins
	// cmsValues := proof.aggregatedRangeProof.GetCommitments()
	for _, outputCoin := range outputCoins {
		if outputCoin.GetPublicKey() == nil || !outputCoin.GetPublicKey().PointValid() {
			return false, fmt.Errorf("validate sanity Public key of output coin failed")
		}

		// check duplicate output addresses
		pubkeyStr := string(outputCoin.GetPublicKey().ToBytesS())
		if _, ok := duplicatePublicKeys[pubkeyStr]; ok {
			return false, fmt.Errorf("Cannot have duplicate publickey ")
		}
		duplicatePublicKeys[pubkeyStr] = true

		if !outputCoin.GetCommitment().PointValid() {
			return false, fmt.Errorf("validate sanity Coin commitment of output coin failed")
		}

		/*// re-compute the commitment if the output coin's address is the burning address*/
		/*// burn TX cannot use confidential asset]*/
		/*// BOOKMARK*/
		/*if common.IsPublicKeyBurningAddress(outputCoins[i].GetPublicKey().ToBytesS()) {*/
		/*value := outputCoin.GetValue()*/
		/*rand := outputCoin.GetRandomness()*/
		/*commitment := operation.PedCom.CommitAtIndex(new(operation.Scalar).FromUint64(value), rand, coin.PedersenValueIndex)*/
		/*outputCoinSpecific, ok := outputCoin.(*coin.CoinV2)*/
		/*if !ok {*/
		/*return false, errors.New("Validate sanity - Cannot cast a coin to v2")*/
		/*}*/
		/*if outputCoinSpecific.GetAssetTag() != nil {*/
		/*com, err := outputCoinSpecific.ComputeCommitmentCA()*/
		/*if err != nil {*/
		/*return false, errors.New("Cannot compute commitment for confidential asset")*/
		/*}*/
		/*commitment = com*/
		/*}*/
		/*if !operation.IsPointEqual(commitment, outputCoin.GetCommitment()) {*/
		/*return false, errors.New("validate sanity Coin commitment of burned coin failed")*/
		/*}*/
		/*}*/
	}
	return true, nil
}

func (p *PaymentProof) Verify() (bool, error) {
	inputCoins := p.inputCoins
	dupMap := make(map[string]bool)
	for _, inCoin := range inputCoins {
		identifier := base64.StdEncoding.EncodeToString(inCoin.GetKeyImage().ToBytesS())
		_, exists := dupMap[identifier]
		if exists {
			return false, fmt.Errorf("Duplicate input inCoin in PaymentProofV2")
		}
		dupMap[identifier] = true
	}
	return true, nil

	//return proof.verifyHasConfidentialAsset(isBatch)
}
