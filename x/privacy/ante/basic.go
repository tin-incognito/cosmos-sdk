package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

type ValidateByItself struct {
	c *Cache
}

func NewValidateByItself(c *Cache) ValidateByItself {
	return ValidateByItself{c: c}
}

func (vbi ValidateByItself) IsPrivacy() bool {
	return true
}

func (vbi ValidateByItself) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	isPrivate, err := tx.IsPrivacy()
	if err != nil {
		return ctx, err
	}
	if !isPrivate {
		return next(ctx, tx, simulate)
	}

	// no need to check index, has been checked before
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		isMintTx, err := IsMintTx(tx)
		if err != nil {
			return ctx, err
		}
		if !isMintTx {
			//TODO: @tin validate metadata by itself
			key, err := common.NewHashFromBytes(msg.Hash)
			if err != nil {
				return ctx, err
			}
			proof, err := vbi.c.GetProof(*key)
			if err != nil {
				proof := repos.NewPaymentProof()
				if err = proof.SetBytes(msg.Proof); err != nil {
					return ctx, err
				}
				if err = vbi.c.AddProof(*key, proof); err != nil {
					return ctx, err
				}
			}

			// TODO: @tin add confidential asset verify later
			isValid, err := proof.Verify()
			if err != nil {
				return ctx, err
			}
			if !isValid {
				return ctx, fmt.Errorf("Verify proof fail")
			}
		}

	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return ctx, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}

	return next(ctx, tx, simulate)
}
