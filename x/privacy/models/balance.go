package models

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func BalanceByKeySet(coins []types.OutputCoin, keySet key.KeySet, serialNumbers map[string]types.SerialNumber) (uint64, error) {
	var res uint64
	for _, outputCoin := range coins {
		o := &coin.Coin{}
		err := o.SetBytes(outputCoin.Value)
		if err != nil {
			return 0, err
		}
		if ok, _ := o.DoesCoinBelongToKeySet(&keySet); !ok {
			continue
		}
		o, err = o.Decrypt(&keySet)
		if err != nil {
			return 0, err
		}
		serialNum := o.GetKeyImage().ToBytesS()
		isConfidentialAsset := o.AssetTag != nil
		hash := common.HashH(append([]byte{common.BoolToByte(isConfidentialAsset)}, serialNum...))
		if _, found := serialNumbers[hash.String()]; !found {
			temp := res + o.GetValue()
			if temp < res {
				return 0, fmt.Errorf("balance is out of range")
			}
			res = temp
		}
	}
	return res, nil
}
