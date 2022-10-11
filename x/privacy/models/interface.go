package models

import "github.com/cosmos/cosmos-sdk/x/privacy/common"

type Metadata interface {
	Hash() common.Hash
	Marshal() ([]byte, error)
	ValidateByItself() (bool, error)
	ValidateByDb() (bool, error)
	ValidateSanity() (bool, error)
}
