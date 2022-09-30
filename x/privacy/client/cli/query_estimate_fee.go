package cli

import (
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdEstimateFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "estimate-fee [privacy-data]",
		Short: "Query estimateFee",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqPrivacyData := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryEstimateFeeRequest{

				PrivacyData: reqPrivacyData,
			}

			res, err := queryClient.EstimateFee(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
