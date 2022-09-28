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
	pk PrivacyKeeper
}

func NewValidateByDbDecorator(privacyKeeper PrivacyKeeper) ValidateByDbDecorator {
	return ValidateByDbDecorator{}
}

func (vbdd ValidateByDbDecorator) IsPrivacy() bool {
	return true
}

func (vbdd ValidateByDbDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {

	// no need to check index, has been checked before
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		proof := repos.NewPaymentProof()
		err := proof.SetBytes(msg.Proof)
		if err != nil {
			return ctx, err
		}
		inputCoins := proof.InputCoins()
		for _, item := range inputCoins {
			isConfidentialAsset := item.AssetTag != nil
			serialNum := item.GetKeyImage().ToBytesS()
			hash := common.HashH(append([]byte{common.BoolToByte(isConfidentialAsset)}, serialNum...))
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
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	return next(ctx, tx, simulate)
}
