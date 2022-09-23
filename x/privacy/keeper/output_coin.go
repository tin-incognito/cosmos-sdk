package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetOutputCoin set a specific outputCoin in the store from its index
func (k Keeper) SetOutputCoin(ctx sdk.Context, outputCoin types.OutputCoin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinKeyPrefix))
	b := k.cdc.MustMarshal(&outputCoin)
	store.Set(types.OutputCoinKey(
		outputCoin.Index,
	), b)
}

// GetOutputCoin returns a outputCoin from its index
func (k Keeper) GetOutputCoin(
	ctx sdk.Context,
	index string,

) (val types.OutputCoin, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinKeyPrefix))

	b := store.Get(types.OutputCoinKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOutputCoin removes a outputCoin from the store
func (k Keeper) RemoveOutputCoin(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinKeyPrefix))
	store.Delete(types.OutputCoinKey(
		index,
	))
}

// GetAllOutputCoin returns all outputCoin
func (k Keeper) GetAllOutputCoin(ctx sdk.Context) (list []types.OutputCoin) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.OutputCoin
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
