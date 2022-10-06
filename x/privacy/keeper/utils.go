package keeper

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func (k Keeper) setPrivacyData(ctx sdk.Context, proof []byte) error {
	outputCoinLength := big.NewInt(0)
	outputCoinSerialNumber, ok := k.GetOutputCoinLength(ctx)
	if ok {
		outputCoinLength.SetBytes(outputCoinSerialNumber.Value)
	}

	// parse data
	serialNumbers, commitments, outputCoins, otaCoins, onetimeAddresses, newOutputCoinLength, err := models.FetchDataFromTx(ctx, proof, *outputCoinLength)
	if err != nil {
		return err
	}

	// store serialNumber
	for _, serialNumber := range serialNumbers {
		k.SetSerialNumber(ctx, serialNumber)
	}

	// store outputCoin
	for key, outputCoin := range outputCoins {
		otas, found := onetimeAddresses[key]
		if !found {
			return fmt.Errorf("Cannot find list otas with key %v", key)
		}
		otaCoin, found := otaCoins[key]
		if !found {
			return fmt.Errorf("Cannot find list otaCoin with key %v", key)
		}
		for i, o := range outputCoin {
			k.SetOutputCoin(ctx, o)
			if i >= len(otaCoin) {
				return fmt.Errorf("Cannot find otaCoin with key %v and index %v", key, i)
			}
			oa := otaCoin[i]
			if _, ok := k.GetOTACoin(ctx, oa.Index); ok {
				return fmt.Errorf("Duplicate ota coin")
			}
			k.SetOTACoin(ctx, oa)
			ota := otas[i]
			if _, ok := k.GetOnetimeAddress(ctx, ota.Index); ok {
				return fmt.Errorf("Duplicate ota")
			}
			k.SetOnetimeAddress(ctx, ota)
		}
	}

	// store commitment
	for _, commitment := range commitments {
		for _, c := range commitment {
			if _, ok := k.GetCommitment(ctx, c.Index); ok {
				return fmt.Errorf("Duplicate commitment")
			}
			k.SetCommitment(ctx, c)
		}
	}

	// store output coin length
	k.SetOutputCoinLength(ctx, types.OutputCoinLength{Value: newOutputCoinLength.Bytes()})

	return nil
}
