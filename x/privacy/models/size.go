package models

import (
	"encoding/json"
	"math"

	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

// GetActualMsgSize returns the size of this msg.
// It is the length of its JSON form.
func GetActualMsgSize(msg *types.MsgPrivacyData) uint64 {
	jsb, err := json.Marshal(msg)
	if err != nil {
		return 0
	}
	return uint64(math.Ceil(float64(len(jsb)) / 1024))
}
