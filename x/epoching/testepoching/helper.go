//go:build app_v1
// +build app_v1

package testepoching

import (
	"math/rand"
	"testing"

	"cosmossdk.io/x/epoching/testepoching/datagen"
	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"

	"cosmossdk.io/math"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	app "cosmossdk.io/simapp"

	"cosmossdk.io/x/epoching/keeper"
	"cosmossdk.io/x/epoching/types"
)

type ValidatorInfo struct {
	BlsKey  bls12381.PrivateKey
	Address sdk.ValAddress
}

// Helper is a structure which wraps the entire app and exposes functionalities for testing the epoching module
type Helper struct {
	t *testing.T

	Ctx            sdk.Context
	App            *app.SimApp
	EpochingKeeper *keeper.Keeper
	MsgSrvr        types.MsgServer
	QueryClient    types.QueryClient
	StakingKeeper  *stakingkeeper.Keeper

	GenAccs        []authtypes.GenesisAccount
	ValBlsPrivKeys []ValidatorInfo
}

// NewHelper creates the helper for testing the epoching module
func NewHelper(t *testing.T) *Helper {
	app := app.Setup(t, false)
	ctx := app.BaseApp.NewContext(false)

	epochingKeeper := app.EpochingKeeper

	// add BLS pubkey to the genesis validator
	valSet := epochingKeeper.GetValidatorSet(ctx, 0)
	require.Len(t, valSet, 1)
	genesisVal := valSet[0]
	blsPrivKey := bls12381.GenPrivKey()
	genesisBLSPubkey := blsPrivKey.PubKey()
	err := app.BLSKeeper.CreateRegistration(ctx, genesisBLSPubkey, genesisVal.Addr)
	require.NoError(t, err)

	querier := keeper.Querier{Keeper: epochingKeeper}
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)
	msgSrvr := keeper.NewMsgServerImpl(epochingKeeper)

	return &Helper{
		t,
		ctx,
		app,
		&epochingKeeper,
		msgSrvr,
		queryClient,
		app.StakingKeeper,
		nil,
		[]ValidatorInfo{ValidatorInfo{
			blsPrivKey,
			genesisVal.Addr,
		}},
	}
}

// NewHelperWithValSet is same as NewHelper, except that it creates a set of validators
func NewHelperWithValSet(t *testing.T) *Helper {
	// generate the validator set with 10 validators
	tmValSet, err := GenTmValidatorSet(10)
	require.NoError(t, err)

	// generate the genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	// ensure the genesis account has a sufficient amount of tokens
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.DefaultPowerReduction.MulRaw(10000000))),
	}
	GenAccs := []authtypes.GenesisAccount{acc}

	// setup the app and ctx
	app := app.SetupWithGenesisValSet(t, tmValSet, GenAccs, balance)
	ctx := app.BaseApp.NewContext(false)

	// get necessary subsets of the app/keeper
	epochingKeeper := app.EpochingKeeper
	valInfos := []ValidatorInfo{}
	// add BLS pubkey to the genesis validator
	valSet := epochingKeeper.GetValidatorSet(ctx, 0)
	for _, val := range valSet {
		blsPrivKey := bls12381.GenPrivKey()
		valInfos = append(valInfos, ValidatorInfo{blsPrivKey, val.Addr})
		blsPubkey := blsPrivKey.PubKey()
		err = app.BLSKeeper.CreateRegistration(ctx, blsPubkey, val.Addr)
		require.NoError(t, err)
	}
	querier := keeper.Querier{Keeper: epochingKeeper}
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, querier)
	queryClient := types.NewQueryClient(queryHelper)
	msgSrvr := keeper.NewMsgServerImpl(epochingKeeper)

	return &Helper{t, ctx, app, &epochingKeeper, msgSrvr, queryClient, app.StakingKeeper, GenAccs, valInfos}
}

// GenAndApplyEmptyBlock generates a new empty block and appends it to the current blockchain
func (h *Helper) GenAndApplyEmptyBlock(r *rand.Rand) sdk.Context {
	newHeight := h.App.LastBlockHeight() + 1
	valSet, err := h.StakingKeeper.GetLastValidators(h.Ctx)
	if err != nil {
		panic(err)
	}
	valhash := CalculateValHash(valSet)
	newHeader := tmproto.Header{
		Height:             newHeight,
		ValidatorsHash:     valhash,
		NextValidatorsHash: valhash,
		AppHash:            datagen.GenRandomByteArray(r, 32),
		LastCommitHash:     datagen.GenRandomLastCommitHash(r),
	}

	//h.App.BeginBlock(abci.RequestBeginBlock{Header: newHeader})
	//h.App.EndBlock(abci.RequestEndBlock{})
	h.App.Commit()

	h.Ctx = h.Ctx.WithBlockHeader(newHeader)
	return h.Ctx
}

//func (h *Helper) BeginBlock(r *rand.Rand) sdk.Context {
//	newHeight := h.App.LastBlockHeight() + 1
//	valSet := h.StakingKeeper.GetLastValidators(h.Ctx)
//	valhash := CalculateValHash(valSet)
//	newHeader := tmproto.Header{
//		Height:             newHeight,
//		AppHash:            datagen.GenRandomByteArray(r, 32),
//		ValidatorsHash:     valhash,
//		NextValidatorsHash: valhash,
//	}
//
//	h.App.BeginBlock(abci.RequestBeginBlock{Header: newHeader})
//	h.Ctx = h.Ctx.WithBlockHeader(newHeader)
//	return h.Ctx
//}
//
//func (h *Helper) EndBlock() sdk.Context {
//	h.App.EndBlock(abci.RequestEndBlock{})
//	h.App.Commit()
//	return h.Ctx
//}

// WrappedDelegate calls handler to delegate stake for a validator
func (h *Helper) WrappedDelegate(delegator sdk.AccAddress, val sdk.ValAddress, amount math.Int) *sdk.Result {
	coin := sdk.NewCoin(sdk.DefaultBondDenom, amount)
	msg := stakingtypes.NewMsgDelegate(delegator.String(), val.String(), coin)
	wmsg := types.NewMsgWrappedDelegate(msg)
	return h.Handle(func(ctx sdk.Context) (proto.Message, error) {
		return h.MsgSrvr.WrappedDelegate(ctx, wmsg)
	})
}

// WrappedDelegateWithPower calls handler to delegate stake for a validator
func (h *Helper) WrappedDelegateWithPower(delegator sdk.AccAddress, val sdk.ValAddress, power int64) *sdk.Result {
	coin := sdk.NewCoin(sdk.DefaultBondDenom, h.StakingKeeper.TokensFromConsensusPower(h.Ctx, power))
	msg := stakingtypes.NewMsgDelegate(delegator.String(), val.String(), coin)
	wmsg := types.NewMsgWrappedDelegate(msg)
	return h.Handle(func(ctx sdk.Context) (proto.Message, error) {
		return h.MsgSrvr.WrappedDelegate(ctx, wmsg)
	})
}

// WrappedUndelegate calls handler to unbound some stake from a validator.
func (h *Helper) WrappedUndelegate(delegator sdk.AccAddress, val sdk.ValAddress, amount math.Int) *sdk.Result {
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, amount)
	msg := stakingtypes.NewMsgUndelegate(delegator.String(), val.String(), unbondAmt)
	wmsg := types.NewMsgWrappedUndelegate(msg)
	return h.Handle(func(ctx sdk.Context) (proto.Message, error) {
		return h.MsgSrvr.WrappedUndelegate(ctx, wmsg)
	})
}

// WrappedBeginRedelegate calls handler to redelegate some stake from a validator to another
func (h *Helper) WrappedBeginRedelegate(delegator sdk.AccAddress, srcVal sdk.ValAddress, dstVal sdk.ValAddress, amount math.Int) *sdk.Result {
	unbondAmt := sdk.NewCoin(sdk.DefaultBondDenom, amount)
	msg := stakingtypes.NewMsgBeginRedelegate(delegator.String(), srcVal.String(), dstVal.String(), unbondAmt)
	wmsg := types.NewMsgWrappedBeginRedelegate(msg)
	return h.Handle(func(ctx sdk.Context) (proto.Message, error) {
		return h.MsgSrvr.WrappedBeginRedelegate(ctx, wmsg)
	})
}

// Handle executes an action function with the Helper's context, wraps the result into an SDK service result, and performs two assertions before returning it
func (h *Helper) Handle(action func(sdk.Context) (proto.Message, error)) *sdk.Result {
	res, err := action(h.Ctx)
	r, _ := sdk.WrapServiceResult(h.Ctx, res, err)
	require.NotNil(h.t, r)
	require.NoError(h.t, err)
	return r
}

// CheckValidator asserts that a validor exists and has a given status (if status!="")
// and if has a right jailed flag.
func (h *Helper) CheckValidator(addr sdk.ValAddress, status stakingtypes.BondStatus, jailed bool) stakingtypes.Validator {
	v, err := h.StakingKeeper.GetValidator(h.Ctx, addr)
	require.NoError(h.t, err)
	require.Equal(h.t, jailed, v.Jailed, "wrong Jalied status")
	if status >= 0 {
		require.Equal(h.t, status, v.Status)
	}
	return v
}

// CheckDelegator asserts that a delegator exists
func (h *Helper) CheckDelegator(delegator sdk.AccAddress, val sdk.ValAddress, found bool) {
	_, ok := h.StakingKeeper.GetDelegation(h.Ctx, delegator, val)
	require.Equal(h.t, ok, found)
}
