package privacy

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	privacyModels "github.com/cosmos/cosmos-sdk/x/privacy/models"
	privacyTypes "github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func isMintTx(tx sdk.Tx) (bool, error) {
	// no need to check length because privacy decorator has been checked before
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *privacyTypes.MsgPrivacyData:
		if msg.TxType != privacyModels.TxMintType {
			return false, nil
		}
		return true, nil
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", privacyTypes.ModuleName, msg)
		return false, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)

	}
}
