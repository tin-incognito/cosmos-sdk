package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/x/privacy/models"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdShield() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shield [payment_address] [amount]",
		Short: "Broadcast message shield",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			paymentAddr := args[0]
			argAmount := args[1]
			amount, err := strconv.ParseUint(argAmount, 10, 64)
			if err != nil {
				return err
			}

			receiver := key.PaymentAddress{}
			err = receiver.ImportFromString(paymentAddr)
			if err != nil {
				panic(err)
			}

			otaReceiver := coin.OTAReceiver{}
			err = otaReceiver.FromAddress(receiver)
			if err != nil {
				return err
			}

			// Get depositor address
			from := clientCtx.GetFromAddress()

			msg, err := models.BuildShieldTx(from, otaReceiver, amount, nil)
			if err != nil {
				return err
			}
			clientCtx.GenerateOnly = true

			_, err = tx.GenerateOrBroadcastPrivacyTxCLI(clientCtx, cmd.Flags(), true, msg)
			//return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
			return err
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
