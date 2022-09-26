package privacy

import sdk "github.com/cosmos/cosmos-sdk/types"

type TransferValidateSanityDataDecorator struct{}

func NewTransferValidateSanityDataDecorator() TransferValidateSanityDataDecorator {
	return TransferValidateSanityDataDecorator{}
}

func (tvsdd TransferValidateSanityDataDecorator) IsPrivacy() bool {
	return true
}

func (tvsdd TransferValidateSanityDataDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if !isMintTx {

	}
	return next(ctx, tx, simulate)
}

type TransferValidateByItselfDecorator struct{}

func NewTransferValidateByItselfDecorator() TransferValidateByItselfDecorator {
	return TransferValidateByItselfDecorator{}
}

func (tvbid TransferValidateByItselfDecorator) IsPrivacy() bool {
	return true
}

func (tvbid TransferValidateByItselfDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if !isMintTx {

	}
	return next(ctx, tx, simulate)
}

type TransferValidateByDbDecorator struct{}

func NewTransferValidateByDbDecorator() MintValidateByDbDecorator {
	return MintValidateByDbDecorator{}
}

func (tvbdd TransferValidateByDbDecorator) IsPrivacy() bool {
	return true
}

func (tvbdd TransferValidateByDbDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isMintTx, err := isMintTx(tx)
	if err != nil {
		return ctx, err
	}
	if !isMintTx {

	}
	return next(ctx, tx, simulate)
}
