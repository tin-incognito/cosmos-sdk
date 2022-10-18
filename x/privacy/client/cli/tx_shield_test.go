package cli

import (
	"context"
	"strconv"
	"testing"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/x/privacy/models"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/coin"
	"github.com/cosmos/cosmos-sdk/x/privacy/repos/key"
	"github.com/spf13/cobra"
)

func TestGenerateRawShieldTx(t *testing.T) {
	//encodingConfig := simappparams.MakeTestEncodingConfig()
	cmd := &cobra.Command{
		Use:  "c",
		Run:  func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) (err error) { return nil },
	}
	cmd.ExecuteContext(context.Background())
	cmd.Flags().String(flags.FlagChainID, "my-test-chain", "The network chain ID")
	clientCtx, err := client.GetClientTxContext(cmd)
	if err != nil {
		panic(err)
	}
	/*clientCtx.*/
	/*WithTxConfig(encodingConfig.TxConfig).*/
	/*WithCodec(encodingConfig.Marshaler)*/

	paymentAddr := "12scjGeftnVsH4Xa2CsRpkqctZaY77asQQMq84Gqx3itNRZaaQxfCDARXfrVZrsSK63pGPC2DzYwdEhLAjAyrUKErGeZfcL2v7HXXTVLee6Gwvr5NsruJRCqiHnQ9aaGsYGDKy8mgTzu1pJfPdv8"
	argAmount := "1000000"
	amount, err := strconv.ParseUint(argAmount, 10, 64)
	if err != nil {
		panic(err)
	}

	receiver := key.PaymentAddress{}
	err = receiver.ImportFromString(paymentAddr)
	if err != nil {
		panic(err)
	}

	otaReceiver := coin.OTAReceiver{}
	err = otaReceiver.FromAddress(receiver)
	if err != nil {
		panic(err)
	}

	// Get depositor address
	from := clientCtx.GetFromAddress()

	msg, err := models.BuildShieldTx(from, otaReceiver, amount, nil)
	if err != nil {
		panic(err)
	}

	clientCtx.GenerateOnly = true

	err = tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
	if err != nil {
		panic(err)
	}
}
