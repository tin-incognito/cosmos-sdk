package cli

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	types2 "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/spf13/cobra"
	"strconv"
)

var _ = strconv.Itoa(0)

func CmdUnshield() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unshield [private_key] [nonprivacy_address] [amount] [gasPrice]",
		Short: "Broadcast message transfer",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			privateKey := args[0]
			toAddress := args[1]
			amountStr := args[2]
			amount, err := strconv.ParseUint(amountStr, 10, 64)
			if err != nil {
				return err
			}
			paymentInfos := []*types.MsgTransfer_PaymentInfo{{"12RxahVABnAVCGP3LGwCn8jkQxgw7z1x14wztHzn455TTVpi1wBq9YGwkRMQg3J4e657AbAnCvYCJSdA9czBUNuCKwGSRQt55Xwz8WA", amount, nil}}

			gasPriceArgs := args[3]
			gasPriceCoin, err := types2.ParseDecCoin(gasPriceArgs)
			if err != nil {
				return err
			}
			gasPrice := gasPriceCoin.Amount
			if err != nil {
				return err
			}

			m := types.MsgTransfer{
				PrivateKey:   privateKey,
				PaymentInfos: paymentInfos,
			}

			msgBytes, err := json.Marshal(m)
			if err != nil {
				return err
			}
			hash := common.HashH(msgBytes)
			keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
			if err != nil {
				return err
			}
			keySet := key.KeySet{}
			err = keySet.InitFromPrivateKeyByte(keyWallet.KeySet.PrivateKey)
			if err != nil {
				return err
			}

			//simulate
			unshield := types.MsgUnShield{
				ToAdrress: toAddress,
				Amount:    amount,
			}
			msg, err := models.BuildTransferTx(keySet, paymentInfos, 1, gasPrice, hash, clientCtx, cmd, unshield)
			if err != nil {
				return err
			}

			pflag := cmd.Flags()
			pflag.Set("gas-adjustment", fmt.Sprint(1.01))
			txf := tx.NewFactoryForPrivacyTxCLI(clientCtx, pflag)
			_, simGasLimit, err := tx.CalculateGas(clientCtx, txf, msg)

			if err != nil {
				return err
			}
			msg, err = models.BuildTransferTx(keySet, paymentInfos, simGasLimit, gasPrice, hash, clientCtx, cmd, unshield)
			if err != nil {
				return err
			}
			msg.TxType = models.TxUnshieldType
			if err != nil {
				return err
			}
			pflag.Set("gas", fmt.Sprint(simGasLimit))
			pflag.Set("gas-prices", fmt.Sprintf("%v", gasPriceArgs))

			return tx.GenerateOrBroadcastPrivacyTxCLI(clientCtx, pflag, msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
