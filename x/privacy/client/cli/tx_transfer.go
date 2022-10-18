package cli

import (
	"fmt"
	"strconv"
	"strings"

	types2 "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [private_key] [payment_infos] [gasPrice]",
		Short: "Broadcast message transfer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			privateKey := args[0]
			paymentInfos := []*types.MsgTransfer_PaymentInfo{}
			infos := strings.Split(args[1], ",")
			for _, v := range infos {
				temp := strings.Split(v, "-")
				paymentInfo := &types.MsgTransfer_PaymentInfo{}
				for i, value := range temp {
					if i == 0 {
						paymentInfo.PaymentAddress = value
					} else if i == 1 {
						amount, err := strconv.ParseUint(value, 10, 64)
						if err != nil {
							return err
						}
						paymentInfo.Amount = amount
					} else if i == 2 {
						paymentInfo.Info = []byte(value)
					} else {
						return fmt.Errorf("Invalid format payment infos %s", v)
					}
				}
				paymentInfos = append(paymentInfos, paymentInfo)
			}
			gasPriceArgs := args[2]
			gasPriceCoin, err := types2.ParseDecCoin(gasPriceArgs)
			if err != nil {
				return err
			}
			gasPrice := gasPriceCoin.Amount
			if err != nil {
				return err
			}

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
			msg, err := models.BuildTransferTx(keySet, paymentInfos, 1, gasPrice, clientCtx, cmd, nil)
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
			msg, err = models.BuildTransferTx(keySet, paymentInfos, simGasLimit, gasPrice, clientCtx, cmd, nil)
			if err != nil {
				return err
			}
			pflag.Set("gas", fmt.Sprint(simGasLimit))
			pflag.Set("gas-prices", fmt.Sprintf("%v", gasPriceArgs))
			clientCtx.GenerateOnly = true

			_, err = tx.GenerateOrBroadcastPrivacyTxCLI(clientCtx, cmd.Flags(), msg)
			return err
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
