package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgPrivacyData = "privacy_data"

var _ sdk.Msg = &MsgPrivacyData{}

func NewMsgPrivacyData(
	lockTime, fee uint64,
	info, sigPubKey, sig, proof []byte,
	txType int32, metadata []byte,
) *MsgPrivacyData {
	return &MsgPrivacyData{
		LockTime:  lockTime,
		Fee:       fee,
		Info:      info,
		SigPubKey: sigPubKey,
		Sig:       sig,
		TxType:    txType,
		Metadata:  metadata,
	}
}

func (msg *MsgPrivacyData) Route() string {
	return RouterKey
}

func (msg *MsgPrivacyData) Type() string {
	return TypeMsgPrivacyData
}

func (msg *MsgPrivacyData) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}

func (msg *MsgPrivacyData) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgPrivacyData) ValidateBasic() error {
	return nil
}

func (msg *MsgPrivacyData) IsPrivacy() bool {
	return true
}
