package ante

import sdk "github.com/cosmos/cosmos-sdk/types"

type ValidateByDbDecorator struct {
	pk PrivacyKeeper
}

func NewValidateByDbDecorator(privacyKeeper PrivacyKeeper) ValidateByDbDecorator {
	return ValidateByDbDecorator{}
}

func (vbdd ValidateByDbDecorator) IsPrivacy() bool {
	return true
}

func (vbdd ValidateByDbDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isPrivate, err := tx.IsPrivacy()
	if err != nil {
		return ctx, err
	}
	if !isPrivate {
		return next(ctx, tx, simulate)
	}

	return next(ctx, tx, simulate)
}
