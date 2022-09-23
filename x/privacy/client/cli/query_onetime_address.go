package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	"github.com/spf13/cobra"
)

func CmdListOnetimeAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list-onetime-address",
		Short: "list all OnetimeAddress",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryAllOnetimeAddressRequest{
				Pagination: pageReq,
			}

			res, err := queryClient.OnetimeAddressAll(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddPaginationFlagsToCmd(cmd, cmd.Use)
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func CmdShowOnetimeAddress() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-onetime-address [index]",
		Short: "shows a OnetimeAddress",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			argIndex := args[0]

			params := &types.QueryGetOnetimeAddressRequest{
				Index: argIndex,
			}

			res, err := queryClient.OnetimeAddress(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
