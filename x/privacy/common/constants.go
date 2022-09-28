package common

const (
	HashSize          = 32 // bytes
	MaxHashStringSize = HashSize * 2
	Uint32Size        = 4 // bytes
	ZeroByte          = byte(0x00)
	CheckSumLen       = 4  // bytes
	PrivateKeySize    = 32 // bytes
	MaxTxSize         = 500
	MaxSizeInfo       = 512
)

const (
	PriKeyType                  = byte(0x0) // Serialize wallet account key into string with only PRIVATE KEY of account keyset
	PaymentAddressType          = byte(0x1) // Serialize wallet account key into string with only PAYMENT ADDRESS of account keyset
	ReadonlyKeyType             = byte(0x2) // Serialize wallet account key into string with only READONLY KEY of account keyset
	OTAKeyType                  = byte(0x3) // Serialize wallet account key into string with only OTA KEY of account keyset
	PrivateReceivingAddressType = byte(0x4) // prefix for marshalled receiving address (coin pk + txRandom), not used for KeyWallet
)
