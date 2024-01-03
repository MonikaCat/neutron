package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"

	"github.com/MonikaCat/neutron/v2/x/dex/types"
)

func CmdWithdrawFilledLimitOrder() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "withdraw-filled-limit-order [tranche-key]",
		Short:   "Broadcast message WithdrawFilledLimitOrder",
		Example: "withdraw-filled-limit-order TRANCHEKEY123 --from alice",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgWithdrawFilledLimitOrder(
				clientCtx.GetFromAddress().String(),
				args[0],
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
