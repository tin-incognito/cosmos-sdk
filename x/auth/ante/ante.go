package ante

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	privacyAnte "github.com/cosmos/cosmos-sdk/x/privacy/ante"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	AccountKeeper   AccountKeeper
	BankKeeper      types.BankKeeper
	FeegrantKeeper  FeegrantKeeper
	PrivacyKeeper   privacyAnte.PrivacyKeeper
	SignModeHandler authsigning.SignModeHandler
	SigGasConsumer  func(meter sdk.GasMeter, sig signing.SignatureV2, params types.Params) error
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrLogic, "sign mode handler is required for ante builder")
	}

	var sigGasConsumer = options.SigGasConsumer
	if sigGasConsumer == nil {
		sigGasConsumer = DefaultSigVerificationGasConsumer
	}

	privacyCache := privacyAnte.NewCache()

	anteDecorators := []sdk.AnteDecorator{
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		NewMempoolFeeDecorator(),
		NewValidateBasicDecorator(),
		NewTxTimeoutHeightDecorator(),
		NewValidateMemoDecorator(options.AccountKeeper),
		NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper),
		NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(options.AccountKeeper),
		NewSigGasConsumeDecorator(options.AccountKeeper, sigGasConsumer),
		NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		NewIncrementSequenceDecorator(options.AccountKeeper),
		// Add validate double spend in mem pool later for mint tx
		privacyAnte.NewValidateSanityDecorator(privacyCache),
		privacyAnte.NewValidateByItself(privacyCache),
		privacyAnte.NewValidateByDbDecorator(options.PrivacyKeeper, privacyCache),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}
