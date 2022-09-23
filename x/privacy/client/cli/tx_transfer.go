package cli

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/privacy/common"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdTransfer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [private_key] [payment_infos] [fee]",
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
			feeArgs := args[2]
			fee, err := strconv.ParseUint(feeArgs, 10, 64)
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

			msg, err := models.BuildTransferTx(keySet, paymentInfos, fee, hash)
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastPrivacyTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
