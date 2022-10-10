package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type ValidateByDbDecorator struct {
	c  *Cache
	pk PrivacyKeeper
}

func NewValidateByDbDecorator(privacyKeeper PrivacyKeeper, c *Cache) ValidateByDbDecorator {
	return ValidateByDbDecorator{pk: privacyKeeper, c: c}
}

func (vbdd ValidateByDbDecorator) IsPrivacy() bool {
	return true
}

func (vbdd ValidateByDbDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isPrivate, err := tx.IsPrivacy()
	if err != nil {
		return ctx, err
	}
	if !isPrivate {
		return next(ctx, tx, simulate)
	}

	// no need to check index, has been checked before
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		isMintTx, err := IsMintTx(tx)
		if err != nil {
			return ctx, err
		}
		if !isMintTx {
			key, err := common.NewHashFromBytes(msg.Hash)
			if err != nil {
				return ctx, err
			}
			proof, err := vbdd.c.GetProof(*key)
			if err != nil {
				proof = repos.NewPaymentProof()
				if err = proof.SetBytes(msg.Proof); err != nil {
					return ctx, err
				}
				if err = vbdd.c.AddProof(*key, proof); err != nil {
					return ctx, err
				}
			} else {
				fmt.Println("err:", err)
			}
			inputCoins := proof.InputCoins()
			for _, item := range inputCoins {
				isConfidentialAsset := item.AssetTag != nil
				serialNum := item.GetKeyImage().ToBytesS()
				hash := common.HashH(append([]byte{common.BoolToByte(isConfidentialAsset)}, serialNum...))
				//fmt.Println("vbdd:", vbdd)
				//fmt.Println("vbdd.pk:", vbdd.pk)
				if _, found := vbdd.pk.GetSerialNumber(ctx, hash.String()); found {
					return ctx, fmt.Errorf("Duplicate serialNumber %s", item.GetKeyImage().String())
				}
			}
			outputCoins := proof.OutputCoins()
			for _, outCoin := range outputCoins {
				isConfidentialAsset := outCoin.AssetTag != nil
				otaPublicKey := outCoin.GetPublicKey().ToBytesS()
				outputCoinBytes := outCoin.Bytes()
				temp := append([]byte{common.BoolToByte(isConfidentialAsset)}, otaPublicKey...)
				hash := common.HashH(append(temp, outputCoinBytes...))
				if _, found := vbdd.pk.GetOnetimeAddress(ctx, hash.String()); found {
					return ctx, fmt.Errorf("Duplicate onetimeAddress at index %s", hash.String())
				}
			}

			//TODO: @tin validate metadata by db
		}
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	return next(ctx, tx, simulate)
}
