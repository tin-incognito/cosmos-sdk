package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/privacy/types"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group privacy queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdQueryParams())
	cmd.AddCommand(CmdListSerialNumber())
	cmd.AddCommand(CmdShowSerialNumber())
	cmd.AddCommand(CmdListOutputCoin())
	cmd.AddCommand(CmdShowOutputCoin())
	cmd.AddCommand(CmdListCommitment())
	cmd.AddCommand(CmdShowCommitment())
	cmd.AddCommand(CmdListToken())
	cmd.AddCommand(CmdShowToken())
	cmd.AddCommand(CmdListOnetimeAddress())
	cmd.AddCommand(CmdShowOnetimeAddress())
	cmd.AddCommand(CmdBalance())

	cmd.AddCommand(CmdListOTACoin())
	cmd.AddCommand(CmdShowOTACoin())
	cmd.AddCommand(CmdEstimateFee())
	// this line is used by starport scaffolding # 1

	return cmd
}
