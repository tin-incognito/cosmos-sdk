package tx

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/spf13/pflag"
)

// GenerateOrBroadcastPrivacyTxCLI will either generate and print and unsigned transaction
// or sign it and broadcast it returning an error upon failure.
func GenerateOrBroadcastPrivacyTxCLI(clientCtx client.Context, flagSet *pflag.FlagSet, isFromCosmosTx bool, msgs ...sdk.Msg) (string, error) {
	txf := NewFactoryForPrivacyTxCLI(clientCtx, flagSet)
	return GenerateOrBroadcastPrivacyTxWithFactory(clientCtx, txf, isFromCosmosTx, msgs...)
}

// NewFactoryCLI creates a new Factory.
func NewFactoryForPrivacyTxCLI(clientCtx client.Context, flagSet *pflag.FlagSet) Factory {
	signModeStr := clientCtx.SignModeStr

	signMode := signing.SignMode_SIGN_MODE_UNSPECIFIED
	switch signModeStr {
	case flags.SignModeDirect:
		signMode = signing.SignMode_SIGN_MODE_DIRECT
	case flags.SignModeLegacyAminoJSON:
		signMode = signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON
	case flags.SignModeEIP191:
		signMode = signing.SignMode_SIGN_MODE_EIP_191
	}

	accNum := uint64(0)
	accSeq := uint64(0)
	memo := ""
	gasAdj, _ := flagSet.GetFloat64(flags.FlagGasAdjustment)

	timeoutHeight, _ := flagSet.GetUint64(flags.FlagTimeoutHeight)

	gasStr, _ := flagSet.GetString(flags.FlagGas)
	gasSetting, _ := flags.ParseGasSetting(gasStr)

	f := Factory{
		txConfig:           clientCtx.TxConfig,
		accountRetriever:   clientCtx.AccountRetriever,
		keybase:            clientCtx.Keyring,
		chainID:            clientCtx.ChainID,
		gas:                gasSetting.Gas,
		simulateAndExecute: gasSetting.Simulate,
		accountNumber:      accNum,
		sequence:           accSeq,
		timeoutHeight:      timeoutHeight,
		gasAdjustment:      gasAdj,
		memo:               memo,
		signMode:           signMode,
	}

	feesStr, _ := flagSet.GetString(flags.FlagFees)
	f = f.WithFees(feesStr)

	gasPricesStr, _ := flagSet.GetString(flags.FlagGasPrices)
	f = f.WithGasPrices(gasPricesStr)
	return f
}

// GenerateOrBroadcastTxWithFactory will either generate and print and unsigned transaction
// or sign it and broadcast it returning an error upon failure.
func GenerateOrBroadcastPrivacyTxWithFactory(clientCtx client.Context, txf Factory, isFromCosmosTx bool, msgs ...sdk.Msg) (string, error) {
	// Validate all msgs before generating or broadcasting the tx.
	// We were calling ValidateBasic separately in each CLI handler before.
	// Right now, we're factorizing that call inside this function.
	// ref: https://github.com/cosmos/cosmos-sdk/pull/9236#discussion_r623803504
	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return "", err
		}
	}

	if clientCtx.GenerateOnly {
		var err error
		if isFromCosmosTx {
			txf, err = prepareFactory(clientCtx, txf)
			if err != nil {
				return "", err
			}
		}
		tx, err := GenerateTx(clientCtx, txf, msgs...)
		if err != nil {
			return "", err
		}
		if isFromCosmosTx {
			tx.SetFeeGranter(clientCtx.GetFeeGranterAddress())
			err = Sign(txf, clientCtx.GetFromName(), tx, true)
			if err != nil {
				return "", err
			}
		}
		txBytes, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return "", err
		}
		clientCtx.PrintString(fmt.Sprintf("%s\n", string(txBytes)))
		return string(txBytes), nil
	}

	return "", BroadcastPrivacyTx(clientCtx, txf, msgs...)
}

func BroadcastRawPrivacyTx(clientCtx client.Context, rawTxs []string) error {
	for _, v := range rawTxs {
		tx, err := clientCtx.TxConfig.TxJSONDecoder()([]byte(v))
		if err != nil {
			return err
		}
		txBytes, err := clientCtx.TxConfig.TxEncoder()(tx)
		if err != nil {
			return err
		}
		// broadcast to a Tendermint node
		res, err := clientCtx.BroadcastTx(txBytes)
		if err != nil {
			return err
		}
		err = clientCtx.PrintProto(res)
		if err != nil {
			return err
		}
	}
	return nil
}

// BroadcastPrivacyTx attempts to generate, sign and broadcast a transaction with the
// given set of messages. It will also simulate gas requirements if necessary.
// It will return an error upon failure.
func BroadcastPrivacyTx(clientCtx client.Context, txf Factory, msgs ...sdk.Msg) error {

	if txf.SimulateAndExecute() || clientCtx.Simulate {
		_, adjusted, err := CalculateGas(clientCtx, txf, msgs...)
		if err != nil {
			return err
		}
		txf = txf.WithGas(adjusted)
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", GasEstimateResponse{GasEstimate: txf.Gas()})
	}

	if clientCtx.Simulate {
		return nil
	}

	tx, err := BuildUnsignedTx(txf, msgs...)
	if err != nil {
		return err
	}

	if !clientCtx.SkipConfirm {
		out, err := clientCtx.TxConfig.TxJSONEncoder()(tx.GetTx())
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", out)

		buf := bufio.NewReader(os.Stdin)
		ok, err := input.GetConfirmation("confirm transaction before signing and broadcasting", buf, os.Stderr)

		if err != nil || !ok {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", "cancelled transaction")
			return err
		}
	}

	txBytes, err := clientCtx.TxConfig.TxEncoder()(tx.GetTx())
	if err != nil {
		return err
	}

	// broadcast to a Tendermint node
	res, err := clientCtx.BroadcastTx(txBytes)
	if err != nil {
		return err
	}

	return clientCtx.PrintProto(res)
}
