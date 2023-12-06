package cli

import (
	"errors"
	"fmt"
	"path/filepath"

	"cosmossdk.io/math"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	flag "github.com/spf13/pflag"

	errorsmod "cosmossdk.io/errors"
	tmconfig "github.com/cometbft/cometbft/config"
	tmos "github.com/cometbft/cometbft/libs/os"

	"github.com/cosmos/cosmos-sdk/privval"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/bls/types"
	cosmoscli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	staketypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// copied from https://github.com/cosmos/cosmos-sdk/blob/7167371f87ae641012549922a292050562821dce/x/staking/client/cli/tx.go#L340
func newBuildCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *staketypes.MsgCreateValidator, error) {
	fAmount, _ := fs.GetString(cosmoscli.FlagAmount)
	amount, err := sdk.ParseCoinNormalized(fAmount)
	if err != nil {
		return txf, nil, err
	}

	valAddr := clientCtx.GetFromAddress()
	pkStr, err := fs.GetString(cosmoscli.FlagPubKey)
	if err != nil {
		return txf, nil, err
	}

	var pk cryptotypes.PubKey
	if err := clientCtx.Codec.UnmarshalInterfaceJSON([]byte(pkStr), &pk); err != nil {
		return txf, nil, err
	}

	moniker, _ := fs.GetString(cosmoscli.FlagMoniker)
	identity, _ := fs.GetString(cosmoscli.FlagIdentity)
	website, _ := fs.GetString(cosmoscli.FlagWebsite)
	security, _ := fs.GetString(cosmoscli.FlagSecurityContact)
	details, _ := fs.GetString(cosmoscli.FlagDetails)
	description := staketypes.NewDescription(
		moniker,
		identity,
		website,
		security,
		details,
	)

	// get the initial validator commission parameters
	rateStr, _ := fs.GetString(cosmoscli.FlagCommissionRate)
	maxRateStr, _ := fs.GetString(cosmoscli.FlagCommissionMaxRate)
	maxChangeRateStr, _ := fs.GetString(cosmoscli.FlagCommissionMaxChangeRate)

	commissionRates, err := buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr)
	if err != nil {
		return txf, nil, err
	}

	// get the initial validator min self delegation
	msbStr, _ := fs.GetString(cosmoscli.FlagMinSelfDelegation)

	minSelfDelegation, ok := math.NewIntFromString(msbStr)
	if !ok {
		return txf, nil, errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "minimum self delegation must be a positive integer")
	}

	msg, err := staketypes.NewMsgCreateValidator(
		sdk.ValAddress(valAddr).String(), pk, amount, description, commissionRates, minSelfDelegation,
	)
	if err != nil {
		return txf, nil, err
	}

	// TODO: temporary skip
	//if err := msg.ValidateBasic(); err != nil {
	//	return txf, nil, err
	//}

	genOnly, _ := fs.GetBool(flags.FlagGenerateOnly)
	if genOnly {
		ip, _ := fs.GetString(cosmoscli.FlagIP)
		nodeID, _ := fs.GetString(cosmoscli.FlagNodeID)

		if nodeID != "" && ip != "" {
			txf = txf.WithMemo(fmt.Sprintf("%s@%s", nodeID, ip))
		}
	}

	return txf, msg, nil
}

// buildWrappedCreateValidatorMsg builds a MsgWrappedCreateValidator that wraps MsgCreateValidator with BLS key
func buildWrappedCreateValidatorMsg(clientCtx client.Context, txf tx.Factory, fs *flag.FlagSet) (tx.Factory, *types.MsgWrappedCreateValidator, error) {
	txf, msg, err := newBuildCreateValidatorMsg(clientCtx, txf, fs)
	if err != nil {
		return txf, nil, err
	}

	home, _ := fs.GetString(flags.FlagHome)
	valKey, err := getValKeyFromFile(home)
	if err != nil {
		return txf, nil, err
	}
	wrappedMsg, err := types.NewMsgWrappedCreateValidator(msg, &valKey.BlsPubkey, valKey.PoP)
	if err != nil {
		return txf, nil, err
	}
	if err := wrappedMsg.ValidateBasic(); err != nil {
		return txf, nil, err
	}

	return txf, wrappedMsg, nil
}

func buildCommissionRates(rateStr, maxRateStr, maxChangeRateStr string) (commission staketypes.CommissionRates, err error) {
	if rateStr == "" || maxRateStr == "" || maxChangeRateStr == "" {
		return commission, errors.New("must specify all validator commission parameters")
	}

	rate, err := math.LegacyNewDecFromStr(rateStr)
	if err != nil {
		return commission, err
	}

	maxRate, err := math.LegacyNewDecFromStr(maxRateStr)
	if err != nil {
		return commission, err
	}

	maxChangeRate, err := math.LegacyNewDecFromStr(maxChangeRateStr)
	if err != nil {
		return commission, err
	}

	commission = staketypes.NewCommissionRates(rate, maxRate, maxChangeRate)

	return commission, nil
}

func getValKeyFromFile(homeDir string) (*privval.ValidatorKeys, error) {
	nodeCfg := tmconfig.DefaultConfig()
	keyPath := filepath.Join(homeDir, nodeCfg.PrivValidatorKeyFile())
	statePath := filepath.Join(homeDir, nodeCfg.PrivValidatorStateFile())
	if !tmos.FileExists(keyPath) {
		return nil, errors.New("validator key file does not exist")
	}
	wrappedPV := privval.LoadWrappedFilePV(keyPath, statePath)

	return privval.NewValidatorKeys(wrappedPV.GetValPrivKey(), wrappedPV.GetBlsPrivKey())
}
