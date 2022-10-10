package models

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/schnorr"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
)

func BuildShieldTx(
	from sdk.AccAddress,
	otaReceiver coin.OTAReceiver,
	amount uint64,
	info []byte,
) (*types.MsgShield, error) {
	outputCoin, err := GenerateOutputCoin(amount, info, otaReceiver)
	if err != nil {
		return nil, err
	}

	proof := repos.NewPaymentProof()
	proof.SetOutputCoins([]*coin.Coin{outputCoin})

	hash := MsgHash(uint64(time.Now().Unix()), 0, proof, nil)

	res := &types.MsgShield{
		Hash:   hash.Bytes(),
		From:   from.String(),
		Amount: amount,
		Proof:  proof.Bytes(),
	}
	return res, nil
}

func BuildMintTx(
	otaReceiver coin.OTAReceiver,
	amount uint64,
	info []byte,
	md []byte,
) (*types.MsgPrivacyData, error) {
	privateKey := GeneratePrivateKey()
	outputCoin, err := GenerateOutputCoin(amount, info, otaReceiver)
	if err != nil {
		return nil, err
	}

	proof := repos.NewPaymentProof()
	proof.SetOutputCoins([]*coin.Coin{outputCoin})
	lockTime := uint64(time.Now().Unix())

	hash := MsgHash(lockTime, 0, proof, nil)

	sig, sigPubKey, err := SignNoPrivacy(&privateKey, hash.Bytes())
	if err != nil {
		return nil, err
	}

	res := &types.MsgPrivacyData{
		Hash:      hash.Bytes(),
		Proof:     proof.Bytes(),
		SigPubKey: sigPubKey,
		Sig:       sig,
		LockTime:  lockTime,
		Info:      info,
		TxType:    TxMintType,
		Metadata:  md,
	}
	return res, nil
}

func SignNoPrivacy(privateKey *key.PrivateKey, hashedMessage []byte) (signatureBytes []byte, sigPubKey []byte, err error) {
	/****** using Schnorr signature *******/
	// sign with sigPrivKey
	// prepare private key for Schnorr
	sk := new(operation.Scalar).FromBytesS(*privateKey)
	r := new(operation.Scalar).FromUint64(0)
	sigKey := new(schnorr.SchnorrPrivateKey)
	sigKey.Set(sk, r)
	signature, err := sigKey.Sign(hashedMessage)
	if err != nil {
		return nil, nil, err
	}

	signatureBytes = signature.Bytes()
	sigPubKey = sigKey.GetPublicKey().GetPublicKey().ToBytesS()
	return signatureBytes, sigPubKey, nil
}

func VerifySigNoPrivacy(sig, sigPubKey, hashedMessage []byte) (bool, error) {
	// check input transaction
	if sig == nil || sigPubKey == nil {
		return false, fmt.Errorf("transaction input must be signed")
	}

	var err error
	/****** verify Schnorr signature *****/
	// prepare Public key for verification
	verifyKey := new(schnorr.SchnorrPublicKey)
	sigPublicKey, err := new(operation.Point).FromBytesS(sigPubKey)

	if err != nil {
		return false, err
	}
	verifyKey.Set(sigPublicKey)

	// convert signature from byte array to SchnorrSign
	signature := new(schnorr.SchnSignature)
	err = signature.SetBytes(sig)
	if err != nil {
		return false, err
	}

	res := verifyKey.Verify(signature, hashedMessage)
	return res, nil
}
