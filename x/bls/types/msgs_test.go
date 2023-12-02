package types_test

import (
	"testing"

	"cosmossdk.io/math"
	appparams "github.com/babylonchain/babylon/app/params"
	"github.com/babylonchain/babylon/privval"
	"github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bls/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	pk1      = ed25519.GenPrivKey().PubKey()
	valAddr1 = sdk.ValAddress(pk1.Address())
)

func TestMsgDecode(t *testing.T) {
	registry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(registry)
	types.RegisterInterfaces(registry)
	stakingtypes.RegisterInterfaces(registry)
	cdc := codec.NewProtoCodec(registry)

	// build MsgWrappedCreateValidator
	msg, err := buildMsgWrappedCreateValidatorWithAmount(
		sdk.AccAddress(valAddr1),
		sdk.TokensFromConsensusPower(10, sdk.DefaultPowerReduction),
	)
	require.NoError(t, err)

	// marshal
	msgBytes, err := cdc.MarshalInterface(msg)
	require.NoError(t, err)

	// unmarshal to sdk.Msg interface
	var msg2 sdk.Msg
	err = cdc.UnmarshalInterface(msgBytes, &msg2)
	require.NoError(t, err)

	// type assertion
	msgWithType, ok := msg2.(*types.MsgWrappedCreateValidator)
	require.True(t, ok)

	// ensure msgWithType.MsgCreateValidator.Pubkey with type Any is unmarshaled successfully
	require.NotNil(t, msgWithType.MsgCreateValidator.Pubkey.GetCachedValue())
}

func buildMsgWrappedCreateValidatorWithAmount(addr sdk.AccAddress, bondTokens math.Int) (*types.MsgWrappedCreateValidator, error) {
	tmValPrivkey := ed25519.GenPrivKey()
	bondCoin := sdk.NewCoin(appparams.DefaultBondDenom, bondTokens)
	description := stakingtypes.NewDescription("foo_moniker", "", "", "", "")
	commission := stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())

	pk, err := cryptocodec.FromTmPubKeyInterface(tmValPrivkey.PubKey())
	if err != nil {
		return nil, err
	}

	createValidatorMsg, err := stakingtypes.NewMsgCreateValidator(
		sdk.ValAddress(addr), pk, bondCoin, description, commission, sdk.OneInt(),
	)
	if err != nil {
		return nil, err
	}
	blsPrivKey := bls12381.GenPrivKey()
	pop, err := privval.BuildPoP(tmValPrivkey, blsPrivKey)
	if err != nil {
		return nil, err
	}
	blsPubKey := blsPrivKey.PubKey()

	return types.NewMsgWrappedCreateValidator(createValidatorMsg, &blsPubKey, pop)
}
