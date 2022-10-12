package models

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func GetMetadataFromMsgPrivacyData(msg *types.MsgPrivacyData) (Metadata, error) {
	var md Metadata
	switch msg.TxType {
	case TxUnshieldType:
		md = &types.MsgUnShield{}
		err := md.Unmarshal(msg.Metadata)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("cannot recognize metadata type")
	}
	return md, nil
}
