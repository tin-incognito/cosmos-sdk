package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetOTACoin set a specific oTACoin in the store from its index
func (k Keeper) SetOTACoin(ctx sdk.Context, oTACoin types.OTACoin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OTACoinKeyPrefix))
	b := k.cdc.MustMarshal(&oTACoin)
	store.Set(types.OTACoinKey(
		oTACoin.Index,
	), b)
}

// GetOTACoin returns a oTACoin from its index
func (k Keeper) GetOTACoin(
	ctx sdk.Context,
	index string,

) (val types.OTACoin, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OTACoinKeyPrefix))

	b := store.Get(types.OTACoinKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOTACoin removes a oTACoin from the store
func (k Keeper) RemoveOTACoin(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OTACoinKeyPrefix))
	store.Delete(types.OTACoinKey(
		index,
	))
}

// GetAllOTACoin returns all oTACoin
func (k Keeper) GetAllOTACoin(ctx sdk.Context) (list []types.OTACoin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OTACoinKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OTACoin
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
