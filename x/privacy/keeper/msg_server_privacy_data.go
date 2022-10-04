package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func (k msgServer) PrivacyData(goCtx context.Context, msg *types.MsgPrivacyData) (*types.MsgPrivacyDataResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.Metadata) != 0 {

	}

	err := k.setPrivacyData(ctx, msg.Proof)
	if err != nil {
		return nil, err
	}

	return &types.MsgPrivacyDataResponse{}, nil
}
