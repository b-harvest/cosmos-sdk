package testckpt

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/babylonchain/babylon/app"
	appparams "github.com/babylonchain/babylon/app/params"
	"github.com/babylonchain/babylon/testutil/datagen"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/gogoproto/proto"
	"github.com/stretchr/testify/require"

	epochingkeeper "github.com/cosmos/cosmos-sdk/x/epoching/keeper"

	"github.com/cosmos/cosmos-sdk/crypto/keys/bls12381"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/bls/keeper"
	"github.com/cosmos/cosmos-sdk/x/bls/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Helper is a structure which wraps the entire app and exposes functionalities for testing the epoching module
type Helper struct {
	t *testing.T

	Ctx                 sdk.Context
	App                 *app.BabylonApp
	CheckpointingKeeper *keeper.Keeper
	MsgSrvr             types.MsgServer
	QueryClient         types.QueryClient
	StakingKeeper       *stakingkeeper.Keeper
	EpochingKeeper      *epochingkeeper.Keeper

	GenAccs []authtypes.GenesisAccount
}

// NewHelper creates the helper for testing the epoching module
func NewHelper(t *testing.T, n int) *Helper {
	accs, balances := datagen.GenRandomAccWithBalance(n)
	app := app.SetupWithGenesisAccounts(accs, balances...)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	checkpointingKeeper := app.CheckpointingKeeper
	epochingKeeper := app.EpochingKeeper
	stakingKeeper := app.StakingKeeper
	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	types.RegisterQueryServer(queryHelper, checkpointingKeeper)
	queryClient := types.NewQueryClient(queryHelper)
	msgSrvr := keeper.NewMsgServerImpl(checkpointingKeeper)

	return &Helper{
		t:                   t,
		Ctx:                 ctx,
		App:                 app,
		CheckpointingKeeper: &checkpointingKeeper,
		MsgSrvr:             msgSrvr,
		QueryClient:         queryClient,
		StakingKeeper:       stakingKeeper,
		EpochingKeeper:      &epochingKeeper,
		GenAccs:             accs,
	}
}

// CreateValidator calls handler to create a new staking validator
func (h *Helper) CreateValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, blsPK *bls12381.PublicKey, pop *types.ProofOfPossession, stakeAmount math.Int, ok bool) {
	coin := sdk.NewCoin(appparams.DefaultBondDenom, stakeAmount)
	h.createValidator(addr, pk, blsPK, pop, coin, ok)
}

// CreateValidatorWithValPower calls handler to create a new staking validator with zero commission
func (h *Helper) CreateValidatorWithValPower(addr sdk.ValAddress, pk cryptotypes.PubKey, blsPK *bls12381.PublicKey, pop *types.ProofOfPossession, valPower int64, ok bool) math.Int {
	amount := h.StakingKeeper.TokensFromConsensusPower(h.Ctx, valPower)
	coin := sdk.NewCoin(appparams.DefaultBondDenom, amount)
	h.createValidator(addr, pk, blsPK, pop, coin, ok)
	return amount
}

// CreateValidatorMsg returns a message used to create validator in this service.
func (h *Helper) CreateValidatorMsg(addr sdk.ValAddress, pk cryptotypes.PubKey, blsPK *bls12381.PublicKey, pop *types.ProofOfPossession, stakeAmount math.Int) *types.MsgWrappedCreateValidator {
	coin := sdk.NewCoin(appparams.DefaultBondDenom, stakeAmount)
	msg, err := stakingtypes.NewMsgCreateValidator(addr, pk, coin, stakingtypes.Description{}, ZeroCommission(), sdk.OneInt())
	require.NoError(h.t, err)
	wmsg, err := types.NewMsgWrappedCreateValidator(msg, blsPK, pop)
	require.NoError(h.t, err)
	return wmsg
}

func (h *Helper) createValidator(addr sdk.ValAddress, pk cryptotypes.PubKey, blsPK *bls12381.PublicKey, pop *types.ProofOfPossession, coin sdk.Coin, ok bool) {
	h.Handle(func(ctx sdk.Context) (proto.Message, error) {
		return h.CreateValidatorMsg(addr, pk, blsPK, pop, coin.Amount), nil
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

// ZeroCommission constructs a commission rates with all zeros.
func ZeroCommission() stakingtypes.CommissionRates {
	return stakingtypes.NewCommissionRates(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
}
