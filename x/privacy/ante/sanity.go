package ante

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type ValidateSanityDecorator struct{}

func NewValidateSanityDecorator() ValidateSanityDecorator {
	return ValidateSanityDecorator{}
}

func (vsd ValidateSanityDecorator) IsPrivacy() bool {
	return true
}

func (vsd ValidateSanityDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isPrivate := tx.IsPrivacy()

	if !isPrivate {
		return next(ctx, tx, simulate)
	}

	// no need to check index, has been checked before
	msg := tx.GetMsgs()[0]
	isValid, err := validateSanity(msg)
	if err != nil {
		return ctx, err
	}
	if !isValid {
		return ctx, fmt.Errorf("Fail to validate sanity")
	}

	// @tin TODO: validate sanity of metadata here
	return next(ctx, tx, simulate)
}

func validateSanity(msg sdk.Msg) (bool, error) {
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		// check LockTime before now
		if msg.LockTime > uint64(time.Now().Unix()) {
			return false, fmt.Errorf("wrong tx locktime %d", msg.LockTime)
		}

		// check tx size
		actualTxSize := models.GetActualMsgSize(msg)
		if actualTxSize > common.MaxTxSize {
			return false, fmt.Errorf("tx size %d kB is too large", actualTxSize)
		}
		proof := repos.NewPaymentProof()
		err := proof.SetBytes(msg.Proof)
		if err != nil {
			return false, err
		}
		valid, err := proof.ValidateSanity()
		if err != nil {
			return false, err
		}
		if !valid {
			return false, fmt.Errorf("cannot validate sanity of proof")
		}

		switch msg.TxType {
		case models.TxMintType, models.TxTransferType: // is valid
		default:
			return false, fmt.Errorf("wrong tx type with %v", msg.TxType)
		}

		info := msg.Info
		if len(info) > common.MaxSizeInfo {
			return false, fmt.Errorf("wrong tx info length %d bytes, only support info with max length <= %d bytes", len(info), 512)
		}
		return true, nil
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return false, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
