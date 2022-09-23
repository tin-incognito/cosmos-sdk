package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetOnetimeAddress set a specific onetimeAddress in the store from its index
func (k Keeper) SetOnetimeAddress(ctx sdk.Context, onetimeAddress types.OnetimeAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OnetimeAddressKeyPrefix))
	b := k.cdc.MustMarshal(&onetimeAddress)
	store.Set(types.OnetimeAddressKey(
		onetimeAddress.Index,
	), b)
}

// GetOnetimeAddress returns a onetimeAddress from its index
func (k Keeper) GetOnetimeAddress(
	ctx sdk.Context,
	index string,

) (val types.OnetimeAddress, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OnetimeAddressKeyPrefix))

	b := store.Get(types.OnetimeAddressKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOnetimeAddress removes a onetimeAddress from the store
func (k Keeper) RemoveOnetimeAddress(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OnetimeAddressKeyPrefix))
	store.Delete(types.OnetimeAddressKey(
		index,
	))
}

// GetAllOnetimeAddress returns all onetimeAddress
func (k Keeper) GetAllOnetimeAddress(ctx sdk.Context) (list []types.OnetimeAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OnetimeAddressKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OnetimeAddress
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
