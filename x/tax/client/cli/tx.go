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
	"github.com/cosmos/cosmos-sdk/x/tax/types"
)

// GetTxCmd returns a root CLI command handler for all x/tax transaction commands.
func GetTxCmd() *cobra.Command {
	taxTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Tax transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	taxTxCmd.AddCommand(
		NewCreateFixedAmountTaxCmd(),
		NewCreateRatioTaxCmd(),
		NewStakeCmd(),
		NewUnstakeCmd(),
		NewClaimCmd(),
	)

	return taxTxCmd
}

func NewCreateFixedAmountTaxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-fixed-tax",
		Aliases: []string{"cf"},
		Args:    cobra.ExactArgs(2),
		Short:   "create fixed amount tax tax",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create fixed amount tax tax.
Example:
$ %s tx %s create-fixed-tax --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			taxCreator := clientCtx.GetFromAddress()

			fmt.Println("taxCreator: ", taxCreator)

			// TODO: replace dummy data
			taxPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochDays := uint32(1)
			epochAmount := sdk.Coins{}

			msg := types.NewMsgCreateFixedAmountTax(
				taxPoolAddr,
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

func NewCreateRatioTaxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create-ratio-tax",
		Aliases: []string{"cr"},
		Args:    cobra.ExactArgs(2),
		Short:   "create ratio tax tax",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Create ratio tax tax.
Example:
$ %s tx %s create-ratio-tax --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			taxCreator := clientCtx.GetFromAddress()

			fmt.Println("taxCreator: ", taxCreator)

			// TODO: replace dummy data
			taxPoolAddr := sdk.AccAddress{}
			stakingCoinWeights := sdk.DecCoins{}
			startTime := time.Time{}
			endTime := time.Time{}
			epochDays := uint32(1)
			epochRatio := sdk.Dec{}

			msg := types.NewMsgCreateRatioTax(
				taxPoolAddr,
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
		Short: "stake coins to the tax tax",
		Long: strings.TrimSpace(
			fmt.Sprintf(`stake coins to the tax tax.
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
			taxCreator := clientCtx.GetFromAddress()

			fmt.Println("taxCreator: ", taxCreator)

			// TODO: replace dummy data
			taxID := uint64(1)
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgStake(
				taxID,
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
		Short: "unstake coins from the tax tax",
		Long: strings.TrimSpace(
			fmt.Sprintf(`unstake coins from the tax tax.
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
			taxCreator := clientCtx.GetFromAddress()

			fmt.Println("taxCreator: ", taxCreator)

			// TODO: replace dummy data
			taxID := uint64(1)
			farmer := sdk.AccAddress{}
			stakingCoins := sdk.Coins{}

			msg := types.NewMsgUnstake(
				taxID,
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
		Short: "claim tax rewards from the tax tax",
		Long: strings.TrimSpace(
			fmt.Sprintf(`claim tax rewards from the tax tax.
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
			taxCreator := clientCtx.GetFromAddress()

			fmt.Println("taxCreator: ", taxCreator)

			// TODO: replace dummy data
			taxID := uint64(1)
			farmer := sdk.AccAddress{}

			msg := types.NewMsgClaim(
				taxID,
				farmer,
			)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}
