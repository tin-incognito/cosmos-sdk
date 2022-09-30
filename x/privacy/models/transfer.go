package models

import (
	"context"
	"fmt"
	"math/big"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/bulletproofs"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/mlsag"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/operation"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/spf13/cobra"
)

func BuildTransferTx(
	keySet key.KeySet,
	msgTransferPaymentInfos []*types.MsgTransfer_PaymentInfo,
	gasLimit uint64, gasPrice types2.Dec, hashedMessage common.Hash,
	clientContext client.Context, cmd *cobra.Command,
) (*types.MsgPrivacyData, error) {
	var amount uint64
	var err error
	for _, paymentInfo := range msgTransferPaymentInfos {
		amount, err = common.AddUint64(amount, paymentInfo.Amount)
		if err != nil {
			return nil, err
		}
	}

	/*flags.AddPaginationFlagsToCmd(cmd, cmd.Use)*/
	/*flags.AddQueryFlagsToCmd(cmd)*/

	pageReq, err := client.ReadPageRequest(cmd.Flags())
	if err != nil {

		return nil, err
	}

	queryClient := types.NewQueryClient(clientContext)

	params := &types.QueryAllOutputCoinRequest{
		Pagination: pageReq,
	}

	outputCoins, err := queryClient.OutputCoinAll(context.Background(), params)
	if err != nil {
		return nil, err
	}
	outcoins := outputCoins.OutputCoin

	coins, paymentInfos, fee, err := chooseCoinsByKeySet(outcoins, keySet, amount, msgTransferPaymentInfos, gasLimit, gasPrice, clientContext, cmd)
	if err != nil {
		return nil, err
	}

	return buildTransferTx(coins, keySet, paymentInfos, fee, hashedMessage, clientContext, cmd)
}

func buildTransferTx(
	inputCoins []*OutputCoin, keySet key.KeySet, paymentInfos []*key.PaymentInfo, fee uint64, hashedMessage common.Hash,
	clientContext client.Context, cmd *cobra.Command,
) (*types.MsgPrivacyData, error) {
	proof, outputCoins, err := Prove(inputCoins, paymentInfos)
	if err != nil {
		return nil, err
	}

	res, err := SignOnMessage(inputCoins, outputCoins, &keySet.PrivateKey, hashedMessage.Bytes(), fee, clientContext, cmd)
	if err != nil {
		return nil, err
	}
	res.Proof = proof.Bytes()
	res.TxType = TxTransferType
	res.Fee = fee

	return res, nil
}

func Prove(
	inputCoins []*OutputCoin,
	paymentInfos []*key.PaymentInfo,
) (*repos.PaymentProof, []*coin.Coin, error) {
	outputCoins, err := GenerateOutputCoinsByPaymentInfos(paymentInfos)
	if err != nil {
		return nil, nil, err
	}

	proof, err := prove(inputCoins, outputCoins, nil, false, paymentInfos)
	if err != nil {
		return nil, nil, err
	}

	return proof, outputCoins, nil
}

func prove(
	inputCoins []*OutputCoin, outputCoins []*coin.Coin,
	sharedSecrets []*operation.Point,
	hasConfidentialAsset bool,
	paymentInfos []*key.PaymentInfo,
) (*repos.PaymentProof, error) {
	proof := repos.NewPaymentProof()
	ic := make([]*coin.Coin, len(inputCoins))
	for i, v := range inputCoins {
		ic[i] = coin.NewCoin()
		*ic[i] = *v.value
	}
	err := proof.SetInputCoins(ic)
	if err != nil {
		return nil, err
	}
	err = proof.SetOutputCoins(outputCoins)
	if err != nil {
		return nil, err
	}

	// Prepare range proofs
	n := len(outputCoins)
	outputValues := make([]uint64, n)
	outputRands := make([]*operation.Scalar, n)
	for i := 0; i < n; i++ {
		outputValues[i] = outputCoins[i].GetValue()
		outputRands[i] = outputCoins[i].GetRandomness()
	}

	wit := new(bulletproofs.AggregatedRangeWitness)
	wit.Set(outputValues, outputRands)
	if hasConfidentialAsset {
		/*blinders := make([]*operation.Scalar, len(sharedSecrets))*/
		/*for i := range sharedSecrets {*/
		/*if sharedSecrets[i] == nil {*/
		/*blinders[i] = new(operation.Scalar).FromUint64(0)*/
		/*} else {*/
		/*blinders[i], err = coin.ComputeAssetTagBlinder(sharedSecrets[i])*/
		/*if err != nil {*/
		/*return nil, err*/
		/*}*/
		/*}*/
		/*}*/
		/*var err error*/
		/*wit, err = bulletproofs.TransformWitnessToCAWitness(wit, blinders)*/
		/*if err != nil {*/
		/*return nil, err*/
		/*}*/

		/*theBase, err := bulletproofs.GetFirstAssetTag(outputCoins)*/
		/*if err != nil {*/
		/*return nil, err*/
		/*}*/
		/*proof.aggregatedRangeProof, err = wit.ProveUsingBase(theBase)*/

		/*outputCommitments := make([]*operation.Point, n)*/
		/*for i := 0; i < n; i++ {*/
		/*com, err := outputCoins[i].ComputeCommitmentCA()*/
		/*if err != nil {*/
		/*return nil, err*/
		/*}*/
		/*outputCommitments[i] = com*/
		/*}*/
		/*proof.aggregatedRangeProof.SetCommitments(outputCommitments)*/
		/*if err != nil {*/
		/*return nil, err*/
		/*}*/
	} else {
		aggregatedRangeProof, err := wit.Prove()
		if err != nil {
			return nil, err
		}
		proof.SetAggregatedRangeProof(aggregatedRangeProof)
	}

	// After Prove, we should hide all information in coin details.
	for i, outputCoin := range proof.OutputCoins() {
		//if !common.IsPublicKeyBurningAddress(outputCoin.GetPublicKey().ToBytesS()) {
		if err = outputCoin.ConcealOutputCoin(paymentInfos[i].PaymentAddress.GetPublicView()); err != nil {
			return nil, err
		}

		// OutputCoin.GetKeyImage should be nil even though we do not have it
		// Because otherwise the RPC server will return the Bytes of [1 0 0 0 0 ...] (the default byte)
		outputCoins[i].SetKeyImage(nil)
		proof.SetOutputCoinAtIndex(i, &outputCoin)
		//}
	}

	for i, inputCoin := range proof.InputCoins() {
		inputCoin.ConcealInputCoin()
		proof.SetInputCoinAtIndex(i, &inputCoin)
	}

	return proof, nil
}

func SignOnMessage(
	inputCoins []*OutputCoin, outputCoins []*coin.Coin,
	privateKey *key.PrivateKey, hashedMessage []byte, fee uint64,
	clientContext client.Context, cmd *cobra.Command,
) (*types.MsgPrivacyData, error) {
	tx := new(types.MsgPrivacyData)
	// Generate Ring
	piBig, err := common.RandBigIntMaxRange(big.NewInt(int64(RingSize)))
	if err != nil {
		return nil, err
	}
	var pi int = int(piBig.Int64())
	ring, indexes, commitmentToZero, err := generateMlsagRingWithIndexes(inputCoins, outputCoins, pi, fee, clientContext, cmd)
	if err != nil {
		return nil, err
	}

	// Set SigPubKey
	txSigPubKey := new(SigPubKey)
	txSigPubKey.Indexes = indexes
	tx.SigPubKey, err = txSigPubKey.Bytes()
	if err != nil {
		return nil, err
	}

	// Set sigPrivKey
	privKeysMlsag, err := createPrivKeyMlsag(inputCoins, outputCoins, privateKey, commitmentToZero)
	if err != nil {
		return nil, err
	}
	sag := mlsag.NewMlsag(privKeysMlsag, ring, pi)

	// Set Signature
	mlsagSignature, err := sag.Sign(hashedMessage)
	if err != nil {
		return nil, err
	}

	// inputCoins already hold keyImage so set to nil to reduce size
	mlsagSignature.SetKeyImages(nil)
	tx.Sig, err = mlsagSignature.ToBytes()

	return tx, nil
}

func generateMlsagRingWithIndexes(
	inputCoins []*OutputCoin, outputCoins []*coin.Coin, pi int, fee uint64,
	clientContext client.Context, cmd *cobra.Command,
) (*mlsag.Ring, [][]*big.Int, *operation.Point, error) {
	outputCoinsAsGeneric := make([]*coin.Coin, len(outputCoins))
	for i := 0; i < len(outputCoins); i++ {
		outputCoinsAsGeneric[i] = outputCoins[i]
	}
	sumOutputsWithFee := CalculateSumOutputsWithFee(outputCoinsAsGeneric, fee)
	indexes := make([][]*big.Int, RingSize)
	ring := make([][]*operation.Point, RingSize)
	var commitmentToZero *operation.Point
	attempts := 0

	queryClient := types.NewQueryClient(clientContext)

	params := &types.QueryGetOutputCoinLengthRequest{}

	outputCoinLength, err := queryClient.OutputCoinLength(context.Background(), params)
	if err != nil {
		return nil, nil, nil, err
	}

	otaCoinsLength := big.NewInt(0).SetBytes(outputCoinLength.OutputCoinLength.Value)
	pageReq, err := client.ReadPageRequest(cmd.Flags())
	if err != nil {
		return nil, nil, nil, err
	}

	queryClient = types.NewQueryClient(clientContext)
	p0 := &types.QueryAllOutputCoinRequest{
		Pagination: pageReq,
	}

	tempOutputCoins, err := queryClient.OutputCoinAll(context.Background(), p0)
	if err != nil {
		return nil, nil, nil, err
	}
	allOutputCoins := tempOutputCoins.OutputCoin

	mOutputCoins := make(map[string]types.OutputCoin)
	for _, v := range allOutputCoins {
		mOutputCoins[v.Index] = v
	}

	pageReq, err = client.ReadPageRequest(cmd.Flags())
	if err != nil {
		return nil, nil, nil, err
	}

	queryClient = types.NewQueryClient(clientContext)
	p1 := &types.QueryAllOTACoinRequest{
		Pagination: pageReq,
	}

	tempOTACoins, err := queryClient.OTACoinAll(context.Background(), p1)
	if err != nil {
		return nil, nil, nil, err
	}
	otaCoins := tempOTACoins.OTACoin
	mOtaCoins := make(map[string]types.OTACoin)
	for _, v := range otaCoins {
		mOtaCoins[v.Index] = v
	}

	for i := 0; i < RingSize; i++ {
		sumInputs := new(operation.Point).Identity()
		sumInputs.Sub(sumInputs, sumOutputsWithFee)

		row := make([]*operation.Point, len(inputCoins))
		rowIndexes := make([]*big.Int, len(inputCoins))
		if i == pi {
			for j := 0; j < len(inputCoins); j++ {
				row[j] = inputCoins[j].value.GetPublicKey()
				c, found := mOutputCoins[inputCoins[j].index]
				if !found {
					return nil, nil, nil, fmt.Errorf("Cannot find inputCoin index")
				}
				rowIndexes[j] = big.NewInt(0).SetBytes(c.SerialNumber)
				sumInputs.Add(sumInputs, inputCoins[j].value.GetCommitment())
			}
		} else {
			for j := 0; j < len(inputCoins); j++ {
				var coinDB *coin.Coin
				for attempts < coin.MaxAttempts { // The chance of infinite loop is negligible
					coinDB = new(coin.Coin)
					rowIndexes[j], _ = common.RandBigIntMaxRange(otaCoinsLength)
					otaCoin, found := mOtaCoins[rowIndexes[j].String()]
					if !found {
						return nil, nil, nil, fmt.Errorf("Cannot find otaCoin")
					}

					c, found := mOutputCoins[otaCoin.OutputCoinIndex]
					if !found {
						return nil, nil, nil, fmt.Errorf("Cannot find outputCoin")
					}
					if err := coinDB.SetBytes(c.Value); err != nil {
						return nil, nil, nil, err
					}

					/*// we do not use burned coins since they will reduce the privacy level of the transaction.*/
					/*if !common.IsPublicKeyBurningAddress(coinDB.GetPublicKey().ToBytesS()) {*/
					/*break*/
					/*}*/
					attempts++
				}
				if attempts == coin.MaxAttempts && coinDB == nil {
					return nil, nil, nil, fmt.Errorf("cannot form decoys")
				}

				row[j] = coinDB.GetPublicKey()
				sumInputs.Add(sumInputs, coinDB.GetCommitment())
			}
		}
		row = append(row, sumInputs)
		if i == pi {
			commitmentToZero = sumInputs
		}
		ring[i] = row
		indexes[i] = rowIndexes
	}
	return mlsag.NewRing(ring), indexes, commitmentToZero, nil
}

func createPrivKeyMlsag(
	inputCoins []*OutputCoin, outputCoins []*coin.Coin,
	privateKey *key.PrivateKey, commitmentToZero *operation.Point,
) ([]*operation.Scalar, error) {
	sumRand := new(operation.Scalar).FromUint64(0)
	for _, in := range inputCoins {
		sumRand.Add(sumRand, in.value.GetRandomness())
	}
	for _, out := range outputCoins {
		sumRand.Sub(sumRand, out.GetRandomness())
	}

	privKeyMlsag := make([]*operation.Scalar, len(inputCoins)+1)
	for i := 0; i < len(inputCoins); i++ {
		var err error
		privKeyMlsag[i], err = inputCoins[i].value.ParsePrivateKeyOfCoin(*privateKey)
		if err != nil {
			return nil, err
		}
	}
	commitmentToZeroRecomputed := new(operation.Point).ScalarMult(operation.PedCom.G[operation.PedersenRandomnessIndex], sumRand)
	match := operation.IsPointEqual(commitmentToZeroRecomputed, commitmentToZero)
	if !match {
		return nil, fmt.Errorf("error : asset tag sum or commitment sum mismatch")
	}
	privKeyMlsag[len(inputCoins)] = sumRand
	return privKeyMlsag, nil
}
