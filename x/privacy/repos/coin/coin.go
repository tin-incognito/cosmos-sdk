package coin

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/common/base58"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"

	proto "github.com/gogo/protobuf/proto"
)

//nolint:revive // skip linter for this struct name
// Coin is the struct that will be stored to db
// If not privacy, mask and amount will be the original randomness and value
// If has privacy, mask and amount will be as paper monero
type Coin struct {
	// Public
	Info       []byte           `protobuf:"bytes,1,opt,name=info,proto3,omitempty" json:"info,omitempty"`
	PublicKey  *operation.Point `protobuf:"bytes,2,opt,name=public_key,proto3,omitempty" json:"public_key,omitempty"`
	Commitment *operation.Point `protobuf:"bytes,3,opt,name=commitment,proto3,omitempty" json:"commitment,omitempty"`
	KeyImage   *operation.Point `protobuf:"bytes,4,opt,name=key_image,proto3,omitempty" json:"key_image,omitempty"`

	// sharedRandom and txRandom is shared secret between receiver and giver
	// sharedRandom is only visible when creating coins, when it is broadcast to network, it will be set to null
	SharedConcealRandom *operation.Scalar `protobuf:"bytes,5,opt,name=shared_conceal_random,proto3,omitempty" json:"shared_conceal_random,omitempty"` //rConceal: shared random when concealing output coin and blinding assetTag
	SharedRandom        *operation.Scalar `protobuf:"bytes,6,opt,name=shared_random,proto3,omitempty" json:"shared_random,omitempty"`                 // rOTA: shared random when creating one-time-address
	TxRandom            *TxRandom         `protobuf:"bytes,7,opt,name=tx_random,proto3,omitempty" json:"tx_random,omitempty"`                         // rConceal*G + rOTA*G + index

	// mask = randomness
	// amount = value
	Mask   *operation.Scalar `protobuf:"bytes,8,opt,name=mask,proto3,omitempty" json:"mask,omitempty"`
	Amount *operation.Scalar `protobuf:"bytes,9,opt,name=amount,proto3,omitempty" json:"creator,omitempty"`
	// tag is nil unless confidential asset
	AssetTag *operation.Point `protobuf:"bytes,10,opt,name=asset_tag,proto3,omitempty" json:"asset_tag,omitempty"`
}

// ParsePrivateKeyOfCoin retrieves the private OTA key of coin from the Master PrivateKey
func (c Coin) ParsePrivateKeyOfCoin(privKey key.PrivateKey) (*operation.Scalar, error) {
	keySet := new(key.KeySet)
	if err := keySet.InitFromPrivateKey(&privKey); err != nil {
		return nil, err
	}
	_, txRandomOTAPoint, index, err := c.GetTxRandomDetail()
	if err != nil {
		return nil, err
	}
	rK := new(operation.Point).ScalarMult(txRandomOTAPoint, keySet.OTAKey.GetOTASecretKey()) // (r_ota*G) * k = r_ota * K
	H := operation.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))       // Hash(r_ota*K, index)

	k := new(operation.Scalar).FromBytesS(privKey)
	return new(operation.Scalar).Add(H, k), nil // Hash(rK, index) + privSpend
}

// ParseKeyImageWithPrivateKey retrieves the keyImage of coin from the Master PrivateKey
func (c Coin) ParseKeyImageWithPrivateKey(privKey key.PrivateKey) (*operation.Point, error) {
	k, err := c.ParsePrivateKeyOfCoin(privKey)
	if err != nil {
		return nil, err
	}
	Hp := operation.HashToPoint(c.GetPublicKey().ToBytesS())
	return new(operation.Point).ScalarMult(Hp, k), nil
}

// Conceal the amount of coin using the publicView of the receiver
//
//	- AdditionalData: must be the publicView of the receiver
func (c *Coin) ConcealOutputCoin(additionalData *operation.Point) error {
	// If this coin is already encrypted or it is created by other person then cannot conceal
	if c.IsEncrypted() || c.GetSharedConcealRandom() == nil {
		return nil
	}
	publicView := additionalData
	rK := new(operation.Point).ScalarMult(publicView, c.GetSharedConcealRandom()) // rK = sharedConcealRandom * publicView

	hash := operation.HashToScalar(rK.ToBytesS()) // hash(rK)
	hash = operation.HashToScalar(hash.ToBytesS())
	mask := new(operation.Scalar).Add(c.GetRandomness(), hash) // mask = c.mask + hash

	hash = operation.HashToScalar(hash.ToBytesS())
	amount := new(operation.Scalar).Add(c.GetAmount(), hash) // amount = c.amout + hash
	c.SetRandomness(mask)
	c.SetAmount(amount)
	c.SetSharedConcealRandom(nil)
	c.SetSharedRandom(nil)
	return nil
}

// Conceal the input coin of a transaction: keep only the keyImage
func (c *Coin) ConcealInputCoin() {
	c.SetValue(0)
	c.SetRandomness(nil)
	c.SetPublicKey(nil)
	c.SetCommitment(nil)
	c.SetTxRandomDetail(new(operation.Point).Identity(), new(operation.Point).Identity(), 0)
}

// Decrypt a coin using the corresponding KeySet
func (c *Coin) Decrypt(keySet *key.KeySet) (*Coin, error) {
	if keySet == nil {
		return nil, fmt.Errorf("cannot Decrypt Coin: Keyset is empty")
	}

	// Must parse keyImage first in any situation
	if len(keySet.PrivateKey) > 0 {
		keyImage, err := c.ParseKeyImageWithPrivateKey(keySet.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("cannot parse key image with privateKey Coin - %v", err)
		}
		c.SetKeyImage(keyImage)
	}

	if !c.IsEncrypted() {
		return c, nil
	}

	viewKey := keySet.ReadonlyKey
	if len(viewKey.Rk) == 0 && len(keySet.PrivateKey) == 0 {
		return nil, fmt.Errorf("cannot Decrypt Coin: Keyset does not contain viewkey or privatekey")
	}

	if viewKey.GetPrivateView() != nil {
		txConcealRandomPoint, err := c.GetTxRandom().GetTxConcealRandomPoint()
		if err != nil {
			return nil, err
		}
		rK := new(operation.Point).ScalarMult(txConcealRandomPoint, viewKey.GetPrivateView())

		// Hash multiple times
		hash := operation.HashToScalar(rK.ToBytesS())
		hash = operation.HashToScalar(hash.ToBytesS())
		randomness := c.GetRandomness().Sub(c.GetRandomness(), hash)

		// Hash 1 more time to get value
		hash = operation.HashToScalar(hash.ToBytesS())
		value := c.GetAmount().Sub(c.GetAmount(), hash)

		commitment := operation.PedCom.CommitAtIndex(value, randomness, operation.PedersenValueIndex)
		// for `confidential asset` coin, we commit differently
		if c.GetAssetTag() != nil {
			com, err := ComputeCommitmentCA(c.GetAssetTag(), randomness, value)
			if err != nil {
				return nil, fmt.Errorf("cannot recompute commitment when decrypting")
			}
			commitment = com
		}
		if !operation.IsPointEqual(commitment, c.GetCommitment()) {
			return nil, fmt.Errorf("cannot Decrypt Coin: Commitment is not the same after decrypt")
		}
		c.SetRandomness(randomness)
		c.SetAmount(value)
	}
	return c, nil
}

func NewCoin() *Coin {
	c := new(Coin)
	c.Info = []byte{}
	c.PublicKey = new(operation.Point).Identity()
	c.Commitment = new(operation.Point).Identity()
	c.KeyImage = new(operation.Point).Identity()
	c.SharedRandom = new(operation.Scalar)
	c.TxRandom = NewTxRandom()
	c.Mask = new(operation.Scalar)
	c.Amount = new(operation.Scalar)
	return c
}

func (c Coin) IsEncrypted() bool {
	if c.Mask == nil || c.Amount == nil {
		return true
	}
	tempCommitment := operation.PedCom.CommitAtIndex(c.Amount, c.Mask, operation.PedersenValueIndex)
	if c.GetAssetTag() != nil {
		// err is only for nil parameters, which we already checked, so here it is ignored
		com, _ := c.ComputeCommitmentCA()
		tempCommitment = com
	}
	return !operation.IsPointEqual(tempCommitment, c.Commitment)
}

func (c Coin) GetRandomness() *operation.Scalar          { return c.Mask }
func (c Coin) GetAmount() *operation.Scalar              { return c.Amount }
func (c Coin) GetSharedRandom() *operation.Scalar        { return c.SharedRandom }
func (c Coin) GetSharedConcealRandom() *operation.Scalar { return c.SharedConcealRandom }
func (c Coin) GetPublicKey() *operation.Point            { return c.PublicKey }
func (c Coin) GetCommitment() *operation.Point           { return c.Commitment }
func (c Coin) GetKeyImage() *operation.Point             { return c.KeyImage }
func (c Coin) GetInfo() []byte                           { return c.Info }
func (c Coin) GetAssetTag() *operation.Point             { return c.AssetTag }
func (c Coin) GetValue() uint64 {
	if c.IsEncrypted() {
		return 0
	}
	return c.Amount.ToUint64Little()
}
func (c Coin) GetTxRandom() *TxRandom { return c.TxRandom }
func (c Coin) GetTxRandomDetail() (*operation.Point, *operation.Point, uint32, error) {
	txRandomOTAPoint, err1 := c.TxRandom.GetTxOTARandomPoint()
	txRandomConcealPoint, err2 := c.TxRandom.GetTxConcealRandomPoint()
	index, err3 := c.TxRandom.GetIndex()
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, nil, 0, fmt.Errorf("cannot Get TxRandom: point or index is wrong")
	}
	return txRandomConcealPoint, txRandomOTAPoint, index, nil
}

func (c Coin) GetCoinDetailEncrypted() []byte {
	return c.GetAmount().ToBytesS()
}

func (c *Coin) GetCoinID() [operation.Ed25519KeySize]byte {
	if c.PublicKey != nil {
		return c.PublicKey.ToBytes()
	}
	return [operation.Ed25519KeySize]byte{}
}

func (c *Coin) SetRandomness(mask *operation.Scalar)           { c.Mask = mask }
func (c *Coin) SetAmount(amount *operation.Scalar)             { c.Amount = amount }
func (c *Coin) SetSharedRandom(sharedRandom *operation.Scalar) { c.SharedRandom = sharedRandom }
func (c *Coin) SetSharedConcealRandom(sharedConcealRandom *operation.Scalar) {
	c.SharedConcealRandom = sharedConcealRandom
}
func (c *Coin) SetTxRandom(txRandom *TxRandom) {
	if txRandom == nil {
		c.TxRandom = nil
	} else {
		c.TxRandom = NewTxRandom()
		_ = c.TxRandom.SetBytes(txRandom.Bytes())
	}
}
func (c *Coin) SetTxRandomDetail(txConcealRandomPoint, txRandomPoint *operation.Point, index uint32) {
	res := new(TxRandom)
	res.SetTxConcealRandomPoint(txConcealRandomPoint)
	res.SetTxOTARandomPoint(txRandomPoint)
	res.SetIndex(index)
	c.TxRandom = res
}

func (c *Coin) SetPublicKey(publicKey *operation.Point)   { c.PublicKey = publicKey }
func (c *Coin) SetCommitment(commitment *operation.Point) { c.Commitment = commitment }
func (c *Coin) SetKeyImage(keyImage *operation.Point)     { c.KeyImage = keyImage }
func (c *Coin) SetValue(value uint64)                     { c.Amount = new(operation.Scalar).FromUint64(value) }
func (c *Coin) SetInfo(b []byte) {
	c.Info = make([]byte, len(b))
	copy(c.Info, b)
}
func (c *Coin) SetAssetTag(at *operation.Point) { c.AssetTag = at }

func (c Coin) Bytes() []byte {
	coinBytes := []byte{}
	info := c.GetInfo()
	byteLengthInfo := byte(getMin(len(info), MaxSizeInfoCoin))
	coinBytes = append(coinBytes, byteLengthInfo)
	coinBytes = append(coinBytes, info[:byteLengthInfo]...)

	if c.PublicKey != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.PublicKey.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.Commitment != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.Commitment.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.KeyImage != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.KeyImage.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.SharedRandom != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.SharedRandom.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.SharedConcealRandom != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.SharedConcealRandom.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.TxRandom != nil {
		coinBytes = append(coinBytes, TxRandomGroupSize)
		coinBytes = append(coinBytes, c.TxRandom.Bytes()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.Mask != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.Mask.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.Amount != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.Amount.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	if c.AssetTag != nil {
		coinBytes = append(coinBytes, byte(operation.Ed25519KeySize))
		coinBytes = append(coinBytes, c.AssetTag.ToBytesS()...)
	} else {
		coinBytes = append(coinBytes, byte(0))
	}

	return coinBytes
}

func (c *Coin) SetBytes(coinBytes []byte) error {
	var err error
	if c == nil {
		return fmt.Errorf("cannot set bytes for unallocated Coin")
	}
	if len(coinBytes) == 0 {
		return fmt.Errorf("coinBytes is empty")
	}

	offset := 0
	c.Info, err = parseInfoForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin info error: %v", err)
	}

	c.PublicKey, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin publicKey error: %v", err)
	}
	c.Commitment, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin commitment error: %v", err)
	}
	c.KeyImage, err = parsePointForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin keyImage error: %v", err)
	}
	c.SharedRandom, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin mask error: %v", err)
	}

	c.SharedConcealRandom, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin mask error: %v", err)
	}

	if offset >= len(coinBytes) {
		return fmt.Errorf("offset is larger than len(bytes), cannot parse txRandom")
	}
	if coinBytes[offset] != TxRandomGroupSize {
		return fmt.Errorf("setBytes Coin TxRandomGroup error: length of TxRandomGroup is not correct")
	}
	offset++
	if offset+TxRandomGroupSize > len(coinBytes) {
		return fmt.Errorf("setBytes Coin TxRandomGroup error: length of coinBytes is too small")
	}
	c.TxRandom = NewTxRandom()
	err = c.TxRandom.SetBytes(coinBytes[offset : offset+TxRandomGroupSize])
	if err != nil {
		return fmt.Errorf("setBytes Coin TxRandomGroup error: %v", err)
	}
	offset += TxRandomGroupSize

	if err != nil {
		return fmt.Errorf("setBytes Coin txRandom error: %v", err)
	}
	c.Mask, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin mask error: %v", err)
	}
	c.Amount, err = parseScalarForSetBytes(&coinBytes, &offset)
	if err != nil {
		return fmt.Errorf("setBytes Coin amount error: %v", err)
	}

	if offset >= len(coinBytes) {
		// for parsing old serialization, which does not have assetTag field
		c.AssetTag = nil
	} else {
		c.AssetTag, err = parsePointForSetBytes(&coinBytes, &offset)
		if err != nil {
			return fmt.Errorf("setBytes Coin assetTag error: %v", err)
		}
	}
	return nil
}

// HashH returns the SHA3-256 hashing of coin bytes array
func (c *Coin) HashH() *common.Hash {
	hash := common.HashH(c.Bytes())
	return &hash
}

func (c Coin) MarshalJSON() ([]byte, error) {
	data := c.Bytes()
	temp := base58.Base58Check{}.Encode(data, common.ZeroByte)
	return json.Marshal(temp)
}

func (c *Coin) UnmarshalJSON(data []byte) error {
	dataStr := ""
	_ = json.Unmarshal(data, &dataStr)
	temp, _, err := base58.Base58Check{}.Decode(dataStr)
	if err != nil {
		return err
	}
	err = c.SetBytes(temp)
	if err != nil {
		return err
	}
	return nil
}

func (c *Coin) CheckCoinValid(paymentAdd key.PaymentAddress, sharedRandom []byte, amount uint64) bool {
	if c.GetValue() != amount {
		return false
	}
	// check one-time address is corresponding to paymentaddress
	r := new(operation.Scalar).FromBytesS(sharedRandom)
	if !r.ScalarValid() {
		return false
	}

	publicOTA := paymentAdd.GetOTAPublicKey()
	if publicOTA == nil {
		return false
	}
	rK := new(operation.Point).ScalarMult(publicOTA, r)
	_, txOTARandomPoint, index, err := c.GetTxRandomDetail()
	if err != nil {
		return false
	}
	if !operation.IsPointEqual(new(operation.Point).ScalarMultBase(r), txOTARandomPoint) {
		return false
	}

	hash := operation.HashToScalar(append(rK.ToBytesS(), common.Uint32ToBytes(index)...))
	HrKG := new(operation.Point).ScalarMultBase(hash)
	tmpPubKey := new(operation.Point).Add(HrKG, paymentAdd.GetPublicSpend())
	return bytes.Equal(tmpPubKey.ToBytesS(), c.PublicKey.ToBytesS())
}

// Check whether the utxo is from this keyset
func (c *Coin) DoesCoinBelongToKeySet(keySet *key.KeySet) (bool, *operation.Point) {
	_, txOTARandomPoint, index, err1 := c.GetTxRandomDetail()
	if err1 != nil {
		return false, nil
	}

	// Check if the utxo belong to this one-time-address
	rK := new(operation.Point).ScalarMult(txOTARandomPoint, keySet.OTAKey.GetOTASecretKey())

	hashed := operation.HashToScalar(
		append(rK.ToBytesS(), common.Uint32ToBytes(index)...),
	)

	HnG := new(operation.Point).ScalarMultBase(hashed)
	KCheck := new(operation.Point).Sub(c.GetPublicKey(), HnG)

	// // Retrieve the sharedConcealRandomPoint for generating the blinded assetTag
	// var rSharedConcealPoint *operation.Point
	// if keySet.ReadonlyKey.GetPrivateView() != nil {
	// 	rSharedConcealPoint = new(operation.Point).ScalarMult(txConcealRandomPoint, keySet.ReadonlyKey.GetPrivateView())
	// }

	return operation.IsPointEqual(KCheck, keySet.OTAKey.GetPublicSpend()), rK
}

func (c *Coin) Reset()         { c = NewCoin() }
func (c *Coin) String() string { return proto.CompactTextString(c) }
func (c *Coin) ProtoMessage()  {}
