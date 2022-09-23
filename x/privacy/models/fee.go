package models

import (
	"math"

	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

func EstimateFee(estimateFeeCoinPerKb uint64, numInputCoins, numPayments int, metadata []byte) uint64 {
	if estimateFeeCoinPerKb == 0 || estimateFeeCoinPerKb < DefaultFeePerKb {
		estimateFeeCoinPerKb = DefaultFeePerKb
	}
	estimateTxSizeInKb := estimateTxSize(numInputCoins, numPayments, metadata)
	return uint64(estimateFeeCoinPerKb) * uint64(estimateTxSizeInKb)
}

func toB64Len(numOfBytes uint64) uint64 {
	l := (numOfBytes*4 + 2) / 3
	l = ((l + 3) / 4) * 4
	return l
}

func estimateProofSize(numIn, numOut uint64) uint64 {
	coinSizeBound := uint64(257) + (operation.Ed25519KeySize+1)*7 + coin.TxRandomGroupSize + 1
	ipProofLRLen := uint64(math.Log2(float64(numOut))) + 1
	aggProofSizeBound := uint64(4) + 1 + operation.Ed25519KeySize*uint64(7+numOut) + 1 + uint64(2*ipProofLRLen+3)*operation.Ed25519KeySize
	// add 10 for rounding
	result := uint64(1) + (coinSizeBound+1)*uint64(numIn+numOut) + 2 + aggProofSizeBound + 10
	return toB64Len(result)
}

func estimateTxSize(numInputCoins, numPayments int, metadata []byte) uint64 {
	jsonKeysSizeBound := uint64(20*10 + 2)
	sizeVersion := uint64(1)      // int8
	sizeType := uint64(5)         // string, max : 5
	sizeLockTime := uint64(8) * 3 // int64
	sizeFee := uint64(8) * 3      // uint64
	sizeInfo := toB64Len(uint64(512))

	numIn := uint64(numInputCoins)
	numOut := uint64(numPayments)

	sizeSigPubKey := uint64(numIn)*RingSize*9 + 2
	sizeSigPubKey = toB64Len(sizeSigPubKey)
	sizeSig := uint64(1) + numIn + (numIn+2)*RingSize
	sizeSig = sizeSig*33 + 3

	sizeProof := estimateProofSize(numIn, numOut)

	sizePubKeyLastByte := uint64(1) * 3
	sizeMetadata := uint64(0)
	if len(metadata) != 0 {
		//sizeMetadata += estimateTxSizeParam.metadata.CalculateSize()
	}

	sizeTx := jsonKeysSizeBound + sizeVersion + sizeType + sizeLockTime + sizeFee + sizeInfo + sizeSigPubKey + sizeSig + sizeProof + sizePubKeyLastByte + sizeMetadata
	/*if estimateTxSizeParam.privacyCustomTokenParams != nil {*/
	/*tokenKeysSizeBound := uint64(20*8 + 2)*/
	/*tokenSize := toB64Len(uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyID)))*/
	/*tokenSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertySymbol))*/
	/*tokenSize += uint64(len(estimateTxSizeParam.privacyCustomTokenParams.PropertyName))*/
	/*tokenSize += 2*/
	/*numIn = uint64(len(estimateTxSizeParam.privacyCustomTokenParams.TokenInput))*/
	/*numOut = uint64(len(estimateTxSizeParam.privacyCustomTokenParams.Receiver))*/

	/*// shadow variable names*/
	/*sizeSigPubKey := uint64(numIn)*privacy.RingSize*9 + 2*/
	/*sizeSigPubKey = toB64Len(sizeSigPubKey)*/
	/*sizeSig := uint64(1) + numIn + (numIn+2)*privacy.RingSize*/
	/*sizeSig = sizeSig*33 + 3*/

	/*sizeProof := EstimateProofSizeV2(numIn, numOut)*/
	/*tokenSize += tokenKeysSizeBound + sizeSigPubKey + sizeSig + sizeProof*/
	/*sizeTx += tokenSize*/
	/*}*/
	return sizeTx
}
