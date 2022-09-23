package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetSerialNumber set a specific serialNumber in the store from its index
func (k Keeper) SetSerialNumber(ctx sdk.Context, serialNumber types.SerialNumber) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SerialNumberKeyPrefix))
	b := k.cdc.MustMarshal(&serialNumber)
	store.Set(types.SerialNumberKey(
		serialNumber.Index,
	), b)
}

// GetSerialNumber returns a serialNumber from its index
func (k Keeper) GetSerialNumber(
	ctx sdk.Context,
	index string,

) (val types.SerialNumber, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SerialNumberKeyPrefix))

	b := store.Get(types.SerialNumberKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveSerialNumber removes a serialNumber from the store
func (k Keeper) RemoveSerialNumber(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SerialNumberKeyPrefix))
	store.Delete(types.SerialNumberKey(
		index,
	))
}

// GetAllSerialNumber returns all serialNumber
func (k Keeper) GetAllSerialNumber(ctx sdk.Context) (list []types.SerialNumber) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.SerialNumberKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.SerialNumber
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
