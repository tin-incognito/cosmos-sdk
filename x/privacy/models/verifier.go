package models

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/mlsag"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
)

func VerifySig(ctx sdk.Context, sig, sigPubKey []byte, proof *repos.PaymentProof, fee uint64, msgHash common.Hash, outputCoinReader OutputCoinReader) (bool, error) {
	// Reform Ring
	sumOutputsWithFee := CalculateSumOutputsWithFee(proof.OutputCoinsReadOnly(), fee)
	ring, err := getRingFromSigPubKeyAndLastColumnCommitmentV2(ctx, sumOutputsWithFee, sigPubKey, outputCoinReader)
	if err != nil {
		return false, err
	}

	// Reform MLSAG Signature
	inputCoins := proof.InputCoins()
	keyImages := make([]*operation.Point, len(inputCoins)+1)
	for i := 0; i < len(inputCoins); i++ {
		if inputCoins[i].GetKeyImage() == nil {
			return false, err
		}
		keyImages[i] = inputCoins[i].GetKeyImage()
	}
	// The last column is gone, so just fill in any value
	keyImages[len(inputCoins)] = operation.RandomPoint()
	mlsagSignature, err := getMLSAGSigFromTxSigAndKeyImages(sig, keyImages)
	if err != nil {
		return false, err
	}
	return mlsag.Verify(mlsagSignature, ring, msgHash[:])
}

func VerifySigCA() (bool, error) {
	// TODO: @tin add confidential asset verify later
	return true, nil
}

// Retrieve ring from database using sigpubkey and last column commitment (last column = sumOutputCoinCommitment + fee)
func getRingFromSigPubKeyAndLastColumnCommitmentV2(ctx sdk.Context, sumOutputsWithFee *operation.Point, sigPubKey []byte, outputCoinReader OutputCoinReader) (*mlsag.Ring, error) {
	txSigPubKey := new(SigPubKey)
	if err := txSigPubKey.SetBytes(sigPubKey); err != nil {
		return nil, fmt.Errorf("error when parsing bytes of txSigPubKey %v", err)
	}
	indexes := txSigPubKey.Indexes
	OTAData, err := loadOTAData(ctx, sigPubKey, outputCoinReader)
	if err != nil {
		return nil, err
	}
	n := len(indexes)
	if n == 0 {
		return nil, fmt.Errorf("cannot get ring from Indexes: Indexes is empty")
	}
	m := len(indexes[0])
	if m*n != len(OTAData) {
		return nil, fmt.Errorf("cached OTA data not match with indexes")
	}

	ring := make([][]*operation.Point, n)
	for i := 0; i < n; i++ {
		sumCommitment := new(operation.Point).Identity()
		sumCommitment.Sub(sumCommitment, sumOutputsWithFee)
		row := make([]*operation.Point, m+1)
		for j := 0; j < m; j++ {
			randomCoinBytes := OTAData[i*m+j]
			randomCoin := new(coin.Coin)
			if err := randomCoin.SetBytes(randomCoinBytes); err != nil {
				return nil, err
			}
			row[j] = randomCoin.GetPublicKey()
			sumCommitment.Add(sumCommitment, randomCoin.GetCommitment())
		}
		row[m] = new(operation.Point).Set(sumCommitment)
		ring[i] = row
	}
	return mlsag.NewRing(ring), nil
}

func getMLSAGSigFromTxSigAndKeyImages(txSig []byte, keyImages []*operation.Point) (*mlsag.Sig, error) {
	mlsagSig, err := new(mlsag.Sig).FromBytes(txSig)
	if err != nil {
		return nil, err
	}

	return mlsag.NewMlsagSig(mlsagSig.GetC(), keyImages, mlsagSig.GetR())
}

func loadOTAData(ctx sdk.Context, sigPubKey []byte, outputCoinReader OutputCoinReader) ([][]byte, error) {
	txSigPubKey := new(SigPubKey)
	if err := txSigPubKey.SetBytes(sigPubKey); err != nil {
		return nil, err
	}
	indexes := txSigPubKey.Indexes
	n := len(indexes)
	if n == 0 {
		return nil, fmt.Errorf("cannot get ring from Indexes: Indexes is empty")
	}
	m := len(indexes[0])

	data := make([][]byte, m*n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			index := indexes[i][j]
			otaCoin, ok := outputCoinReader.GetOTACoin(ctx, index.String())
			if !ok {
				return nil, fmt.Errorf("cannot find ota coin")
			}
			outputCoin, ok := outputCoinReader.GetOutputCoin(ctx, otaCoin.OutputCoinIndex)
			if !ok {
				return nil, fmt.Errorf("cannot find outputcoin")
			}
			randomCoinBytes := outputCoin.Value
			data[i*m+j] = randomCoinBytes
			randomCoin := new(coin.Coin)
			if err := randomCoin.SetBytes(randomCoinBytes); err != nil {
				return nil, err
			}
		}
	}
	return data, nil
}
