package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func IsPrivacyTx(tx sdk.Tx) (bool, error) {
	msgs := tx.GetMsgs()
	if len(msgs) != 1 {
		return false, nil
	}
	msg := tx.GetMsgs()[0]
	switch msg := msg.(type) {
	case *types.MsgPrivacyData:
		return true, nil
	default:
		errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
		return false, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
	}
}
