package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"

	"cosmossdk.io/core/address"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmoscli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"cosmossdk.io/x/bls/types"
)

//nolint:unused
const (
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	listSeparator              = ","
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(ac address.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdTxAddBlsSig())
	cmd.AddCommand(CmdWrappedCreateValidator(ac))

	return cmd
}

func CmdTxAddBlsSig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit [epoch_number] [last_commit_hash] [bls_sig] [signer address]",
		Short: "submit a BLS signature",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			epoch_num, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			lch, err := types.NewLastCommitHashFromHex(args[1])
			if err != nil {
				return err
			}

			blsSig, err := bls12381.NewBLSSigFromHex(args[2])
			if err != nil {
				return err
			}

			addr, err := sdk.ValAddressFromBech32(args[3])
			if err != nil {
				return err
			}

			msg := types.NewMsgAddBlsSig(clientCtx.GetFromAddress(), epoch_num, lch, blsSig, addr)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// LEGACY

func flagSetDescriptionCreate() *flag.FlagSet {
	fs := flag.NewFlagSet("", flag.ContinueOnError)

	fs.String(cosmoscli.FlagMoniker, "", "The validator's name")
	fs.String(cosmoscli.FlagIdentity, "", "The optional identity signature (ex. UPort or Keybase)")
	fs.String(cosmoscli.FlagWebsite, "", "The validator's (optional) website")
	fs.String(cosmoscli.FlagSecurityContact, "", "The validator's (optional) security contact email")
	fs.String(cosmoscli.FlagDetails, "", "The validator's (optional) details")

	return fs
}

// NewCreateValidatorCmd returns a CLI command handler for creating a MsgCreateValidator transaction.
func NewCreateValidatorCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-validator",
		Short: "create new validator initialized with a self-delegation to it",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, cmd.Flags())
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
		},
	}

	cmd.Flags().AddFlagSet(cosmoscli.FlagSetPublicKey())
	cmd.Flags().AddFlagSet(cosmoscli.FlagSetAmount())
	cmd.Flags().AddFlagSet(flagSetDescriptionCreate())
	cmd.Flags().AddFlagSet(cosmoscli.FlagSetCommissionCreate())
	cmd.Flags().AddFlagSet(cosmoscli.FlagSetMinSelfDelegation())

	cmd.Flags().String(cosmoscli.FlagIP, "", fmt.Sprintf("The node's public IP. It takes effect only when used in combination with --%s", flags.FlagGenerateOnly))
	cmd.Flags().String(cosmoscli.FlagNodeID, "", "The node's ID")
	flags.AddTxFlagsToCmd(cmd)

	_ = cmd.MarkFlagRequired(flags.FlagFrom)
	_ = cmd.MarkFlagRequired(cosmoscli.FlagAmount)
	_ = cmd.MarkFlagRequired(cosmoscli.FlagPubKey)
	_ = cmd.MarkFlagRequired(cosmoscli.FlagMoniker)

	return cmd
}

func CmdWrappedCreateValidator(ac address.Codec) *cobra.Command {
	//cmd := cosmoscli.NewCreateValidatorCmd(ac)
	cmd := NewCreateValidatorCmd()
	cmd.Long = strings.TrimSpace(
		string(`create-validator will create a new validator initialized
with a self-delegation to it using the BLS key generated for the validator (e.g., via babylond create-bls-key).

This command creates a MsgWrappedCreateValidator message which is a wrapper of cosmos-sdk's
create validator with a pair of BLS key. The BLS key should exist in priv_validator_key.json
before running the command (e.g., via babylond create-bls-key).`))
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		clientCtx, err := client.GetClientTxContext(cmd)
		if err != nil {
			return err
		}

		txf, err := tx.NewFactoryCLI(clientCtx, cmd.Flags())

		if err != nil {
			return err
		}

		txf = txf.WithTxConfig(clientCtx.TxConfig).WithAccountRetriever(clientCtx.AccountRetriever)

		txf, msg, err := buildWrappedCreateValidatorMsg(clientCtx, txf, cmd.Flags())
		if err != nil {
			return err
		}

		return tx.GenerateOrBroadcastTxWithFactory(clientCtx, txf, msg)
	}
	// HACK: test cases need to setup the path where the priv validator BLS key is going to be set
	// so we redefine the FlagHome here. Since we can't import `app` due to a cyclic dependency,
	// we have to duplicate the definition here.
	// If this changes, the `DefaultHomeDir` flag at `app/app.go` needs to change as well.
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	// TODO: fix
	defaultNodeHome := filepath.Join(userHomeDir, ".babylond")
	cmd.Flags().String(flags.FlagHome, defaultNodeHome, "The node home directory")

	return cmd
}
