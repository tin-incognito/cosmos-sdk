package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// SetCommitment set a specific commitment in the store from its index
func (k Keeper) SetCommitment(ctx sdk.Context, commitment types.Commitment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentKeyPrefix))
	b := k.cdc.MustMarshal(&commitment)
	store.Set(types.CommitmentKey(
		commitment.Index,
	), b)
}

// GetCommitment returns a commitment from its index
func (k Keeper) GetCommitment(
	ctx sdk.Context,
	index string,

) (val types.Commitment, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentKeyPrefix))

	b := store.Get(types.CommitmentKey(
		index,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveCommitment removes a commitment from the store
func (k Keeper) RemoveCommitment(
	ctx sdk.Context,
	index string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentKeyPrefix))
	store.Delete(types.CommitmentKey(
		index,
	))
}

// GetAllCommitment returns all commitment
func (k Keeper) GetAllCommitment(ctx sdk.Context) (list []types.Commitment) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.CommitmentKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Commitment
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
