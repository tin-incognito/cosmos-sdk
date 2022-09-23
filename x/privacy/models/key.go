package models

import (
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
)

func GeneratePrivateKey() key.PrivateKey {
	b := common.RandBytes(common.HashSize)
	return key.GeneratePrivateKey(b)

}
