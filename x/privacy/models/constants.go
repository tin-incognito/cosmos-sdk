package models

const (
	RingSize    = 2
	MaxSizeByte = (1 << 8) - 1
)

// type of transaction mint or transfer
const (
	TxMintType = iota + 1
	TxTransferType
	TxUnshieldType
)

const (
	DefaultFeePerKb = 1
)
