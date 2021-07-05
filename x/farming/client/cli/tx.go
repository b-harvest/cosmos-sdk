package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/farming/types"

	"github.com/spf13/cobra"
)

// GetTxCmd returns a root CLI command handler for all x/farming transaction commands.
func GetTxCmd() *cobra.Command {
	farmingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Farming transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	farmingTxCmd.AddCommand(
		NewCreateFixedAmountPlanCmd(),
		NewCreateRatioPlanCmd(),
		NewStakeCmd(),
		NewUnstakeCmd(),
		NewClaimCmd(),
	)

	return farmingTxCmd
}

func NewCreateFixedAmountPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-fixed-plan ",
		Aliases: []string{"cf"},
		Args:    cobra.ExactArgs(2),
		Short:   "create fixed amount farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create fixed amount farming plan.
Example:
$ %s tx %s create-fixed-plan --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farming_pool_address := sdk.AccAddress{}
			staking_coins_weight := sdk.DecCoins{}
			start_time := time.Time{}
			end_time := time.Time{}
			epoch_days := uint32(1)
			epoch_amount := sdk.Coins{}

			msg := types.NewMsgCreateFixedAmountPlan(
				farming_pool_address,
				staking_coins_weight,
				&start_time,
				&end_time,
				epoch_days,
				epoch_amount,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewCreateRatioPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-ratio-plan",
		Aliases: []string{"cr"},
		Args:    cobra.ExactArgs(2),
		Short:   "create ratio farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create ratio farming plan.
Example:
$ %s tx %s create-ratio-plan --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			farming_pool_address := sdk.AccAddress{}
			staking_coins_weight := sdk.DecCoins{}
			start_time := time.Time{}
			end_time := time.Time{}
			epoch_days := uint32(1)
			epoch_ratio := sdk.Dec{}

			msg := types.NewMsgCreateRatioPlan(
				farming_pool_address,
				staking_coins_weight,
				&start_time,
				&end_time,
				epoch_days,
				epoch_ratio,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake",
		Args:  cobra.ExactArgs(2),
		Short: "stake coins to the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`stake coins to the farming plan.
Example:
$ %s tx %s stake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			plan_id := uint64(1)
			farmer := sdk.AccAddress{}
			staking_coins := sdk.Coins{}

			msg := types.NewMsgStake(
				plan_id,
				farmer,
				staking_coins,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unstake",
		Args:  cobra.ExactArgs(2),
		Short: "unstake coins from the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`unstake coins from the farming plan.
Example:
$ %s tx %s unstake --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			plan_id := uint64(1)
			farmer := sdk.AccAddress{}
			staking_coins := sdk.Coins{}

			msg := types.NewMsgUnstake(
				plan_id,
				farmer,
				staking_coins,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func NewClaimCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "claim",
		Args:  cobra.ExactArgs(2),
		Short: "claim farming rewards from the farming plan",
		Long: strings.TrimSpace(
			fmt.Sprintf(`claim farming rewards from the farming plan.
Example:
$ %s tx %s claim --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			planCreator := clientCtx.GetFromAddress()

			fmt.Println("planCreator: ", planCreator)

			// TODO: replace dummy data
			plan_id := uint64(1)
			farmer := sdk.AccAddress{}

			msg := types.NewMsgClaim(
				plan_id,
				farmer,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
