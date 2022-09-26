package privacy

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MintValidateSanityDataDecorator struct{}

func NewMintValidateSanityDataDecorator() MintValidateSanityDataDecorator {
	return MintValidateSanityDataDecorator{}
}

func (mvsdd MintValidateSanityDataDecorator) IsPrivacy() bool {
	return true
}

func (mvsdd MintValidateSanityDataDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if isMintTx {

	}
	return next(ctx, tx, simulate)
}

type MintValidateByItselfDecorator struct{}

func NewMintValidateByItselfDecorator() MintValidateByItselfDecorator {
	return MintValidateByItselfDecorator{}
}

func (mvbid MintValidateByItselfDecorator) IsPrivacy() bool {
	return true
}

func (mvbid MintValidateByItselfDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if isMintTx {

	}
	return next(ctx, tx, simulate)
}

type MintValidateByDbDecorator struct{}

func NewMintValidateByDbDecorator() MintValidateByDbDecorator {
	return MintValidateByDbDecorator{}
}

func (mvbdd MintValidateByDbDecorator) IsPrivacy() bool {
	return true
}

func (mvbdd MintValidateByDbDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if isMintTx {

	}
	return next(ctx, tx, simulate)
}
