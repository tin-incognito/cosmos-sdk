package keeper

import (
	"context"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func (k msgServer) ShieldCoin(shieldCtx context.Context, shield *types.MsgShield) (*types.MsgShieldResponse, error) {
	fromAcc := shield.GetSigners()[0]
	i := types2.NewInt(int64(shield.Amount))

	ctx := types2.UnwrapSDKContext(shieldCtx)

	k.bankKeeper.SendCoinsFromAccountToModule(ctx, fromAcc, types.ModuleName, types2.Coins{}.Add(types2.Coin{"stake", i}))

	err := k.setPrivacyData(ctx, shield.GetProof())
	if err != nil {
		return nil, err
	}

	return &types.MsgShieldResponse{}, nil
}
