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

func (k Keeper) OutputCoinAll(c context.Context, req *types.QueryAllOutputCoinRequest) (*types.QueryAllOutputCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var outputCoins []types.OutputCoin
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	outputCoinStore := prefix.NewStore(store, types.KeyPrefix(types.OutputCoinKeyPrefix))

	pageRes, err := query.Paginate(outputCoinStore, req.Pagination, func(key []byte, value []byte) error {
		var outputCoin types.OutputCoin
		if err := k.cdc.Unmarshal(value, &outputCoin); err != nil {
			return err
		}

		outputCoins = append(outputCoins, outputCoin)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOutputCoinResponse{OutputCoin: outputCoins, Pagination: pageRes}, nil
}

func (k Keeper) OutputCoin(c context.Context, req *types.QueryGetOutputCoinRequest) (*types.QueryGetOutputCoinResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOutputCoin(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOutputCoinResponse{OutputCoin: val}, nil
}
