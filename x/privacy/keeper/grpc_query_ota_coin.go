package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OTACoinAll(c context.Context, req *types.QueryAllOTACoinRequest) (*types.QueryAllOTACoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var oTACoins []types.OTACoin
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	oTACoinStore := prefix.NewStore(store, types.KeyPrefix(types.OTACoinKeyPrefix))

	pageRes, err := query.Paginate(oTACoinStore, req.Pagination, func(key []byte, value []byte) error {
		var oTACoin types.OTACoin
		if err := k.cdc.Unmarshal(value, &oTACoin); err != nil {
			return err
		}

		oTACoins = append(oTACoins, oTACoin)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOTACoinResponse{OTACoin: oTACoins, Pagination: pageRes}, nil
}

func (k Keeper) OTACoin(c context.Context, req *types.QueryGetOTACoinRequest) (*types.QueryGetOTACoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOTACoin(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOTACoinResponse{OTACoin: val}, nil
}
