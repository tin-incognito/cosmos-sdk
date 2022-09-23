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

func (k Keeper) OnetimeAddressAll(c context.Context, req *types.QueryAllOnetimeAddressRequest) (*types.QueryAllOnetimeAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var onetimeAddresss []types.OnetimeAddress
	ctx := sdk.UnwrapSDKContext(c)

	store := ctx.KVStore(k.storeKey)
	onetimeAddressStore := prefix.NewStore(store, types.KeyPrefix(types.OnetimeAddressKeyPrefix))

	pageRes, err := query.Paginate(onetimeAddressStore, req.Pagination, func(key []byte, value []byte) error {
		var onetimeAddress types.OnetimeAddress
		if err := k.cdc.Unmarshal(value, &onetimeAddress); err != nil {
			return err
		}

		onetimeAddresss = append(onetimeAddresss, onetimeAddress)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllOnetimeAddressResponse{OnetimeAddress: onetimeAddresss, Pagination: pageRes}, nil
}

func (k Keeper) OnetimeAddress(c context.Context, req *types.QueryGetOnetimeAddressRequest) (*types.QueryGetOnetimeAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	val, found := k.GetOnetimeAddress(
		ctx,
		req.Index,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetOnetimeAddressResponse{OnetimeAddress: val}, nil
}
