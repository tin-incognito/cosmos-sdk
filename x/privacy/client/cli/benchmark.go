package cli

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdBenchmark() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "benchmark [file_path]",
		Short: "Benchmark with list available raw transactions from json file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			filePath := args[0]
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return err
			}
			var rawTxs []string
			err = json.Unmarshal(data, &rawTxs)
			if err != nil {
				return err
			}
			// broadcast raw txs
			err = tx.BroadcastRawPrivacyTx(clientCtx, rawTxs)
			if err != nil {
				return err
			}
			return os.Remove(filePath)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
