package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/incognitochain/go-incognito-sdk-v2/wallet"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdAirdrop() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "airdrop [private_key] [amount]",
		Short: "Broadcast message airdrop",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			privateKey := args[0]
			argAmount := args[1]
			amount, err := strconv.ParseUint(argAmount, 10, 64)
			if err != nil {
				return err
			}
			keyWallet, err := wallet.Base58CheckDeserialize(privateKey)
			if err != nil {
				panic(err)
			}
			keySet := key.KeySet{}
			err = keySet.InitFromPrivateKeyByte(keyWallet.KeySet.PrivateKey)
			if err != nil {
				panic(err)
			}
			otaReceiver := coin.OTAReceiver{}
			err = otaReceiver.FromAddress(keySet.PaymentAddress)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg, err := models.BuildMintTx(otaReceiver, amount, nil, nil)
			if err != nil {
				return err
			}
			clientCtx.GenerateOnly = true

			_, err = tx.GenerateOrBroadcastPrivacyTxCLI(clientCtx, cmd.Flags(), false, msg)
			return err
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
