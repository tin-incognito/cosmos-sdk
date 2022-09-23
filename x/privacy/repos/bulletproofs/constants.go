package bulletproofs

import "github.com/cosmos/cosmos-sdk/x/privacy/common"

const (
	precompPedGValIndex = iota
	precompPedGRandIndex
	precompUIndex
	precompGIndex
)

const (
	aggParamNMax = common.MaxOutputCoin * common.MaxExp
)

func precompHIndex(paramN int) int {
	return precompGIndex + paramN
}

const (
	MaxExp = 64
)
