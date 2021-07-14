package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/farming/types"
	"github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
		Use:     "create-fixed-plan",
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
			farmingPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochDays := uint32(1)
			epochAmount := sdk.Coins{}

			msg := types.NewMsgCreateFixedAmountPlan(
				farmingPoolAddr,
				stakingCoinWeights,
				startTime,
				endTime,
				epochDays,
				epochAmount,
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
			farmingPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochDays := uint32(1)
			epochRatio := sdk.Dec{}

			msg := types.NewMsgCreateRatioPlan(
				farmingPoolAddr,
				stakingCoinWeights,
				startTime,
				endTime,
				epochDays,
				epochRatio,
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
			planID := uint64(1)
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgStake(
				planID,
				farmer,
				stakingCoins,
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
			planID := uint64(1)
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgUnstake(
				planID,
				farmer,
				stakingCoins,
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
			planID := uint64(1)
			farmer := sdk.AccAddress{}

			msg := types.NewMsgClaim(
				planID,
				farmer,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

// GetCmdSubmitProposal implements a command handler for submitting a public farming plan creation transaction.
func GetCmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "public-farming-plan [proposal-file] [flags]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a public farming plan creation",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a a public farming plan creation along with an initial deposit.
The proposal details must be supplied via a JSON file.

Example:
$ %s tx gov submit-proposal public-farming-plan <path/to/proposal.json> --from=<key_or_address> --deposit=<deposit_amount>

Where proposal.json contains:

{
    "title": "Public Farming Plan",
    "description": "Here goes first public farming plan!",
    "plans": [
        {
            "@type": "/cosmos.farming.v1beta1.FixedAmountPlan",
            "base_plan": {
                "id": "1",
                "type": "PLAN_TYPE_PUBLIC",
                "farming_pool_address": "",
                "reward_pool_address": "",
                "termination_address": "",
                "staking_reserve_address": "",
                "staking_coin_weights": [
	                {
	                    "denom": "poolCoinDenom",
	                    "amount": "1.000000000000000000"
	                }
                ],
                "start_time": "2021-07-15T08:41:21.662422Z",
                "end_time": "2022-07-16T08:41:21.662422Z",
                "epoch_days": 1
            },
            "epoch_amount": [
                {
	                "denom": "uatom",
	                "amount": "1"
                }
            ]
        }
    ]
}
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			depositStr, err := cmd.Flags().GetString(cli.FlagDeposit)
			if err != nil {
				return err
			}

			deposit, err := sdk.ParseCoinsNormalized(depositStr)
			if err != nil {
				return err
			}

			proposal, err := ParsePublicPlanProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			plans, err := types.UnpackPlans(proposal.Plans)
			if err != nil {
				return err
			}

			for _, plan := range plans {
				if plan.GetType() != types.PlanTypePublic {
					return types.ErrPlanTypeNotAllowed
				}
				// TODO: farming module account.
				// plan.SetFarmingPoolAddress(modAcc)
				// plan.SetTerminationAddress(modAcc)
			}

			content, err := types.NewPublicPlanProposal(proposal.Title, proposal.Description, plans)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(cli.FlagDeposit, "", "deposit of proposal")

	return cmd
}
