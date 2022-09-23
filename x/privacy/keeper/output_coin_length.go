package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetOutputCoinLength set outputCoinLength in the store
func (k Keeper) SetOutputCoinLength(ctx sdk.Context, outputCoinLength types.OutputCoinLength) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinLengthKey))
	b := k.cdc.MustMarshal(&outputCoinLength)
	store.Set([]byte{0}, b)
}

// GetOutputCoinLength returns outputCoinLength
func (k Keeper) GetOutputCoinLength(ctx sdk.Context) (val types.OutputCoinLength, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinLengthKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveOutputCoinLength removes outputCoinLength from the store
func (k Keeper) RemoveOutputCoinLength(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.OutputCoinLengthKey))
	store.Delete([]byte{0})
}
