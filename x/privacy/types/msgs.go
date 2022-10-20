package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgPrivacyData = "privacy_data"
const TypeMsgShieldData = "shield_data"

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

/*
Shield message
*/

func (m *MsgShield) ValidateBasic() error {
	return nil
}

func (m *MsgShield) GetSigners() []sdk.AccAddress {
	fromAcc, _ := sdk.AccAddressFromBech32(m.GetFrom())
	return []sdk.AccAddress{fromAcc}
}

func (m *MsgShield) IsPrivacy() bool {
	return false
}

func (msg *MsgShield) Route() string {
	return RouterKey
}

func (msg *MsgShield) Type() string {
	return TypeMsgShieldData
}

func (msg *MsgShield) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
