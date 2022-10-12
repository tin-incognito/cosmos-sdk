package models

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type Metadata interface {
	Hash() common.Hash
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	ValidateByItself() (bool, error)
	ValidateByDb() (bool, error)
	ValidateSanity() (bool, error)
}

type OutputCoinReader interface {
	GetOTACoin(ctx sdk.Context, index string) (types.OTACoin, bool)
	GetOutputCoin(ctx sdk.Context, index string) (types.OutputCoin, bool)
}
