package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func (k msgServer) PrivacyData(goCtx context.Context, msg *types.MsgPrivacyData) (*types.MsgPrivacyDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	err := k.setPrivacyData(ctx, msg.Proof)
	if err != nil {
		return nil, err
	}

	if msg.TxType == models.TxUnshieldType {
		unshieldData := &types.MsgUnShield{}
		unshieldData.Unmarshal(msg.Metadata)
		toAccount, err := sdk.AccAddressFromBech32(unshieldData.ToAdrress)
		if err != nil {
			return nil, err
		}
		i := sdk.NewInt(int64(unshieldData.Amount))
		err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAccount, sdk.Coins{}.Add(sdk.Coin{"prv", i}))
		if err != nil {
			return nil, err
		}
	}
	return &types.MsgPrivacyDataResponse{}, nil
}
