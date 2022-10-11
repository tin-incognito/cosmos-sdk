package types

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
)

func (m *MsgUnShield) Hash() common.Hash {
	data, _ := json.Marshal(m)
	return common.HashH(data)
}
