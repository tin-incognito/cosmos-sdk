package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type PrivacyKeeper interface {
	GetSerialNumber(ctx sdk.Context, index string) (types.SerialNumber, bool)
	GetOnetimeAddress(ctx sdk.Context, index string) (types.OnetimeAddress, bool)
	GetOTACoin(ctx sdk.Context, index string) (types.OTACoin, bool)
	GetOutputCoin(ctx sdk.Context, index string) (types.OutputCoin, bool)
}
