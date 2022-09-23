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

func (k Keeper) SerialNumberAll(c context.Context, req *types.QueryAllSerialNumberRequest) (*types.QueryAllSerialNumberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var serialNumbers []types.SerialNumber
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	serialNumberStore := prefix.NewStore(store, types.KeyPrefix(types.SerialNumberKeyPrefix))

	pageRes, err := query.Paginate(serialNumberStore, req.Pagination, func(key []byte, value []byte) error {
		var serialNumber types.SerialNumber
		if err := k.cdc.Unmarshal(value, &serialNumber); err != nil {
			return err
		}

		serialNumbers = append(serialNumbers, serialNumber)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllSerialNumberResponse{SerialNumber: serialNumbers, Pagination: pageRes}, nil
}

func (k Keeper) SerialNumber(c context.Context, req *types.QueryGetSerialNumberRequest) (*types.QueryGetSerialNumberResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetSerialNumber(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetSerialNumberResponse{SerialNumber: val}, nil
}
