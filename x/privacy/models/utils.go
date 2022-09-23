package models

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

func ArrayScalarToBytes(arr *[]*operation.Scalar) ([]byte, error) {
	scalarArr := *arr

	n := len(scalarArr)
	if n > 255 {
		return nil, fmt.Errorf("arrayScalarToBytes: length of scalar array is too big")
	}
	b := make([]byte, 1)
	b[0] = byte(n)

	for _, sc := range scalarArr {
		b = append(b, sc.ToBytesS()...)
	}
	return b, nil
}

func CalculateSumOutputsWithFee(outputCoins []*coin.Coin, fee uint64) *operation.Point {
	sumOutputsWithFee := new(operation.Point).Identity()
	for i := 0; i < len(outputCoins); i++ {
		sumOutputsWithFee.Add(sumOutputsWithFee, outputCoins[i].GetCommitment())
	}
	feeCommitment := new(operation.Point).ScalarMult(
		operation.PedCom.G[operation.PedersenValueIndex],
		new(operation.Scalar).FromUint64(fee),
	)
	sumOutputsWithFee.Add(sumOutputsWithFee, feeCommitment)
	return sumOutputsWithFee
}

func DebugCoins(coins []*coin.Coin) {
	for _, v := range coins {
		fmt.Println("commitment:", v.GetCommitment().String())
	}
}

func DebugCoins1(coins []coin.Coin) {
	for _, v := range coins {
		fmt.Println("commitment:", v.GetCommitment().String())
	}
}
